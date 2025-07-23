package cardhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	cardapierrors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/card"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type cardCommandHandleApi struct {
	card pb.CardCommandServiceClient
	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardCommandResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// cardHandleDeps contains the dependencies required to initialize the card handler.
//
// This struct is typically passed to a constructor function to set up routes
// and initialize the `cardCommandHandleDeps`.
type cardCommandHandleApiDeps struct {
	// client is the gRPC CardService client used for communication.
	client pb.CardCommandServiceClient

	// router is the Echo HTTP router instance used to register endpoints.
	router *echo.Echo

	// logger is used for logging inside the card handler.
	logger logger.LoggerInterface

	// mapper provides a way to transform internal gRPC data into HTTP response models.
	mapper apimapper.CardCommandResponseMapper
}

// NewCardCommandHandleApi initializes a new instance of cardCommandHandleApi with the provided parameters.
// It sets up Prometheus metrics for counting and measuring the duration of card command requests.
//
// Parameters:
// - params: A pointer to cardCommandHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardCommandHandleApi.
func NewCardCommandHandleApi(params *cardCommandHandleApiDeps) *cardCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_command_handler_requests_total",
			Help: "Total number of Card command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_command_handler_request_duration_seconds",
			Help:    "Duration of Card command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	cardHandler := &cardCommandHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-command")

	routerCard.POST("/create", cardHandler.CreateCard)
	routerCard.POST("/update/:id", cardHandler.UpdateCard)
	routerCard.POST("/trashed/:id", cardHandler.TrashedCard)
	routerCard.POST("/restore/:id", cardHandler.RestoreCard)
	routerCard.DELETE("/permanent/:id", cardHandler.DeleteCardPermanent)

	routerCard.POST("/restore/all", cardHandler.RestoreAllCard)
	routerCard.POST("/permanent/all", cardHandler.DeleteAllCardPermanent)

	return cardHandler
}

// @Security Bearer
// @Summary Create a new card
// @Tags Card-Command
// @Description Create a new card for a user
// @Accept json
// @Produce json
// @Param CreateCardRequest body requests.CreateCardRequest true "Create card request"
// @Success 200 {object} response.ApiResponseCard "Created card"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create card"
// @Router /api/card/create [post]
func (h *cardCommandHandleApi) CreateCard(c echo.Context) error {
	const method = "CreateCard"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateCardRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateCard request", err, zap.Error(err))
		return cardapierrors.ErrApiBindCreateCard(c)
	}

	if err := body.Validate(); err != nil {
		logError("Validation failed for CreateCard", err, zap.Error(err))
		return cardapierrors.ErrApiValidateCreateCard(c)
	}

	req := &pb.CreateCardRequest{
		UserId:       int32(body.UserID),
		CardType:     body.CardType,
		ExpireDate:   timestamppb.New(body.ExpireDate),
		Cvv:          body.CVV,
		CardProvider: body.CardProvider,
	}

	res, err := h.card.CreateCard(ctx, req)

	if err != nil {
		logError("CreateCard service failed", err)
		return cardapierrors.ErrApiFailedCreateCard(c)
	}

	response := h.mapper.ToApiResponseCard(res)

	logSuccess("Successfully created card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Update a card
// @Tags Card-Command
// @Description Update a card for a user
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param UpdateCardRequest body requests.UpdateCardRequest true "Update card request"
// @Success 200 {object} response.ApiResponseCard "Updated card"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update card"
// @Router /api/card/update/{id} [post]
func (h *cardCommandHandleApi) UpdateCard(c echo.Context) error {
	const method = "UpdateCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid card id", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidCardID(c)
	}

	var body requests.UpdateCardRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateCard request", err, zap.Error(err))

		return cardapierrors.ErrApiBindUpdateCard(c)
	}

	if err := body.Validate(); err != nil {
		logError("Validation failed for UpdateCard", err, zap.Error(err))

		return cardapierrors.ErrApiValidateUpdateCard(c)
	}

	req := &pb.UpdateCardRequest{
		CardId:       int32(idInt),
		UserId:       int32(body.UserID),
		CardType:     body.CardType,
		ExpireDate:   timestamppb.New(body.ExpireDate),
		Cvv:          body.CVV,
		CardProvider: body.CardProvider,
	}

	res, err := h.card.UpdateCard(ctx, req)

	if err != nil {
		logError("UpdateCard service failed", err, zap.Error(err))

		return cardapierrors.ErrApiFailedUpdateCard(c)
	}

	so := h.mapper.ToApiResponseCard(res)

	logSuccess("Successfully updated card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Trashed a card
// @Tags Card-Command
// @Description Trashed a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Trashed card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed card"
// @Router /api/card/trashed/{id} [post]
func (h *cardCommandHandleApi) TrashedCard(c echo.Context) error {
	const method = "TrashedCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid card id", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.TrashedCard(ctx, req)

	if err != nil {
		logError("TrashedCard service failed", err, zap.Error(err))

		return cardapierrors.ErrApiFailedTrashCard(c)
	}

	so := h.mapper.ToApiResponseCardDeleteAt(res)

	logSuccess("Successfully trashed card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore a card
// @Tags Card-Command
// @Description Restore a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Restored card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore card"
// @Router /api/card/restore/{id} [post]
func (h *cardCommandHandleApi) RestoreCard(c echo.Context) error {
	const method = "RestoreCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid card id", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.RestoreCard(ctx, req)

	if err != nil {
		logError("RestoreCard service failed", err, zap.Error(err))

		return cardapierrors.ErrApiFailedRestoreCard(c)
	}

	so := h.mapper.ToApiResponseCard(res)

	logSuccess("Successfully restored card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete a card permanently
// @Tags Card-Command
// @Description Delete a card by its ID permanently
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCardDelete "Deleted card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete card"
// @Router /api/card/permanent/{id} [delete]
func (h *cardCommandHandleApi) DeleteCardPermanent(c echo.Context) error {
	const method = "DeleteCardPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid card id", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.DeleteCardPermanent(ctx, req)

	if err != nil {
		logError("Failed to delete card permanently", err, zap.Error(err))

		return cardapierrors.ErrApiFailedDeleteCardPermanent(c)
	}

	so := h.mapper.ToApiResponseCardDelete(res)

	logSuccess("Successfully deleted card permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore all card records
// @Tags Card-Command
// @Description Restore all card records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCardAll "Successfully restored all card records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all card records"
// @Router /api/card/restore/all [post]
func (h *cardCommandHandleApi) RestoreAllCard(c echo.Context) error {
	const method = "RestoreAllCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.card.RestoreAllCard(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all cards", err, zap.Error(err))

		return cardapierrors.ErrApiFailedRestoreAllCard(c)
	}

	h.logger.Debug("Successfully restored all cards")

	so := h.mapper.ToApiResponseCardAll(res)

	logSuccess("Successfully restored all cards", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer.
// @Summary Permanently delete all card records
// @Tags Card-Command
// @Description Permanently delete all card records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCardAll "Successfully deleted all card records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all card records"
// @Router /api/card/permanent/all [post]
func (h *cardCommandHandleApi) DeleteAllCardPermanent(c echo.Context) error {
	const method = "DeleteAllCardPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.card.DeleteAllCardPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all cards permanently", err, zap.Error(err))

		return cardapierrors.ErrApiFailedDeleteAllCardPermanent(c)
	}

	so := h.mapper.ToApiResponseCardAll(res)

	logSuccess("Successfully deleted all cards permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *cardCommandHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (
	end func(),
	logSuccess func(string, ...zap.Field),
	logError func(string, error, ...zap.Field),
) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	status := "success"

	end = func() {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess = func(msg string, fields ...zap.Field) {
		status = "success"
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError = func(msg string, err error, fields ...zap.Field) {
		status = "error"
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *cardCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

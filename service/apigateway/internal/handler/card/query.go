package cardhandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	cardapierrors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/card"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardQueryHandleApi struct {
	// card is the gRPC client used to interact with the CardService.
	card pb.CardQueryServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardQueryResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// cardQueryHandleApiDeps contains the necessary dependencies for initializing a cardQueryHandleApi.
type cardQueryHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardService.
	client pb.CardQueryServiceClient

	// router is the Echo HTTP router used to register endpoints.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardQueryResponseMapper
}

// NewCardQueryHandleApi initializes a new instance of cardQueryHandleApi with the provided parameters.
//
// It sets up Prometheus metrics for counting and measuring the duration of card query requests.
//
// Parameters:
// - params: A pointer to cardQueryHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardQueryHandleApi.
func NewCardQueryHandleApi(
	params *cardQueryHandleApiDeps,
) *cardQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_query_handler_requests_total",
			Help: "Total number of Card Query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_query_handler_request_duration_seconds",
			Help:    "Duration of Card Query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	cardQueryHandler := &cardQueryHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-query")

	routerCard.GET("", cardQueryHandler.FindAll)
	routerCard.GET("/:id", cardQueryHandler.FindById)

	routerCard.GET("/user", cardQueryHandler.FindByUserID)
	routerCard.GET("/active", cardQueryHandler.FindByActive)
	routerCard.GET("/trashed", cardQueryHandler.FindByTrashed)
	routerCard.GET("/card_number/:card_number", cardQueryHandler.FindByCardNumber)

	return cardQueryHandler
}

// FindAll godoc
// @Summary Retrieve all cards
// @Tags Card-Query
// @Security Bearer
// @Description Retrieve all cards with pagination
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Number of data per page"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCard "Card data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card data"
// @Router /api/card [get]
func (h *cardQueryHandleApi) FindAll(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAll"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	cards, err := h.card.FindAllCard(ctx, req)

	if err != nil {
		logError("Failed to find all cards", err,
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
		return cardapierrors.ErrApiFailedFindAllCards(c)
	}

	response := h.mapper.ToApiResponsesCard(cards)

	logSuccess("Successfully retrieved card list",
		zap.Int("count", len(response.Data)),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Bool("success", true),
	)

	return c.JSON(http.StatusOK, response)
}

// FindById godoc
// @Summary Retrieve card by ID
// @Tags Card-Query
// @Security Bearer
// @Description Retrieve a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid card ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/{id} [get]
func (h *cardQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid card ID", err, zap.Error(err))
		return cardapierrors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(id),
	}

	card, err := h.card.FindByIdCard(ctx, req)
	if err != nil {
		logError("FindByIdCard failed", err, zap.Int("card.id", id), zap.Error(err))
		return cardapierrors.ErrApiFailedFindByIdCard(c)
	}

	response := h.mapper.ToApiResponseCard(card)

	logSuccess("Successfully retrieved card record", zap.Int("card.id", id), zap.Bool("success", true))

	return c.JSON(http.StatusOK, response)
}

// FindByUserID godoc
// @Summary Retrieve cards by user ID
// @Tags Card-Query
// @Security Bearer
// @Description Retrieve a list of cards associated with a user by their ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/user [get]
func (h *cardQueryHandleApi) FindByUserID(c echo.Context) error {
	const method = "FindByUserID"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int32)
	if !ok {
		err := errors.New("user id not found in context")

		logError("Invalid user id", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidUserID(c)
	}

	req := &pb.FindByUserIdCardRequest{
		UserId: userID,
	}

	card, err := h.card.FindByUserIdCard(ctx, req)

	if err != nil {
		logError("FindByUserIdCard failed", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindByUserIdCard(c)
	}

	so := h.mapper.ToApiResponseCard(card)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active card by Saldo ID
// @Tags Card-Query
// @Description Retrieve an active card associated with a Saldo ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationCardDeleteAt "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid Saldo ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/active [get]
func (h *cardQueryHandleApi) FindByActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActive"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.card.FindByActiveCard(ctx, req)

	if err != nil {
		logError("Failed to retrieve card record", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindByActiveCard(c)
	}

	so := h.mapper.ToApiResponsesCardDeletedAt(res)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed cards
// @Tags Card-Query
// @Security Bearer
// @Description Retrieve a list of trashed cards
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCardDeleteAt "Card data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/trashed [get]
func (h *cardQueryHandleApi) FindByTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashed"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.card.FindByTrashedCard(ctx, req)

	if err != nil {
		logError("Failed to retrieve card record", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindByTrashedCard(c)
	}

	so := h.mapper.ToApiResponsesCardDeletedAt(res)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve card by card number
// @Tags Card-Query
// @Description Retrieve a card by its card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card number"
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Failed to fetch card record"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/{card_number} [get]
func (h *cardQueryHandleApi) FindByCardNumber(c echo.Context) error {
	const method = "FindByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber := c.Param("card_number")

	if cardNumber == "" {
		err := errors.New("invalid card number")

		logError("Failed to fetch card record", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidCardNumber(c)
	}

	req := &pbhelpers.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.card.FindByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve card record", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindByCardNumber(c)
	}

	so := h.mapper.ToApiResponseCard(res)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *cardQueryHandleApi) startTracingAndLogging(
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

func (s *cardQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

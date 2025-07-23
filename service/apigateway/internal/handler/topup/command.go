package topuphandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/topup"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupCommandHandleApi struct {
	client pb.TopupCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TopupCommandResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type topupCommandHandleDeps struct {
	client pb.TopupCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TopupCommandResponseMapper
}

func NewTopupCommandHandleApi(params *topupCommandHandleDeps) *topupCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_command_handler_requests_total",
			Help: "Total number of topup command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_command_handler_request_duration_seconds",
			Help:    "Duration of topup command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	topupHandler := &topupCommandHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("topup-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTopup := params.router.Group("/api/topup-command")

	routerTopup.POST("/create", topupHandler.Create)
	routerTopup.POST("/update/:id", topupHandler.Update)
	routerTopup.POST("/trashed/:id", topupHandler.TrashTopup)
	routerTopup.POST("/restore/:id", topupHandler.RestoreTopup)
	routerTopup.DELETE("/permanent/:id", topupHandler.DeleteTopupPermanent)

	routerTopup.POST("/trashed/all", topupHandler.DeleteAllTopupPermanent)
	routerTopup.POST("/restore/all", topupHandler.RestoreAllTopup)
	return topupHandler
}

// @Summary Create topup
// @Tags Topup Command
// @Security Bearer
// @Description Create a new topup record
// @Accept json
// @Produce json
// @Param CreateTopupRequest body requests.CreateTopupRequest true "Create topup request"
// @Success 200 {object} response.ApiResponseTopup "Created topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: "
// @Failure 500 {object} response.ErrorResponse "Failed to create topup: "
// @Router /api/topup-command/create [post]
func (h *topupCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateTopupRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateTopup request", err, zap.Error(err))

		return topup_errors.ErrApiBindCreateTopup(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate CreateTopup request", err, zap.Error(err))

		return topup_errors.ErrApiValidateCreateTopup(c)
	}

	res, err := h.client.CreateTopup(ctx, &pb.CreateTopupRequest{
		CardNumber:  body.CardNumber,
		TopupAmount: int32(body.TopupAmount),
		TopupMethod: body.TopupMethod,
	})

	if err != nil {
		logError("Failed to create topup", err, zap.Error(err))

		return topup_errors.ErrApiFailedCreateTopup(c)
	}

	so := h.mapper.ToApiResponseTopup(res)

	logSuccess("success create topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update topup
// @Tags Topup Command
// @Security Bearer
// @Description Update an existing topup record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Param UpdateTopupRequest body requests.UpdateTopupRequest true "Update topup request"
// @Success 200 {object} response.ApiResponseTopup "Updated topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid input data"
// @Failure 500 {object} response.ErrorResponse "Failed to update topup: "
// @Router /api/topup-command/update/{id} [post]
func (h *topupCommandHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	idint, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	var body requests.UpdateTopupRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateTopup request", err, zap.Error(err))

		return topup_errors.ErrApiBindUpdateTopup(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate UpdateTopup request", err, zap.Error(err))

		return topup_errors.ErrApiValidateUpdateTopup(c)
	}

	res, err := h.client.UpdateTopup(ctx, &pb.UpdateTopupRequest{
		TopupId:     int32(idint),
		CardNumber:  body.CardNumber,
		TopupAmount: int32(body.TopupAmount),
		TopupMethod: body.TopupMethod,
	})

	if err != nil {
		logError("Failed to update topup", err, zap.Error(err))

		return topup_errors.ErrApiFailedUpdateTopup(c)
	}

	so := h.mapper.ToApiResponseTopup(res)

	logSuccess("success update topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a topup
// @Tags Topup Command
// @Security Bearer
// @Description Trash a topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Successfully trashed topup record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trash topup:"
// @Router /api/topup-command/trash/{id} [post]
func (h *topupCommandHandleApi) TrashTopup(c echo.Context) error {
	const method = "TrashTopup"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	res, err := h.client.TrashedTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		logError("Failed to trash topup", err, zap.Error(err))

		return topup_errors.ErrApiFailedTrashTopup(c)
	}

	so := h.mapper.ToApiResponseTopupDeleteAt(res)

	logSuccess("success trash topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed topup
// @Tags Topup Command
// @Security Bearer
// @Description Restore a trashed topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Successfully restored topup record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore topup:"
// @Router /api/topup-command/restore/{id} [post]
func (h *topupCommandHandleApi) RestoreTopup(c echo.Context) error {
	const method = "RestoreTopup"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	res, err := h.client.RestoreTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		logError("Failed to restore topup", err, zap.Error(err))

		return topup_errors.ErrApiFailedRestoreTopup(c)
	}

	so := h.mapper.ToApiResponseTopup(res)

	logSuccess("success restore topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a topup
// @Tags Topup Command
// @Security Bearer
// @Description Permanently delete a topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopupDelete "Successfully deleted topup record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete topup:"
// @Router /api/topup-command/permanent/{id} [delete]
func (h *topupCommandHandleApi) DeleteTopupPermanent(c echo.Context) error {
	const method = "DeleteTopupPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	res, err := h.client.DeleteTopupPermanent(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		logError("Failed to delete topup", err, zap.Error(err))

		return topup_errors.ErrApiFailedDeletePermanentTopup(c)
	}

	so := h.mapper.ToApiResponseTopupDelete(res)

	logSuccess("success delete topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore all topup records
// @Tags Topup Command
// @Security Bearer
// @Description Restore all topup records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTopupAll "Successfully restored all topup records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all topup records"
// @Router /api/topup-command/restore/all [post]
func (h *topupCommandHandleApi) RestoreAllTopup(c echo.Context) error {
	const method = "RestoreAllTopup"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.RestoreAllTopup(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all topup", err, zap.Error(err))

		return topup_errors.ErrApiFailedRestoreAllTopup(c)
	}

	so := h.mapper.ToApiResponseTopupAll(res)

	logSuccess("success restore all topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete all topup records
// @Tags Topup Command
// @Security Bearer
// @Description Permanently delete all topup records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTopupAll "Successfully deleted all topup records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all topup records"
// @Router /api/topup-command/permanent/all [post]
func (h *topupCommandHandleApi) DeleteAllTopupPermanent(c echo.Context) error {
	const method = "DeleteAllTopupPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.DeleteAllTopupPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all topup permanently", err, zap.Error(err))

		return topup_errors.ErrApiFailedDeleteAllTopupPermanent(c)
	}

	so := h.mapper.ToApiResponseTopupAll(res)

	logSuccess("success delete all topup", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *topupCommandHandleApi) startTracingAndLogging(
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

func (s *topupCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

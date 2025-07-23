package transferhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transfer"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transferCommandHandleApi struct {
	client pb.TransferCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransferCommandResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transferCommandHandleDeps struct {
	client pb.TransferCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransferCommandResponseMapper
}

func NewTransferCommandHandleApi(params *transferCommandHandleDeps) *transferCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_command_handler_requests_total",
			Help: "Total number of transfer command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_command_handler_request_duration_seconds",
			Help:    "Duration of transfer command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	transferCommandHandleApi := &transferCommandHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transfer-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransfer := params.router.Group("/api/transfer-command")

	routerTransfer.POST("/create", transferCommandHandleApi.CreateTransfer)
	routerTransfer.POST("/update/:id", transferCommandHandleApi.UpdateTransfer)
	routerTransfer.POST("/trashed/:id", transferCommandHandleApi.TrashTransfer)
	routerTransfer.POST("/restore/:id", transferCommandHandleApi.RestoreTransfer)
	routerTransfer.DELETE("/permanent/:id", transferCommandHandleApi.DeleteTransferPermanent)

	routerTransfer.POST("/restore/all", transferCommandHandleApi.RestoreAllTransfer)
	routerTransfer.POST("/permanent/all", transferCommandHandleApi.DeleteAllTransferPermanent)

	return transferCommandHandleApi
}

// @Summary Create a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Create a new transfer record
// @Accept json
// @Produce json
// @Param body body requests.CreateTransferRequest true "Transfer request"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Validation Error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transfer"
// @Router /api/transfer-command/create [post]
func (h *transferCommandHandleApi) CreateTransfer(c echo.Context) error {
	const method = "CreateTransfer"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateTransferRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateTransfer request", err, zap.Error(err))

		return transfer_errors.ErrApiBindCreateTransfer(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate CreateTransfer request", err, zap.Error(err))

		return transfer_errors.ErrApiValidateCreateTransfer(c)
	}

	res, err := h.client.CreateTransfer(ctx, &pb.CreateTransferRequest{
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	})

	if err != nil {
		logError("Failed to create transfer", err, zap.Error(err))

		return transfer_errors.ErrApiFailedCreateTransfer(c)
	}

	so := h.mapper.ToApiResponseTransfer(res)

	logSuccess("Successfully created transfer", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Update an existing transfer record
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param body body requests.UpdateTransferRequest true "Transfer request"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Validation Error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transfer"
// @Router /api/transfer-command/update/{id} [post]
func (h *transferCommandHandleApi) UpdateTransfer(c echo.Context) error {
	const method = "UpdateTransfer"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transfer_errors.ErrApiTransferInvalidID(c)
	}

	var body requests.UpdateTransferRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateTransfer request", err, zap.Error(err))

		return transfer_errors.ErrApiBindUpdateTransfer(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate UpdateTransfer request", err, zap.Error(err))

		return transfer_errors.ErrApiValidateUpdateTransfer(c)
	}

	res, err := h.client.UpdateTransfer(ctx, &pb.UpdateTransferRequest{
		TransferId:     int32(idInt),
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	})

	if err != nil {
		logError("Failed to update transfer", err, zap.Error(err))

		return transfer_errors.ErrApiFailedUpdateTransfer(c)
	}

	so := h.mapper.ToApiResponseTransfer(res)

	logSuccess("Successfully updated transfer", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Soft delete a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Soft delete a transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Successfully trashed transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transfer"
// @Router /api/transfer-command/trash/{id} [post]
func (h *transferCommandHandleApi) TrashTransfer(c echo.Context) error {
	const method = "TrashTransfer"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transfer_errors.ErrApiTransferInvalidID(c)
	}

	res, err := h.client.TrashedTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		logError("Failed to trashed transfer", err, zap.Error(err))

		return transfer_errors.ErrApiFailedTrashedTransfer(c)
	}

	so := h.mapper.ToApiResponseTransferDeleteAt(res)

	logSuccess("Successfully trashed transfer", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Restore a trashed transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Successfully restored transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transfer:"
// @Router /api/transfer-command/restore/{id} [post]
func (h *transferCommandHandleApi) RestoreTransfer(c echo.Context) error {
	const method = "RestoreTransfer"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transfer_errors.ErrApiTransferInvalidID(c)
	}

	res, err := h.client.RestoreTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		logError("Failed to restore transfer", err, zap.Error(err))

		return transfer_errors.ErrApiFailedRestoreTransfer(c)
	}

	so := h.mapper.ToApiResponseTransfer(res)

	logSuccess("Successfully restored transfer", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Permanently delete a transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransferDelete "Successfully deleted transfer record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transfer:"
// @Router /api/transfer-command/permanent/{id} [delete]
func (h *transferCommandHandleApi) DeleteTransferPermanent(c echo.Context) error {
	const method = "DeleteTransferPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transfer_errors.ErrApiTransferInvalidID(c)
	}

	res, err := h.client.DeleteTransferPermanent(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		logError("Failed to delete transfer", err, zap.Error(err))

		return transfer_errors.ErrApiFailedDeleteTransferPermanent(c)
	}

	so := h.mapper.ToApiResponseTransferDelete(res)

	logSuccess("Successfully deleted transfer permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Restore a trashed transfer all
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransferAll "Successfully restored transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transfer:"
// @Router /api/transfer-command/restore/all [post]
func (h *transferCommandHandleApi) RestoreAllTransfer(c echo.Context) error {
	const method = "RestoreAllTransfer"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.RestoreAllTransfer(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all transfer", err, zap.Error(err))

		return transfer_errors.ErrApiFailedRestoreAllTransfer(c)
	}

	so := h.mapper.ToApiResponseTransferAll(res)

	logSuccess("Successfully restored all transfer", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Permanently delete a transfer record all.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransferAll "Successfully deleted transfer all"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transfer:"
// @Router /api/transfer-command/permanent/all [post]
func (h *transferCommandHandleApi) DeleteAllTransferPermanent(c echo.Context) error {
	const method = "DeleteAllTransferPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.DeleteAllTransferPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all transfer permanently", err, zap.Error(err))

		return transfer_errors.ErrApiFailedDeleteAllTransferPermanent(c)
	}

	so := h.mapper.ToApiResponseTransferAll(res)

	logSuccess("Successfully deleted all transfer permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transferCommandHandleApi) startTracingAndLogging(
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

func (s *transferCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

package withdrawhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/withdraw"
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

type withdrawCommandHandleApi struct {
	client pb.WithdrawCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.WithdrawCommandResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type withdrawCommandHandleDeps struct {
	client pb.WithdrawCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.WithdrawCommandResponseMapper
}

func NewWithdrawCommandHandleApi(params *withdrawCommandHandleDeps) *withdrawCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_command_handler_requests_total",
			Help: "Total number of withdraw command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_command_handler_request_duration_seconds",
			Help:    "Duration of withdraw command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	withdrawCommandHandleApi := &withdrawCommandHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("withdraw-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerWithdraw := params.router.Group("/api/withdraw-command")

	routerWithdraw.POST("/create", withdrawCommandHandleApi.Create)
	routerWithdraw.POST("/update/:id", withdrawCommandHandleApi.Update)

	routerWithdraw.POST("/trashed/:id", withdrawCommandHandleApi.TrashWithdraw)
	routerWithdraw.POST("/restore/:id", withdrawCommandHandleApi.RestoreWithdraw)
	routerWithdraw.DELETE("/permanent/:id", withdrawCommandHandleApi.DeleteWithdrawPermanent)

	routerWithdraw.POST("/restore/all", withdrawCommandHandleApi.RestoreAllWithdraw)
	routerWithdraw.POST("/permanent/all", withdrawCommandHandleApi.DeleteAllWithdrawPermanent)

	return withdrawCommandHandleApi
}

// @Summary Create a new withdraw
// @Tags Withdraw Command
// @Security Bearer
// @Description Create a new withdraw record with the provided details.
// @Accept json
// @Produce json
// @Param CreateWithdrawRequest body requests.CreateWithdrawRequest true "Create Withdraw Request"
// @Success 200 {object} response.ApiResponseWithdraw "Successfully created withdraw record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create withdraw"
// @Router /api/withdraw-command/create [post]
func (h *withdrawCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateWithdrawRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateWithdraw request", err, zap.Error(err))

		return withdraw_errors.ErrApiBindCreateWithdraw(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate CreateWithdraw request", err, zap.Error(err))

		return withdraw_errors.ErrApiValidateCreateWithdraw(c)
	}

	res, err := h.client.CreateWithdraw(ctx, &pb.CreateWithdrawRequest{
		CardNumber:     body.CardNumber,
		WithdrawAmount: int32(body.WithdrawAmount),
		WithdrawTime:   timestamppb.New(body.WithdrawTime),
	})

	if err != nil {
		logError("Failed to create withdraw", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedCreateWithdraw(c)
	}

	so := h.mapper.ToApiResponseWithdraw(res)

	logSuccess("Success create withdraw", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update an existing withdraw
// @Tags Withdraw Command
// @Security Bearer
// @Description Update an existing withdraw record with the provided details.
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Param UpdateWithdrawRequest body requests.UpdateWithdrawRequest true "Update Withdraw Request"
// @Success 200 {object} response.ApiResponseWithdraw "Successfully updated withdraw record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update withdraw"
// @Router /api/withdraw-command/update/{id} [post]
func (h *withdrawCommandHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid withdraw ID", err, zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	var body requests.UpdateWithdrawRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateWithdraw request", err, zap.Error(err))

		return withdraw_errors.ErrApiBindUpdateWithdraw(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate UpdateWithdraw request", err, zap.Error(err))

		return withdraw_errors.ErrApiValidateUpdateWithdraw(c)
	}

	res, err := h.client.UpdateWithdraw(ctx, &pb.UpdateWithdrawRequest{
		WithdrawId:     int32(id),
		CardNumber:     body.CardNumber,
		WithdrawAmount: int32(body.WithdrawAmount),
		WithdrawTime:   timestamppb.New(body.WithdrawTime),
	})

	if err != nil {
		logError("Failed to update withdraw", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedUpdateWithdraw(c)
	}

	so := h.mapper.ToApiResponseWithdraw(res)

	logSuccess("Success update withdraw", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a withdraw by ID
// @Tags Withdraw Command
// @Security Bearer
// @Description Trash a withdraw using its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdaw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trash withdraw"
// @Router /api/withdraw-command/trashed/{id} [post]
func (h *withdrawCommandHandleApi) TrashWithdraw(c echo.Context) error {
	const method = "TrashWithdraw"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid withdraw ID", err, zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	res, err := h.client.TrashedWithdraw(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		logError("Failed to trash withdraw", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedTrashedWithdraw(c)
	}

	so := h.mapper.ToApiResponseWithdrawDeleteAt(res)

	logSuccess("Success trash withdraw", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a withdraw by ID
// @Tags Withdraw Command
// @Security Bearer
// @Description Restore a withdraw by its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore withdraw"
// @Router /api/withdraw-command/restore/{id} [post]
func (h *withdrawCommandHandleApi) RestoreWithdraw(c echo.Context) error {
	const method = "RestoreWithdraw"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid withdraw ID", err, zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	res, err := h.client.RestoreWithdraw(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		logError("Failed to restore withdraw", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedRestoreWithdraw(c)
	}

	so := h.mapper.ToApiResponseWithdraw(res)

	logSuccess("Success restore withdraw", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a withdraw by ID
// @Tags Withdraw Command
// @Security Bearer
// @Description Permanently delete a withdraw by its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdrawDelete "Successfully deleted withdraw permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete withdraw permanently:"
// @Router /api/withdraw-command/permanent/{id} [delete]
func (h *withdrawCommandHandleApi) DeleteWithdrawPermanent(c echo.Context) error {
	const method = "DeleteWithdrawPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid withdraw ID", err, zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	res, err := h.client.DeleteWithdrawPermanent(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		logError("Failed to delete withdraw permanently", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedDeleteWithdrawPermanent(c)
	}

	so := h.mapper.ToApiResponseWithdrawDelete(res)

	logSuccess("Success delete withdraw permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a withdraw all
// @Tags Withdraw Command
// @Security Bearer
// @Description Restore a withdraw all
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseWithdrawAll "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore withdraw"
// @Router /api/withdraw-command/restore/all [post]
func (h *withdrawCommandHandleApi) RestoreAllWithdraw(c echo.Context) error {
	const method = "RestoreAllWithdraw"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.RestoreAllWithdraw(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all withdraw", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedRestoreAllWithdraw(c)
	}

	so := h.mapper.ToApiResponseWithdrawAll(res)

	logSuccess("Success restore all withdraw", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a withdraw by ID
// @Tags Withdraw Command
// @Security Bearer
// @Description Permanently delete a withdraw by its ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseWithdrawAll "Successfully deleted withdraw permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete withdraw permanently:"
// @Router /api/withdraw-command/permanent/all [post]
func (h *withdrawCommandHandleApi) DeleteAllWithdrawPermanent(c echo.Context) error {
	const method = "DeleteAllWithdrawPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.DeleteAllWithdrawPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all withdraw permanently", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedDeleteAllWithdrawPermanent(c)
	}

	so := h.mapper.ToApiResponseWithdrawAll(res)

	logSuccess("Success delete all withdraw permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *withdrawCommandHandleApi) startTracingAndLogging(
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

func (s *withdrawCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

package merchanthandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchant"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantCommandHandleApi struct {
	// client is the gRPC client used to interact with the MerchantService.
	client pb.MerchantCommandServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.MerchantCommandResponseMapper

	// trace provides distributed tracing capabilities.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// merchantCommandHandleDeps contains the necessary dependencies for the NewMerchantCommandHandleApi function.
type merchantCommandHandleDeps struct {
	// client is the gRPC client used to interact with the MerchantCommandServiceClient.
	client pb.MerchantCommandServiceClient

	// router is the Echo router used to register HTTP routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.MerchantCommandResponseMapper
}

// NewMerchantCommandHandleApi initializes a new merchantCommandHandleApi and sets up the routes for merchant-related operations.
//
// This function registers various HTTP endpoints related to merchant management, including creation, updating, deletion, and retrieval of merchants.
// It also collects Prometheus metrics for tracking the number of requests and their durations, helping to monitor the performance and reliability of the handler.
// The routes are grouped under "/api/merchant-command".
//
// Parameters:
// - params: A pointer to merchantCommandHandleDeps containing the necessary dependencies such as router, logger, client, and mapper.
//
// Returns:
// - A pointer to a newly created merchantCommandHandleApi.
func NewMerchantCommandHandleApi(params *merchantCommandHandleDeps) *merchantCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_command_handler_requests_total",
			Help: "Total number of merchant command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_command_handler_request_duration_seconds",
			Help:    "Duration of merchant command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantHandler := &merchantCommandHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("merchant-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := params.router.Group("/api/merchant-command")

	routerMerchant.POST("/create", merchantHandler.Create)
	routerMerchant.POST("/updates/:id", merchantHandler.Update)
	routerMerchant.POST("/update-status/:id", merchantHandler.UpdateStatus)

	routerMerchant.POST("/trashed/:id", merchantHandler.TrashedMerchant)
	routerMerchant.POST("/restore/:id", merchantHandler.RestoreMerchant)
	routerMerchant.DELETE("/permanent/:id", merchantHandler.Delete)

	routerMerchant.POST("/restore/all", merchantHandler.RestoreAllMerchant)
	routerMerchant.POST("/permanent/all", merchantHandler.DeleteAllMerchantPermanent)

	return merchantHandler
}

// Create godoc
// @Summary Create a new merchant
// @Tags Merchant
// @Security Bearer
// @Description Create a new merchant with the given name and user ID
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantRequest true "Create merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Created merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant"
// @Router /api/merchant-command/create [post]
func (h *merchantCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateMerchantRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateMerchant request", err, zap.Error(err))

		return merchant_errors.ErrApiBindCreateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		logError("Validation Error", err, zap.Error(err))

		return merchant_errors.ErrApiValidateCreateMerchant(c)
	}

	req := &pb.CreateMerchantRequest{
		Name:   body.Name,
		UserId: int32(body.UserID),
	}

	res, err := h.client.CreateMerchant(ctx, req)

	if err != nil {
		logError("Failed to create merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedCreateMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchant(res)

	logSuccess("Merchant created successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Update godoc
// @Summary Update a merchant
// @Tags Merchant
// @Security Bearer
// @Description Update a merchant with the given ID
// @Accept json
// @Produce json
// @Param body body requests.UpdateMerchantRequest true "Update merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Updated merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant"
// @Router /api/merchant-command/update/{id} [post]
func (h *merchantCommandHandleApi) Update(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	var body requests.UpdateMerchantRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateMerchant request", err, zap.Error(err))

		return merchant_errors.ErrApiBindUpdateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		logError("Validation Error", err, zap.Error(err))

		return merchant_errors.ErrApiValidateUpdateMerchant(c)
	}

	req := &pb.UpdateMerchantRequest{
		MerchantId: int32(id),
		Name:       body.Name,
		UserId:     int32(body.UserID),
		Status:     body.Status,
	}

	res, err := h.client.UpdateMerchant(ctx, req)

	if err != nil {
		logError("Failed to update merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedUpdateMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchant(res)

	logSuccess("Merchant updated successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// UpdateStatus godoc
// @Summary Update merchant status
// @Tags Merchant
// @Security Bearer
// @Description Update the status of a merchant with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param body body requests.UpdateMerchantStatusRequest true "Update merchant status request"
// @Success 200 {object} response.ApiResponseMerchant "Updated merchant status"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant status"
// @Router /api/merchant-command/update-status/{id} [post]
func (h *merchantCommandHandleApi) UpdateStatus(c echo.Context) error {
	const method = "UpdateStatus"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	var body requests.UpdateMerchantStatusRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateMerchantStatus request", err, zap.Error(err))

		return merchant_errors.ErrApiBindUpdateMerchantStatus(c)
	}

	if err := body.Validate(); err != nil {
		logError("Validation Error", err, zap.Error(err))

		return merchant_errors.ErrApiValidateUpdateMerchantStatus(c)
	}

	req := &pb.UpdateMerchantStatusRequest{
		MerchantId: int32(id),
		Status:     body.Status,
	}

	res, err := h.client.UpdateMerchantStatus(ctx, req)

	if err != nil {
		logError("Failed to update merchant status", err, zap.Error(err))

		return merchant_errors.ErrApiFailedUpdateMerchantStatus(c)
	}

	so := h.mapper.ToApiResponseMerchant(res)

	logSuccess("Merchant status updated successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// TrashedMerchant godoc
// @Summary Trashed a merchant
// @Tags Merchant
// @Security Bearer
// @Description Trashed a merchant by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed merchant"
// @Router /api/merchant-command/trashed/{id} [post]
func (h *merchantCommandHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	res, err := h.client.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		logError("Failed to trashed merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedTrashMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchantDeleteAt(res)

	logSuccess("Merchant trashed successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreMerchant godoc
// @Summary Restore a merchant
// @Tags Merchant
// @Security Bearer
// @Description Restore a merchant by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Restored merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-command/restore/{id} [post]
func (h *merchantCommandHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	res, err := h.client.RestoreMerchant(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		logError("Failed to restore merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedRestoreMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchant(res)

	logSuccess("Merchant restored successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Delete godoc
// @Summary Delete a merchant permanently
// @Tags Merchant
// @Security Bearer
// @Description Delete a merchant by its ID permanently
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Deleted merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant"
// @Router /api/merchant-command/{id} [delete]
func (h *merchantCommandHandleApi) Delete(c echo.Context) error {
	const method = "Delete"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	res, err := h.client.DeleteMerchantPermanent(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		logError("Failed to delete merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedDeleteMerchantPermanent(c)
	}

	so := h.mapper.ToApiResponseMerchantDelete(res)

	logSuccess("Merchant deleted successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllMerchant godoc.
// @Summary Restore all merchant records
// @Tags Merchant
// @Security Bearer
// @Description Restore all merchant records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all merchant records"
// @Router /api/merchant-command/restore/all [post]
func (h *merchantCommandHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.RestoreAllMerchant(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedRestoreAllMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchantAll(res)

	logSuccess("Merchant restored successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DeleteAllMerchantPermanent godoc.
// @Summary Permanently delete all merchant records
// @Tags Merchant
// @Security Bearer
// @Description Permanently delete all merchant records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all merchant records"
// @Router /api/merchant-command/permanent/all [post]
func (h *merchantCommandHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "DeleteAllMerchantPermanent"

	end, logSuccess, logError := h.startTracingAndLogging(c.Request().Context(), method)
	defer func() {

		end()
	}()

	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all merchant permanently", err, zap.Error(err))

		return merchant_errors.ErrApiFailedDeleteAllMerchantPermanent(c)
	}

	so := h.mapper.ToApiResponseMerchantAll(res)

	logSuccess("Merchant deleted successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantCommandHandleApi) startTracingAndLogging(
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

func (s *merchantCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

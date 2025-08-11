package merchanthandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchant"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatsMethodHandleApi struct {
	// client is the client for the merchant stats method service.
	client pb.MerchantStatsMethodServiceClient

	// router is the router for the merchant stats method handler.
	logger logger.LoggerInterface

	// mapper is the mapper for the merchant stats method response.
	mapper apimapper.MerchantStatsMethodResponseMapper

	// trace is the tracer for the merchant stats method handler.
	trace trace.Tracer

	// requestCounter is the counter for the merchant stats method handler.
	requestCounter *prometheus.CounterVec

	// requestDuration is the duration for the merchant stats method handler.
	requestDuration *prometheus.HistogramVec
}

// merchantStatsMethodHandleDeps is a struct that holds the necessary dependencies for the merchantStatsMethodHandleApi.
type merchantStatsMethodHandleDeps struct {
	// client is the client for the merchant stats method service.
	client pb.MerchantStatsMethodServiceClient

	// router is the router for the merchant stats method handler.
	router *echo.Echo

	// logger is the logger for the merchant stats method handler.
	logger logger.LoggerInterface

	// mapper is the mapper for the merchant stats method response.
	mapper apimapper.MerchantStatsMethodResponseMapper
}

// NewMerchantStatsMethodHandleApi initializes a new merchantStatsMethodHandleApi and sets up the routes for merchant stats method-related operations.
//
// This function registers various HTTP endpoints related to merchant statistics method management, including retrieval of monthly and yearly transaction methods.
// It also collects Prometheus metrics for tracking the number of requests and their durations, helping to monitor the performance and reliability of the handler.
// The routes are grouped under "/api/merchant-stats-method".
//
// Parameters:
// - params: A pointer to merchantStatsMethodHandleDeps containing the necessary dependencies such as router, logger, client, and mapper.
//
// Returns:
// - A pointer to a newly created merchantStatsMethodHandleApi.
func NewMerchantStatsMethodHandleApi(params *merchantStatsMethodHandleDeps) *merchantStatsMethodHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_stats_method_handler_requests_total",
			Help: "Total number of merchant stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_stats_method_handler_request_duration_seconds",
			Help:    "Duration of merchant stats amount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantHandler := &merchantStatsMethodHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("merchant-stats-amount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := params.router.Group("/api/merchant-stats-amount")

	routerMerchant.GET("/monthly-payment-methods", merchantHandler.FindMonthlyPaymentMethodsMerchant)
	routerMerchant.GET("/yearly-payment-methods", merchantHandler.FindYearlyPaymentMethodMerchant)

	routerMerchant.GET("/monthly-payment-methods-by-merchant", merchantHandler.FindMonthlyPaymentMethodByMerchants)
	routerMerchant.GET("/yearly-payment-methods-by-merchant", merchantHandler.FindYearlyPaymentMethodByMerchants)
	routerMerchant.GET("/monthly-payment-methods-by-apikey", merchantHandler.FindMonthlyPaymentMethodByApikeys)
	routerMerchant.GET("/yearly-payment-methods-by-apikey", merchantHandler.FindYearlyPaymentMethodByApikeys)

	return merchantHandler
}

// FindMonthlyPaymentMethodsMerchant godoc
// @Summary Find monthly payment methods for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchant-stats-amount/monthly-payment-methods [get]
func (h *merchantStatsMethodHandleApi) FindMonthlyPaymentMethodsMerchant(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodsMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.client.FindMonthlyPaymentMethodsMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve monthly payment methods", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyPaymentMethodsMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyPaymentMethods(res)

	logSuccess("monthly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodMerchant godoc.
// @Summary Find yearly payment methods for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchant-stats-amount/yearly-payment-methods [get]
func (h *merchantStatsMethodHandleApi) FindYearlyPaymentMethodMerchant(c echo.Context) error {
	const method = "FindYearlyPaymentMethodMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.client.FindYearlyPaymentMethodMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve yearly payment methods", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyPaymentMethodMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyPaymentMethods(res)

	logSuccess("yearly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodByMerchants godoc.
// @Summary Find monthly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchant-stats-amount/monthly-payment-methods-by-merchant [get]
func (h *merchantStatsMethodHandleApi) FindMonthlyPaymentMethodByMerchants(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantIDStr := c.QueryParam("merchant_id")

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil || merchantID <= 0 {
		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.client.FindMonthlyPaymentMethodByMerchants(ctx, req)

	if err != nil {
		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyPaymentMethodByMerchants(c)
	}

	so := h.mapper.ToApiResponseMonthlyPaymentMethods(res)

	logSuccess("monthly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodByMerchants godoc.
// @Summary Find yearly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchant-stats-amount/yearly-payment-methods-by-merchant [get]
func (h *merchantStatsMethodHandleApi) FindYearlyPaymentMethodByMerchants(c echo.Context) error {
	const method = "FindYearlyPaymentMethodByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantIDStr := c.QueryParam("merchant_id")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.client.FindYearlyPaymentMethodByMerchants(ctx, req)

	if err != nil {
		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyPaymentMethodByMerchants(c)
	}

	so := h.mapper.ToApiResponseYearlyPaymentMethods(res)

	logSuccess("yearly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodByApikeys godoc.
// @Summary Find monthly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchant-stats-amount/monthly-payment-methods-by-apikey [get]
func (h *merchantStatsMethodHandleApi) FindMonthlyPaymentMethodByApikeys(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	api_key := c.QueryParam("api_key")

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.client.FindMonthlyPaymentMethodByApikey(ctx, req)

	if err != nil {
		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyPaymentMethodsMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyPaymentMethods(res)

	logSuccess("monthly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodByApikeys godoc.
// @Summary Find yearly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchant-stats-amount/yearly-payment-methods-by-apikey [get]
func (h *merchantStatsMethodHandleApi) FindYearlyPaymentMethodByApikeys(c echo.Context) error {
	const method = "FindYearlyPaymentMethodByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	api_key := c.QueryParam("api_key")

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.client.FindYearlyPaymentMethodByApikey(ctx, req)

	if err != nil {
		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyPaymentMethodMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyPaymentMethods(res)

	logSuccess("yearly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantStatsMethodHandleApi) startTracingAndLogging(
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

func (s *merchantStatsMethodHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

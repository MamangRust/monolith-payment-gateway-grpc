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

// merchantStatsAmountHandleApi is a struct that handles merchant stats amount requests.
type merchantStatsAmountHandleApi struct {
	// client is the client for the merchant stats amount service.
	client pb.MerchantStatsAmountServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper is the mapper for the merchant stats amount response.
	mapper apimapper.MerchantStatsAmountResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP requests handled by this service.
	requestDuration *prometheus.HistogramVec
}

// merchantStatsAmountHandleDeps is a struct that holds the necessary dependencies for the merchantStatsAmountHandleApi.
type merchantStatsAmountHandleDeps struct {
	// client is the client for the merchant stats amount service.
	client pb.MerchantStatsAmountServiceClient

	// router is the router for the merchant stats amount service.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper is the mapper for the merchant stats amount response.
	mapper apimapper.MerchantStatsAmountResponseMapper
}

// NewMerchantStatsAmountHandleApi initializes a new merchantStatsAmountHandleApi and sets up the routes for merchant statistics amount-related operations.
//
// This function registers various HTTP endpoints related to merchant statistics amount management, including retrieval of monthly and yearly transaction amounts.
// It also collects Prometheus metrics for tracking the number of requests and their durations, helping to monitor the performance and reliability of the handler.
// The routes are grouped under "/api/merchant-stats-amount".
//
// Parameters:
// - params: A pointer to merchantStatsAmountHandleDeps containing the necessary dependencies such as router, logger, client, and mapper.
//
// Returns:
// - A pointer to a newly created merchantStatsAmountHandleApi.
func NewMerchantStatsAmountHandleApi(params *merchantStatsAmountHandleDeps) *merchantStatsAmountHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_stats_amount_handler_requests_total",
			Help: "Total number of merchant stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_stats_amount_handler_request_duration_seconds",
			Help:    "Duration of merchant stats amount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	merchantHandler := &merchantStatsAmountHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("merchant-stats-amount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := params.router.Group("/api/merchant-stats-amount")

	routerMerchant.GET("/monthly-amount", merchantHandler.FindMonthlyAmountMerchant)
	routerMerchant.GET("/yearly-amount", merchantHandler.FindYearlyAmountMerchant)

	routerMerchant.GET("/monthly-amount-by-merchant", merchantHandler.FindMonthlyAmountByMerchants)
	routerMerchant.GET("/yearly-amount-by-merchant", merchantHandler.FindYearlyAmountByMerchants)

	routerMerchant.GET("/monthly-amount-by-apikey", merchantHandler.FindMonthlyAmountByApikeys)
	routerMerchant.GET("/yearly-amount-by-apikey", merchantHandler.FindYearlyAmountByApikeys)

	return merchantHandler
}

// FindMonthlyAmountMerchant godoc
// @Summary Find monthly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchant-stats-amount/monthly-amount [get]
func (h *merchantStatsAmountHandleApi) FindMonthlyAmountMerchant(c echo.Context) error {
	const method = "FindMonthlyAmountMerchant"
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

	res, err := h.client.FindMonthlyAmountMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve monthly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("monthly transaction amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountMerchant godoc.
// @Summary Find yearly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchant-stats-amount/yearly-amount [get]
func (h *merchantStatsAmountHandleApi) FindYearlyAmountMerchant(c echo.Context) error {
	const method = "FindYearlyAmountMerchant"
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

	res, err := h.client.FindYearlyAmountMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve yearly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("yearly transaction amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountByMerchants godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchant-stats-amount/monthly-amount-by-merchant [get]
func (h *merchantStatsAmountHandleApi) FindMonthlyAmountByMerchants(c echo.Context) error {
	const method = "FindMonthlyAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantIDStr := c.QueryParam("merchant_id")

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil || merchantID <= 0 {
		logError("failed to find monthly amount by merchant", err, zap.Error(err))

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

	res, err := h.client.FindMonthlyAmountByMerchants(ctx, req)

	if err != nil {
		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyAmountByMerchants(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("monthly amount retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountByMerchants godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchant-stats-amount/yearly-amount-by-merchant [get]
func (h *merchantStatsAmountHandleApi) FindYearlyAmountByMerchants(c echo.Context) error {
	const method = "FindYearlyAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantIDStr := c.QueryParam("merchant_id")

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil || merchantID <= 0 {
		logError("failed to find yearly amount by merchant", err, zap.Error(err))

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

	res, err := h.client.FindYearlyAmountByMerchants(ctx, req)

	if err != nil {
		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyAmountByMerchants(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("yearly amount retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountByApikeys godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchant-stats-amount/monthly-amount-by-apikey [get]
func (h *merchantStatsAmountHandleApi) FindMonthlyAmountByApikeys(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
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

	res, err := h.client.FindMonthlyAmountByApikey(ctx, req)

	if err != nil {
		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("monthly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountByApikeys godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchant-stats-amount/yearly-amount-by-apikey [get]
func (h *merchantStatsAmountHandleApi) FindYearlyAmountByApikeys(c echo.Context) error {
	const method = "FindYearlyAmountByApikeys"
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

	res, err := h.client.FindYearlyAmountByApikey(ctx, req)

	if err != nil {
		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("yearly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantStatsAmountHandleApi) startTracingAndLogging(
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

func (s *merchantStatsAmountHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

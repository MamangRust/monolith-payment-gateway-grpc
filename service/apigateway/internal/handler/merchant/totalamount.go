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

type merchantStatsTotalAmountHandleApi struct {
	client pb.MerchantStatsTotalAmountServiceClient

	logger logger.LoggerInterface

	mapper apimapper.MerchantStatsTotalAmountResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type merchantStatsTotalAmountHandleDeps struct {
	client pb.MerchantStatsTotalAmountServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.MerchantStatsTotalAmountResponseMapper
}

func NewMerchantStatsTotalAmountHandleApi(params *merchantStatsTotalAmountHandleDeps) *merchantStatsTotalAmountHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_stats_totalamount_handler_requests_total",
			Help: "Total number of merchant stats totalamount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_stats_totalamount_handler_request_duration_seconds",
			Help:    "Duration of merchant stats totalamount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantHandler := &merchantStatsTotalAmountHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("merchant-stats-totalamount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := params.router.Group("/api/merchant-stats-totalamount")

	routerMerchant.GET("/monthly-totalamount", merchantHandler.FindMonthlyTotalAmountMerchant)
	routerMerchant.GET("/yearly-totalamount", merchantHandler.FindYearlyTotalAmountMerchant)

	routerMerchant.GET("/monthly-totalamount-by-merchant", merchantHandler.FindMonthlyTotalAmountByMerchants)
	routerMerchant.GET("/yearly-totalamount-by-merchant", merchantHandler.FindYearlyTotalAmountByMerchants)

	routerMerchant.GET("/monthly-totalamount-by-apikey", merchantHandler.FindMonthlyTotalAmountByApikeys)
	routerMerchant.GET("/yearly-totalamount-by-apikey", merchantHandler.FindYearlyTotalAmountByApikeys)

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
// @Router /api/merchant-stats-totalamount/monthly-total-amount [get]
func (h *merchantStatsTotalAmountHandleApi) FindMonthlyTotalAmountMerchant(c echo.Context) error {
	const method = "FindMonthlyTotalAmountMerchant"
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

	res, err := h.client.FindMonthlyTotalAmountMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve monthly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyTotalAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyTotalAmounts(res)

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
// @Router /api/merchant-stats-totalamount/yearly-total-amount [get]
func (h *merchantStatsTotalAmountHandleApi) FindYearlyTotalAmountMerchant(c echo.Context) error {
	const method = "FindYearlyTotalAmountMerchant"
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

	res, err := h.client.FindYearlyTotalAmountMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve yearly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyTotalAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyTotalAmounts(res)

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
// @Router /api/merchant-stats-totalamount/monthly-totalamount-by-merchant [get]
func (h *merchantStatsTotalAmountHandleApi) FindMonthlyTotalAmountByMerchants(c echo.Context) error {
	const method = "FindMonthlyTotalAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantIDStr := c.QueryParam("merchant_id")

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil || merchantID <= 0 {
		logError("failed to find monthly total amount by merchant", err, zap.Error(err))

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

	res, err := h.client.FindMonthlyTotalAmountByMerchants(ctx, req)

	if err != nil {
		logError("failed to find monthly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyTotalAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyTotalAmounts(res)

	logSuccess("monthly total amount retrieved successfully", zap.Bool("success", true))

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
// @Router /api/merchant-stats-totalamount/yearly-totalamount-by-merchant [get]
func (h *merchantStatsTotalAmountHandleApi) FindYearlyTotalAmountByMerchants(c echo.Context) error {
	const method = "FindYearlyTotalAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantIDStr := c.QueryParam("merchant_id")

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil || merchantID <= 0 {
		logError("failed to find yearly total amount by merchant", err, zap.Error(err))

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

	res, err := h.client.FindYearlyTotalAmountByMerchants(ctx, req)

	if err != nil {
		logError("failed to find yearly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyTotalAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyTotalAmounts(res)

	logSuccess("yearly total amount retrieved successfully", zap.Bool("success", true))

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
// @Router /api/merchant-stats-totalamount/monthly-totalamount-by-apikey [get]
func (h *merchantStatsTotalAmountHandleApi) FindMonthlyTotalAmountByApikeys(c echo.Context) error {
	const method = "FindMonthlyTotalAmountByApikeys"
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

	res, err := h.client.FindMonthlyTotalAmountByApikey(ctx, req)

	if err != nil {
		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyTotalAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseMonthlyTotalAmounts(res)

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
// @Router /api/merchant-stats-totalamount/yearly-totalamount-by-apikey [get]
func (h *merchantStatsTotalAmountHandleApi) FindYearlyTotalAmountByApikeys(c echo.Context) error {
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

	res, err := h.client.FindYearlyTotalAmountByApikey(ctx, req)

	if err != nil {
		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyTotalAmountMerchant(c)
	}

	so := h.mapper.ToApiResponseYearlyTotalAmounts(res)

	logSuccess("yearly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantStatsTotalAmountHandleApi) startTracingAndLogging(
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

func (s *merchantStatsTotalAmountHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

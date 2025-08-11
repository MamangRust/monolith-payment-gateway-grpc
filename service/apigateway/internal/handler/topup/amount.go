package topuphandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/topup"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStatsAmountHandleApi struct {
	client pb.TopupStatsAmountServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TopupStatsAmountResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type topupStatsAmountHandleDeps struct {
	client pb.TopupStatsAmountServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TopupStatsAmountResponseMapper
}

func NewTopupStatsAmountHandleApi(params *topupStatsAmountHandleDeps) *topupStatsAmountHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_amount_handler_requests_total",
			Help: "Total number of topup stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_amount_handler_request_duration_seconds",
			Help:    "Duration of topup stats amount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	topupHandler := &topupStatsAmountHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("topup-stats-amount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTopup := params.router.Group("/api/topup-stats-amount")

	routerTopup.GET("/monthly-amounts", topupHandler.FindMonthlyTopupAmounts)
	routerTopup.GET("/yearly-amounts", topupHandler.FindYearlyTopupAmounts)

	routerTopup.GET("/monthly-amounts-by-card", topupHandler.FindMonthlyTopupAmountsByCardNumber)
	routerTopup.GET("/yearly-amounts-by-card", topupHandler.FindYearlyTopupAmountsByCardNumber)

	return topupHandler
}

// FindMonthlyTopupAmounts retrieves the monthly top-up amounts for a specific year.
// @Summary Get monthly top-up amounts
// @Tags Topup Amount
// @Security Bearer
// @Description Retrieve the monthly top-up amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount "Monthly top-up amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up amounts"
// @Router /api/topup/monthly-amounts [get]
func (h *topupStatsAmountHandleApi) FindMonthlyTopupAmounts(c echo.Context) error {
	const method = "FindMonthlyTopupAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTopupAmounts(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly top-up amounts", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupAmounts(c)
	}

	so := h.mapper.ToApiResponseTopupMonthAmount(res)

	logSuccess("success find monthly top-up amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmounts retrieves the yearly top-up amounts for a specific year.
// @Summary Get yearly top-up amounts
// @Tags Topup Amount
// @Security Bearer
// @Description Retrieve the yearly top-up amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount "Yearly top-up amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up amounts"
// @Router /api/topup-stats-amount/yearly-amounts [get]
func (h *topupStatsAmountHandleApi) FindYearlyTopupAmounts(c echo.Context) error {
	const method = "FindYearlyTopupAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupAmounts(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly top-up amounts", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupAmounts(c)
	}

	so := h.mapper.ToApiResponseTopupYearAmount(res)

	logSuccess("success find yearly top-up amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupAmountsByCardNumber retrieves the monthly top-up amounts for a specific card number and year.
// @Summary Get monthly top-up amounts by card number
// @Tags Topup Amount
// @Security Bearer
// @Description Retrieve the monthly top-up amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount "Monthly top-up amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up amounts by card number"
// @Router /api/topup-stats-amount/monthly-amounts-by-card [get]
func (h *topupStatsAmountHandleApi) FindMonthlyTopupAmountsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTopupAmountsByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTopupAmountsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly top-up amounts by card number", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupAmountsByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTopupMonthAmount(res)

	logSuccess("success find monthly top-up amounts by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmountsByCardNumber retrieves the yearly top-up amounts for a specific card number and year.
// @Summary Get yearly top-up amounts by card number
// @Tags Topup Amount
// @Security Bearer
// @Description Retrieve the yearly top-up amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount "Yearly top-up amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up amounts by card number"
// @Router /api/topup-stats-amount/yearly-amounts-by-card [get]
func (h *topupStatsAmountHandleApi) FindYearlyTopupAmountsByCardNumber(c echo.Context) error {
	const method = "FindAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupAmountsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly top-up amounts by card number", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupAmountsByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTopupYearAmount(res)

	logSuccess("success find yearly top-up amounts by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *topupStatsAmountHandleApi) startTracingAndLogging(
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

func (s *topupStatsAmountHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

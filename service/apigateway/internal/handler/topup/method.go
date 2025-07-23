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

type topupStatsMethodHandleApi struct {
	client pb.TopupStatsMethodServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TopupStatsMethodResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type topupStatsMethodHandleDeps struct {
	client pb.TopupStatsMethodServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TopupStatsMethodResponseMapper
}

func NewTopupStatsMethodHandleApi(params *topupStatsMethodHandleDeps) *topupStatsMethodHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_method_handler_requests_total",
			Help: "Total number of topup stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_method_handler_request_duration_seconds",
			Help:    "Duration of topup stats method requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	topupHandler := &topupStatsMethodHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("topup-stats-method-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTopup := params.router.Group("/api/topup-stats-method")

	routerTopup.GET("/monthly-methods", topupHandler.FindMonthlyTopupMethods)
	routerTopup.GET("/yearly-methods", topupHandler.FindYearlyTopupMethods)

	routerTopup.GET("/monthly-methods-by-card", topupHandler.FindMonthlyTopupMethodsByCardNumber)
	routerTopup.GET("/yearly-methods-by-card", topupHandler.FindYearlyTopupMethodsByCardNumber)

	return topupHandler
}

// FindMonthlyTopupMethods retrieves the monthly top-up methods for a specific year.
// @Summary Get monthly top-up methods
// @Tags Topup Method
// @Security Bearer
// @Description Retrieve the monthly top-up methods for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod "Monthly top-up methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up methods"
// @Router /api/topup-stats-method/monthly-methods [get]
func (h *topupStatsMethodHandleApi) FindMonthlyTopupMethods(c echo.Context) error {
	const method = "FindMonthlyTopupMethods"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTopupMethods(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly top-up methods", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupMethods(c)
	}

	so := h.mapper.ToApiResponseTopupMonthMethod(res)

	logSuccess("success find monthly top-up methods", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupMethods retrieves the yearly top-up methods for a specific year.
// @Summary Get yearly top-up methods
// @Tags Topup Method
// @Security Bearer
// @Description Retrieve the yearly top-up methods for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod "Yearly top-up methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up methods"
// @Router /api/topup-stats-method/yearly-methods [get]
func (h *topupStatsMethodHandleApi) FindYearlyTopupMethods(c echo.Context) error {
	const method = "FindYearlyTopupMethods"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupMethods(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly top-up methods", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupMethods(c)
	}

	so := h.mapper.ToApiResponseTopupYearMethod(res)

	logSuccess("success find yearly top-up methods", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupMethodsByCardNumber retrieves the monthly top-up methods for a specific card number and year.
// @Summary Get monthly top-up methods by card number
// @Tags Topup Method
// @Security Bearer
// @Description Retrieve the monthly top-up methods for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod "Monthly top-up methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up methods by card number"
// @Router /api/topup-stats-method/monthly-methods-by-card [get]
func (h *topupStatsMethodHandleApi) FindMonthlyTopupMethodsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTopupMethodsByCardNumber"
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

	res, err := h.client.FindMonthlyTopupMethodsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly top-up methods by card number", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupMethodsByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTopupMonthMethod(res)

	logSuccess("success find monthly top-up methods by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupMethodsByCardNumber retrieves the yearly top-up methods for a specific card number and year.
// @Summary Get yearly top-up methods by card number
// @Tags Topup Method
// @Security Bearer
// @Description Retrieve the yearly top-up methods for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod "Yearly top-up methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up methods by card number"
// @Router /api/topup-stats-method/yearly-methods-by-card [get]
func (h *topupStatsMethodHandleApi) FindYearlyTopupMethodsByCardNumber(c echo.Context) error {
	const method = "FindYearlyTopupMethodsByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupMethodsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly top-up methods by card number", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupMethodsByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTopupYearMethod(res)

	logSuccess("success find yearly top-up methods by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *topupStatsMethodHandleApi) startTracingAndLogging(
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

func (s *topupStatsMethodHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

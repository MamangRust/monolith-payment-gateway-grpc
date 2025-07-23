package cardhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
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

type cardStatsBalanceHandleApi struct {
	// card is the gRPC client used to interact with the CardService.
	card pb.CardStatsBalanceServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsBalanceResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// cardStatsBalanceHandleApiDeps contains the necessary dependencies for the cardStatsBalanceHandleApi to function.
type cardStatsBalanceHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardService.
	client pb.CardStatsBalanceServiceClient

	// router is the Echo HTTP router used to register routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsBalanceResponseMapper
}

// NewCardStatsBalanceHandleApi initializes a new cardStatsBalanceHandleApi and sets up the routes for card stats balance-related operations.
//
// This function registers various HTTP endpoints related to card stats balance management, including retrieval of monthly and yearly balances.
// It also tracks metrics like request count and duration using Prometheus metrics. The routes are grouped under "/api/card-stats-balance".
//
// Parameters:
// - params: A pointer to cardStatsBalanceHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardStatsBalanceHandleApi.
func NewCardStatsBalanceHandleApi(
	params *cardStatsBalanceHandleApiDeps,
) *cardStatsBalanceHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_stats_balance_handler_requests_total",
			Help: "Total number of Card Stats Balance requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_stats_balance_handler_request_duration_seconds",
			Help:    "Duration of Card Stats Balance requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	cardStatsBalanceHandler := &cardStatsBalanceHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-stats-balance-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-stats-balance")

	routerCard.GET("/monthly-balance", cardStatsBalanceHandler.FindMonthlyBalance)
	routerCard.GET("/yearly-balance", cardStatsBalanceHandler.FindYearlyBalance)

	routerCard.GET("/monthly-balance-by-card", cardStatsBalanceHandler.FindMonthlyBalanceByCardNumber)
	routerCard.GET("/yearly-balance-by-card", cardStatsBalanceHandler.FindYearlyBalanceByCardNumber)

	return cardStatsBalanceHandler
}

// FindMonthlyBalance godoc
// @Summary Get monthly balance data
// @Description Retrieve monthly balance data for a specific year
// @Tags Card-Stats-Balance
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-balance [get]
func (h *cardStatsBalanceHandleApi) FindMonthlyBalance(c echo.Context) error {
	const method = "FindMonthlyBalance"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearBalance{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyBalance(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly balance data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyBalance(c)
	}

	so := h.mapper.ToApiResponseMonthlyBalances(res)

	logSuccess("Success retrieve monthly balance data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyBalance godoc
// @Summary Get yearly balance data
// @Description Retrieve yearly balance data for a specific year
// @Tags Card-Stats-Balance
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-balance [get]
func (h *cardStatsBalanceHandleApi) FindYearlyBalance(c echo.Context) error {
	const method = "FindYearlyBalance"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearBalance{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyBalance(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly balance data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyBalance(c)
	}

	so := h.mapper.ToApiResponseYearlyBalances(res)

	logSuccess("Success retrieve yearly balance data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyBalanceByCardNumber godoc
// @Summary Get monthly balance data by card number
// @Description Retrieve monthly balance data for a specific year and card number
// @Tags Card-Stats-Balance
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-balance-by-card [get]
func (h *cardStatsBalanceHandleApi) FindMonthlyBalanceByCardNumber(c echo.Context) error {
	const method = "FindMonthlyBalanceByCardNumber"

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

	req := &pb.FindYearBalanceCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyBalanceByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly balance data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyBalanceByCard(c)
	}

	so := h.mapper.ToApiResponseMonthlyBalances(res)

	logSuccess("Success retrieve monthly balance data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyBalanceByCardNumber godoc
// @Summary Get yearly balance data by card number
// @Description Retrieve yearly balance data for a specific year and card number
// @Tags Card-Stats-Balance
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-balance-by-card [get]
func (h *cardStatsBalanceHandleApi) FindYearlyBalanceByCardNumber(c echo.Context) error {
	const method = "FindYearlyBalanceByCardNumber"

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

	req := &pb.FindYearBalanceCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyBalanceByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly balance data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyBalanceByCard(c)
	}

	so := h.mapper.ToApiResponseYearlyBalances(res)

	logSuccess("Success retrieve yearly balance data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// startTracingAndLogging starts a tracing span and returns functions to log the outcome of the call.
// The returned functions are logSuccess and logError, which log the outcome of the call to the trace span.
// The returned end function records the metrics and ends the trace span.
func (s *cardStatsBalanceHandleApi) startTracingAndLogging(
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

// recordMetrics records a Prometheus metric for the given method and status.
// It increments a counter and records the duration since the provided start time.
func (s *cardStatsBalanceHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

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

type cardStatsWithdrawHandleApi struct {
	// card is the gRPC client used to interact with the CardService.
	card pb.CardStatsWithdrawServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsAmountResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// cardStatsWithdrawHandleApiDeps is a struct that holds the necessary dependencies for the cardStatsWithdrawHandleApi.
type cardStatsWithdrawHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardService.
	client pb.CardStatsWithdrawServiceClient

	// router is the Echo router used to register HTTP routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsAmountResponseMapper
}

// NewCardStatsWithdrawHandleApi initializes a new cardStatsWithdrawHandleApi and sets up the routes for card stats withdraw-related operations.
//
// This function registers various HTTP endpoints related to card stats withdraw management, including retrieval of monthly and yearly withdraw amounts.
// It also tracks metrics like request count and duration using Prometheus metrics. The routes are grouped under "/api/card-stats-withdraw".
//
// Parameters:
// - params: A pointer to cardStatsWithdrawHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardStatsWithdrawHandleApi.
func NewCardStatsWithdrawHandleApi(
	params *cardStatsWithdrawHandleApiDeps,
) *cardStatsWithdrawHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_stats_withdraw_handler_requests_total",
			Help: "Total number of Card Stats Withdraw requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_stats_withdraw_handler_request_duration_seconds",
			Help:    "Duration of Card Stats Withdraw requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	cardStatsWithdrawHandler := &cardStatsWithdrawHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-stats-withdraw-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-stats-withdraw")

	routerCard.GET("/monthly-withdraw-amount", cardStatsWithdrawHandler.FindMonthlyWithdrawAmount)
	routerCard.GET("/yearly-withdraw-amount", cardStatsWithdrawHandler.FindYearlyWithdrawAmount)

	routerCard.GET("/monthly-withdraw-amount-by-card", cardStatsWithdrawHandler.FindMonthlyWithdrawAmountByCardNumber)
	routerCard.GET("/yearly-withdraw-amount-by-card", cardStatsWithdrawHandler.FindYearlyWithdrawAmountByCardNumber)

	return cardStatsWithdrawHandler
}

// FindMonthlyWithdrawAmount
// godoc
// @Summary Get monthly withdraw amount data
// @Description Retrieve monthly withdraw amount data for a specific year
// @Tags Card-Stats-Withdraw
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-withdraw/monthly-withdraw-amount [get]
func (h *cardStatsWithdrawHandleApi) FindMonthlyWithdrawAmount(c echo.Context) error {
	const method = "FindMonthlyWithdrawAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyWithdrawAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly withdraw amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyWithdrawAmount(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly withdraw amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawAmount godoc
// @Summary Get yearly withdraw amount data
// @Description Retrieve yearly withdraw amount data for a specific year
// @Tags Card-Stats-Withdraw
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-withdraw/yearly-withdraw-amount [get]
func (h *cardStatsWithdrawHandleApi) FindYearlyWithdrawAmount(c echo.Context) error {
	const method = "FindYearlyWithdrawAmount"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyWithdrawAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly withdraw amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyWithdrawAmount(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly withdraw amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawAmountByCardNumber godoc
// @Summary Get monthly withdraw amount data by card number
// @Description Retrieve monthly withdraw amount data for a specific year and card number
// @Tags Card-Stats-Withdraw
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-withdraw/monthly-withdraw-amount-by-card [get]
func (h *cardStatsWithdrawHandleApi) FindMonthlyWithdrawAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyWithdrawAmountByCardNumber"

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

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyWithdrawAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly withdraw amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyWithdrawAmountByCard(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly withdraw amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawAmountByCardNumber godoc
// @Summary Get yearly withdraw amount data by card number
// @Description Retrieve yearly withdraw amount data for a specific year and card number
// @Tags Card-Stats-Withdraw
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-withdraw/yearly-withdraw-amount-by-card [get]
func (h *cardStatsWithdrawHandleApi) FindYearlyWithdrawAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyWithdrawAmountByCardNumber"

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

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyWithdrawAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly withdraw amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyWithdrawAmountByCard(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly withdraw amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// startTracingAndLogging starts a tracing span and returns functions to log the outcome of the call.
// The returned functions are logSuccess and logError, which log the outcome of the call to the trace span.
// The returned end function records the metrics and ends the trace span.
func (s *cardStatsWithdrawHandleApi) startTracingAndLogging(
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

// recordMetrics records Prometheus metrics for the specified method and status.
// It increments the request counter for the provided method and status,
// and observes the request duration by calculating the time elapsed since the provided start time.
func (s *cardStatsWithdrawHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

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

type cardStatsTopupHandleApi struct {
	// card is the gRPC client used to interact with the CardService.
	card pb.CardStatsTopupServiceClient

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

// cardStatsTopupHandleApiDeps contains the necessary dependencies for initializing a cardStatsTopupHandleApi.
type cardStatsTopupHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardService.
	client pb.CardStatsTopupServiceClient

	// router is the Echo router used to register HTTP routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsAmountResponseMapper
}

// NewCardStatsTopupHandleApi initializes a new cardStatsTopupHandleApi and sets up the routes for
// card stats topup-related operations.
//
// This function registers various HTTP endpoints related to card stats topup management, including
// retrieval of monthly and yearly topup amounts. It also tracks metrics like request count and
// duration using Prometheus metrics. The routes are grouped under "/api/card-stats-topup".
//
// Parameters:
// - params: A pointer to cardStatsTopupHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardStatsTopupHandleApi.
func NewCardStatsTopupHandleApi(
	params *cardStatsTopupHandleApiDeps,
) *cardStatsTopupHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_stats_topup_handler_requests_total",
			Help: "Total number of Card Stats Topup requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_stats_topup_handler_request_duration_seconds",
			Help:    "Duration of Card Stats Topup requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	cardStatsTopupHandler := &cardStatsTopupHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-stats-topup-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-stats-topup")

	routerCard.GET("/monthly-topup-amount", cardStatsTopupHandler.FindMonthlyTopupAmount)
	routerCard.GET("/yearly-topup-amount", cardStatsTopupHandler.FindYearlyTopupAmount)

	routerCard.GET("/monthly-topup-amount-by-card", cardStatsTopupHandler.FindMonthlyTopupAmountByCardNumber)
	routerCard.GET("/yearly-topup-amount-by-card", cardStatsTopupHandler.FindYearlyTopupAmountByCardNumber)

	return cardStatsTopupHandler
}

// FindMonthlyTopupAmount godoc
// @Summary Get monthly topup amount data
// @Description Retrieve monthly topup amount data for a specific year
// @Tags Card-Stats-Topup
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-topup-amount [get]
func (h *cardStatsTopupHandleApi) FindMonthlyTopupAmount(c echo.Context) error {
	const method = "FindMonthlyTopupAmount"

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

	res, err := h.card.FindMonthlyTopupAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly topup amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTopupAmount(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly topup amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmount godoc
// @Summary Get yearly topup amount data
// @Description Retrieve yearly topup amount data for a specific year
// @Tags Card-Stats-Topup
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/topup/yearly-topup-amount [get]
func (h *cardStatsTopupHandleApi) FindYearlyTopupAmount(c echo.Context) error {
	const method = "FindYearlyTopupAmount"

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

	res, err := h.card.FindYearlyTopupAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly topup amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTopupAmount(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly topup amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupAmountByCardNumber godoc
// @Summary Get monthly topup amount data by card number
// @Description Retrieve monthly topup amount data for a specific year and card number
// @Tags Card-Stats-Topup
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-topup-amount-by-card [get]
func (h *cardStatsTopupHandleApi) FindMonthlyTopupAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTopupAmountByCardNumber"

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

	res, err := h.card.FindMonthlyTopupAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly topup amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTopupAmountByCard(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly topup amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmountByCardNumber godoc
// @Summary Get yearly topup amount data by card number
// @Description Retrieve yearly topup amount data for a specific year and card number
// @Tags Card-Stats-Topup
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-topup-amount-by-card [get]
func (h *cardStatsTopupHandleApi) FindYearlyTopupAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTopupAmountByCardNumber"

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

	res, err := h.card.FindYearlyTopupAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly topup amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTopupAmountByCard(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly topup amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// startTracingAndLogging starts a tracing span and returns functions to log the outcome of the call.
// The returned functions are logSuccess and logError, which log the outcome of the call to the trace span.
// The returned end function records the metrics and ends the trace span.
func (s *cardStatsTopupHandleApi) startTracingAndLogging(
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

// recordMetrics records Prometheus metrics for the given method and status.
// It increments the request counter and observes the request duration
// for the given method and status, using the provided start time.
func (s *cardStatsTopupHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

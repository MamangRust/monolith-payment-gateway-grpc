package cardhandler

import (
	"context"
	"net/http"
	"strconv"
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

type cardStatsTransactionHandleApi struct {
	// card is the gRPC client used to interact with the CardStatsTransactonServiceClient.
	card pb.CardStatsTransactonServiceClient

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

// cardStatsTransactionHandleApiDeps is a struct that holds the dependencies required to initialize a cardStatsTransactionHandleApi.
type cardStatsTransactionHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardStatsTransactonServiceClient.
	client pb.CardStatsTransactonServiceClient

	// router is the Echo router used to register HTTP routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsAmountResponseMapper
}

// NewCardStatsTransactionHandleApi initializes a new cardStatsTransactionHandleApi and sets up the routes for card stats transaction-related operations.
//
// This function registers various HTTP endpoints related to card stats transaction management, including retrieval of monthly and yearly transaction amounts.
// It also tracks metrics like request count and duration using Prometheus metrics. The routes are grouped under "/api/card-stats-transaction".
//
// Parameters:
// - params: A pointer to cardStatsTransactionHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardStatsTransactionHandleApi.
func NewCardStatsTransactionHandleApi(
	params *cardStatsTransactionHandleApiDeps,
) *cardStatsTransactionHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_stats_transaction_handler_requests_total",
			Help: "Total number of Card Stats Transaction requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_stats_transaction_handler_request_duration_seconds",
			Help:    "Duration of Card Stats Transaction requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	cardStatsTransactionHandler := &cardStatsTransactionHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-stats-transaction-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-stats-transaction")

	routerCard.GET("/monthly-transaction-amount", cardStatsTransactionHandler.FindMonthlyTransactionAmount)
	routerCard.GET("/yearly-transaction-amount", cardStatsTransactionHandler.FindYearlyTransactionAmount)

	routerCard.GET("/monthly-transaction-amount-by-card", cardStatsTransactionHandler.FindMonthlyTransactionAmountByCardNumber)
	routerCard.GET("/yearly-transaction-amount-by-card", cardStatsTransactionHandler.FindYearlyTransactionAmountByCardNumber)

	return cardStatsTransactionHandler
}

// FindMonthlyTransactionAmount godoc
// @Summary Get monthly transaction amount data
// @Description Retrieve monthly transaction amount data for a specific year
// @Tags Card-Stats-Transaction
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transaction/monthly-transaction-amount [get]
func (h *cardStatsTransactionHandleApi) FindMonthlyTransactionAmount(c echo.Context) error {
	const method = "FindMonthlyTransactionAmount"

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

	res, err := h.card.FindMonthlyTransactionAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly transaction amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTransactionAmount(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transaction amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionAmount godoc
// @Summary Get yearly transaction amount data
// @Description Retrieve yearly transaction amount data for a specific year
// @Tags Card-Stats-Transaction
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transaction/yearly-transaction-amount [get]
func (h *cardStatsTransactionHandleApi) FindYearlyTransactionAmount(c echo.Context) error {
	const method = "FindYearlyTransactionAmount"
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

	res, err := h.card.FindYearlyTransactionAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly transaction amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTransactionAmount(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transaction amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionAmountByCardNumber godoc
// @Summary Get monthly transaction amount data by card number
// @Description Retrieve monthly transaction amount data for a specific year and card number
// @Tags Card-Stats-Transaction
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transaction/monthly-transaction-amount-by-card [get]
func (h *cardStatsTransactionHandleApi) FindMonthlyTransactionAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransactionAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		logError("Invalid year", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidYear(c)
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransactionAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly transaction amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTransactionAmountByCard(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transaction amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionAmountByCardNumber godoc
// @Summary Get yearly transaction amount data by card number
// @Description Retrieve yearly transaction amount data for a specific year and card number
// @Tags Card-Stats-Transaction
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transaction/yearly-transaction-amount-by-card [get]
func (h *cardStatsTransactionHandleApi) FindYearlyTransactionAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransactionAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		logError("Invalid year", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidYear(c)
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransactionAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly transaction amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTransactionAmountByCard(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transaction amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// startTracingAndLogging initializes tracing and logging for a given method within the context.
// It returns three functions: `end` to conclude the tracing and log metrics,
// `logSuccess` to log successful events, and `logError` to log errors.
//
// Parameters:
// - ctx: The context in which the tracing and logging occur.
// - method: The name of the method being traced and logged.
// - attrs: Optional key-value attributes to be set on the trace span.
func (s *cardStatsTransactionHandleApi) startTracingAndLogging(
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
func (s *cardStatsTransactionHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

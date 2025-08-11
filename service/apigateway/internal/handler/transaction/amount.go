package transactionhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transaction"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatsAmountHandleApi struct {
	client pb.TransactionsStatsAmountServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransactionStatsAmountResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transactionStatsAmountHandleDeps struct {
	client pb.TransactionsStatsAmountServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransactionStatsAmountResponseMapper
}

func NewTransactionStatsAmountHandleApi(params *transactionStatsAmountHandleDeps) *transactionStatsAmountHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_amount_handler_requests_total",
			Help: "Total number of transaction stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_amount_handler_request_duration_seconds",
			Help:    "Duration of transaction stats amount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	transactionStatsAmountHandleApi := &transactionStatsAmountHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transaction-stats-amount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransaction := params.router.Group("/api/transaction-stats-amount")

	routerTransaction.GET("/monthly-amounts-by-card", transactionStatsAmountHandleApi.FindMonthlyAmountsByCardNumber)
	routerTransaction.GET("/yearly-amounts-by-card", transactionStatsAmountHandleApi.FindYearlyAmountsByCardNumber)
	routerTransaction.GET("/monthly-amounts", transactionStatsAmountHandleApi.FindMonthlyAmounts)
	routerTransaction.GET("/yearly-amounts", transactionStatsAmountHandleApi.FindYearlyAmounts)

	return transactionStatsAmountHandleApi
}

// FindMonthlyAmounts retrieves the monthly transaction amounts for a specific year.
// @Summary Get monthly transaction amounts
// @Tags Transaction Stats Amount
// @Security Bearer
// @Description Retrieve the monthly transaction amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/transaction-stats-amount/monthly-amounts [get]
func (h *transactionStatsAmountHandleApi) FindMonthlyAmounts(c echo.Context) error {
	const method = "FindMonthlyAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyAmounts(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly amounts", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyAmounts(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthAmount(res)

	logSuccess("success retrieve monthly amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmounts retrieves the yearly transaction amounts for a specific year.
// @Summary Get yearly transaction amounts
// @Tags Transaction Stats Amount
// @Security Bearer
// @Description Retrieve the yearly transaction amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/transaction-stats-amount/yearly-amounts [get]
func (h *transactionStatsAmountHandleApi) FindYearlyAmounts(c echo.Context) error {
	const method = "FindYearlyAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyAmounts(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly amounts", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyAmounts(c)
	}

	so := h.mapper.ToApiResponseTransactionYearAmount(res)

	logSuccess("success retrieve yearly amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountsByCardNumber retrieves the monthly transaction amounts for a specific card number and year.
// @Summary Get monthly transaction amounts by card number
// @Tags Transaction Stats Amount
// @Security Bearer
// @Description Retrieve the monthly transaction amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount "Monthly transaction amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts by card number"
// @Router /api/transaction-stats-amount/monthly-amounts-by-card [get]
func (h *transactionStatsAmountHandleApi) FindMonthlyAmountsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyAmountsByCardNumber"
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

	res, err := h.client.FindMonthlyAmountsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly amounts by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyAmountsByCard(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthAmount(res)

	logSuccess("success retrieve monthly amounts by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountsByCardNumber retrieves the yearly transaction amounts for a specific card number and year.
// @Summary Get yearly transaction amounts by card number
// @Tags Transaction Stats Amount
// @Security Bearer
// @Description Retrieve the yearly transaction amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount "Yearly transaction amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts by card number"
// @Router /api/transaction-stats-amount/yearly-amounts-by-card [get]
func (h *transactionStatsAmountHandleApi) FindYearlyAmountsByCardNumber(c echo.Context) error {
	const method = "FindYearlyAmountsByCardNumber"
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

	res, err := h.client.FindYearlyAmountsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly amounts by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyAmountsByCard(c)
	}

	so := h.mapper.ToApiResponseTransactionYearAmount(res)

	logSuccess("success retrieve yearly amounts by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionStatsAmountHandleApi) startTracingAndLogging(
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

func (s *transactionStatsAmountHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

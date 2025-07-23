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

type transactionStatsMethodHandleApi struct {
	client pb.TransactionStatsMethodServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransactionStatsMethodResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transactionStatsMethodHandleDeps struct {
	client pb.TransactionStatsMethodServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransactionStatsMethodResponseMapper
}

func NewTransactionStatsMethodHandleApi(params *transactionStatsMethodHandleDeps) *transactionStatsMethodHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_method_handler_requests_total",
			Help: "Total number of transaction stats method requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_method_handler_request_duration_seconds",
			Help:    "Duration of transaction stats method requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	transactionStatsMethodHandleApi := &transactionStatsMethodHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transaction-stats-method-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransaction := params.router.Group("/api/transaction-stats-method")

	routerTransaction.GET("/monthly-methods", transactionStatsMethodHandleApi.FindMonthlyPaymentMethods)
	routerTransaction.GET("/yearly-methods", transactionStatsMethodHandleApi.FindYearlyPaymentMethods)
	routerTransaction.GET("/monthly-methods-by-card", transactionStatsMethodHandleApi.FindMonthlyPaymentMethodsByCardNumber)
	routerTransaction.GET("/yearly-methods-by-card", transactionStatsMethodHandleApi.FindYearlyPaymentMethodsByCardNumber)

	return transactionStatsMethodHandleApi
}

// FindMonthlyPaymentMethods retrieves the monthly payment methods for transactions.
// @Summary Get monthly payment methods
// @Tags Transaction Stats Method
// @Security Bearer
// @Description Retrieve the monthly payment methods for transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/transaction-stats-method/monthly-payment-methods [get]
func (h *transactionStatsMethodHandleApi) FindMonthlyPaymentMethods(c echo.Context) error {
	const method = "FindMonthlyPaymentMethods"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyPaymentMethods(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("failed to retrieve monthly payment methods", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyMethods(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthMethod(res)

	logSuccess("success retrieve monthly payment methods", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethods retrieves the yearly payment methods for transactions.
// @Summary Get yearly payment methods
// @Tags Transaction Stats Method
// @Security Bearer
// @Description Retrieve the yearly payment methods for transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/transaction-stats-method/yearly-payment-methods [get]
func (h *transactionStatsMethodHandleApi) FindYearlyPaymentMethods(c echo.Context) error {
	const method = "FindYearlyPaymentMethods"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyPaymentMethods(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly payment methods", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyMethods(c)
	}

	so := h.mapper.ToApiResponseTransactionYearMethod(res)

	logSuccess("success retrieve yearly payment methods", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodsByCardNumber retrieves the monthly payment methods for transactions by card number and year.
// @Summary Get monthly payment methods by card number
// @Tags Transaction Stats Method
// @Security Bearer
// @Description Retrieve the monthly payment methods for transactions by card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod "Monthly payment methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods by card number"
// @Router /api/transaction-stats-method/monthly-payment-methods-by-card [get]
func (h *transactionStatsMethodHandleApi) FindMonthlyPaymentMethodsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodsByCardNumber"
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

	res, err := h.client.FindMonthlyPaymentMethodsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("failed to retrieve monthly payment methods by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyMethodsByCard(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthMethod(res)

	logSuccess("success retrieve monthly payment methods by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodsByCardNumber retrieves the yearly payment methods for transactions by card number and year.
// @Summary Get yearly payment methods by card number
// @Tags Transaction Stats Method
// @Security Bearer
// @Description Retrieve the yearly payment methods for transactions by card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod "Yearly payment methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods by card number"
// @Router /api/transaction-stats-method/yearly-payment-methods-by-card [get]
func (h *transactionStatsMethodHandleApi) FindYearlyPaymentMethodsByCardNumber(c echo.Context) error {
	const method = "FindYearlyPaymentMethodsByCardNumber"
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

	res, err := h.client.FindYearlyPaymentMethodsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly payment methods by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyMethodsByCard(c)
	}

	so := h.mapper.ToApiResponseTransactionYearMethod(res)

	logSuccess("success retrieve yearly payment methods by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionStatsMethodHandleApi) startTracingAndLogging(
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

func (s *transactionStatsMethodHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

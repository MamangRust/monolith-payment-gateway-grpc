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

type transactionStatsStatusHandleApi struct {
	client pb.TransactionStatsStatusServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransactionStatsStatusResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transactionStatsStatusHandleDeps struct {
	client pb.TransactionStatsStatusServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransactionStatsStatusResponseMapper
}

func NewTransactionStatsStatusHandleApi(params *transactionStatsStatusHandleDeps) *transactionStatsStatusHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_status_handler_requests_total",
			Help: "Total number of transaction stats status requests",
		},
		[]string{"status", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_status_handler_request_duration_seconds",
			Help:    "Duration of transaction stats status requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "status"},
	)

	transactionStatsStatusHandleApi := &transactionStatsStatusHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transaction-stats-status-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransaction := params.router.Group("/api/transaction-stats-status")

	routerTransaction.GET("/monthly-success", transactionStatsStatusHandleApi.FindMonthlyTransactionStatusSuccess)
	routerTransaction.GET("/yearly-success", transactionStatsStatusHandleApi.FindYearlyTransactionStatusSuccess)
	routerTransaction.GET("/monthly-failed", transactionStatsStatusHandleApi.FindMonthlyTransactionStatusFailed)
	routerTransaction.GET("/yearly-failed", transactionStatsStatusHandleApi.FindYearlyTransactionStatusFailed)

	routerTransaction.GET("/monthly-success-by-card", transactionStatsStatusHandleApi.FindMonthlyTransactionStatusSuccessByCardNumber)
	routerTransaction.GET("/yearly-success-by-card", transactionStatsStatusHandleApi.FindYearlyTransactionStatusSuccessByCardNumber)
	routerTransaction.GET("/monthly-failed-by-card", transactionStatsStatusHandleApi.FindMonthlyTransactionStatusFailedByCardNumber)
	routerTransaction.GET("/yearly-failed-by-card", transactionStatsStatusHandleApi.FindYearlyTransactionStatusFailedByCardNumber)

	return transactionStatsStatusHandleApi
}

// FindMonthlyTransactionStatusSuccess retrieves the monthly transaction status for successful transactions.
// @Summary Get monthly transaction status for successful transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the monthly transaction status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess "Monthly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for successful transactions"
// @Router /api/transaction-stats-status/monthly-success [get]
func (h *transactionStatsStatusHandleApi) FindMonthlyTransactionStatusSuccess(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransactionStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlySuccess(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthStatusSuccess(res)

	logSuccess("success retrieve monthly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
// @Summary Get yearly transaction status for successful transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the yearly transaction status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusSuccess "Yearly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for successful transactions"
// @Router /api/transaction-stats-status/yearly-success [get]
func (h *transactionStatsStatusHandleApi) FindYearlyTransactionStatusSuccess(c echo.Context) error {
	const method = "FindYearlyTransactionStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransactionStatusSuccess(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlySuccess(c)
	}

	so := h.mapper.ToApiResponseTransactionYearStatusSuccess(res)

	logSuccess("success retrieve yearly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusFailed retrieves the monthly transaction status for failed transactions.
// @Summary Get monthly transaction status for failed transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the monthly transaction status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed "Monthly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for failed transactions"
// @Router /api/transaction-stats-status/monthly-failed [get]
func (h *transactionStatsStatusHandleApi) FindMonthlyTransactionStatusFailed(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransactionStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyFailed(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthStatusFailed(res)

	logSuccess("success retrieve monthly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusFailed retrieves the yearly transaction status for failed transactions.
// @Summary Get yearly transaction status for failed transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the yearly transaction status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed "Yearly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for failed transactions"
// @Router /api/transaction-stats-status/yearly-failed [get]
func (h *transactionStatsStatusHandleApi) FindYearlyTransactionStatusFailed(c echo.Context) error {
	const method = "FindYearlyTransactionStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransactionStatusFailed(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyFailed(c)
	}

	so := h.mapper.ToApiResponseTransactionYearStatusFailed(res)

	logSuccess("success retrieve yearly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusSuccess retrieves the monthly transaction status for successful transactions.
// @Summary Get monthly transaction status for successful transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the monthly transaction status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess "Monthly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for successful transactions"
// @Router /api/transaction-stats-status/monthly-success-by-card [get]
func (h *transactionStatsStatusHandleApi) FindMonthlyTransactionStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusSuccessByCardNumber"
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

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransactionStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
		Month:      int32(month),
	})

	if err != nil {
		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlySuccess(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthStatusSuccess(res)

	logSuccess("success retrieve monthly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
// @Summary Get yearly transaction status for successful transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the yearly transaction status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param cardNumber query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionYearStatusSuccess "Yearly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for successful transactions"
// @Router /api/transaction-stats-status/yearly-success-by-card [get]
func (h *transactionStatsStatusHandleApi) FindYearlyTransactionStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransactionStatusSuccessByCardNumber"
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

	res, err := h.client.FindYearlyTransactionStatusSuccessByCardNumber(ctx, &pb.FindYearTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlySuccess(c)
	}

	so := h.mapper.ToApiResponseTransactionYearStatusSuccess(res)

	logSuccess("success retrieve yearly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusFailed retrieves the monthly transaction status for failed transactions.
// @Summary Get monthly transaction status for failed transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the monthly transaction status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param cardNumber query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed "Monthly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for failed transactions"
// @Router /api/transaction-stats-status/monthly-failed-by-card [get]
func (h *transactionStatsStatusHandleApi) FindMonthlyTransactionStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusFailedByCardNumber"
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

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransactionStatusFailedByCardNumber(ctx, &pb.FindMonthlyTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
		Month:      int32(month),
	})

	if err != nil {
		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyFailed(c)
	}

	so := h.mapper.ToApiResponseTransactionMonthStatusFailed(res)

	logSuccess("success retrieve monthly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusFailedByCardNumber retrieves the yearly transaction status for failed transactions.
// @Summary Get yearly transaction status for failed transactions
// @Tags Transaction Stats Status
// @Security Bearer
// @Description Retrieve the yearly transaction status for failed transactions by year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed "Yearly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for failed transactions"
// @Router /api/transaction-stats-status/yearly-failed-by-card [get]
func (h *transactionStatsStatusHandleApi) FindYearlyTransactionStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransactionStatusFailedByCardNumber"
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

	res, err := h.client.FindYearlyTransactionStatusFailedByCardNumber(ctx, &pb.FindYearTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyFailed(c)
	}

	so := h.mapper.ToApiResponseTransactionYearStatusFailed(res)

	logSuccess("success retrieve yearly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionStatsStatusHandleApi) startTracingAndLogging(
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

func (s *transactionStatsStatusHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

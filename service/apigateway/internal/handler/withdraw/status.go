package withdrawhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/withdraw"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawStatsStatusHandleApi struct {
	client pb.WithdrawStatsStatusClient

	logger logger.LoggerInterface

	mapper apimapper.WithdrawStatsStatusResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type withdrawStatsStatusHandleDeps struct {
	client pb.WithdrawStatsStatusClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.WithdrawStatsStatusResponseMapper
}

func NewWithdrawStatsStatusHandleApi(params *withdrawStatsStatusHandleDeps) *withdrawStatsStatusHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_stats_status_handler_requests_total",
			Help: "Total number of withdraw stats status requests",
		},
		[]string{"status", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_stats_status_handler_request_duration_seconds",
			Help:    "Duration of withdraw stats status requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "status"},
	)

	withdrawStatsStatusHandleApi := &withdrawStatsStatusHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("withdraw-stats-status-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerWithdraw := params.router.Group("/api/withdraw-stats-status")

	routerWithdraw.GET("/monthly-success", withdrawStatsStatusHandleApi.FindMonthlyWithdrawStatusSuccess)
	routerWithdraw.GET("/yearly-success", withdrawStatsStatusHandleApi.FindYearlyWithdrawStatusSuccess)
	routerWithdraw.GET("/monthly-failed", withdrawStatsStatusHandleApi.FindMonthlyWithdrawStatusFailed)
	routerWithdraw.GET("/yearly-failed", withdrawStatsStatusHandleApi.FindYearlyWithdrawStatusFailed)

	routerWithdraw.GET("/monthly-success-by-card", withdrawStatsStatusHandleApi.FindMonthlyWithdrawStatusSuccessByCardNumber)
	routerWithdraw.GET("/yearly-success-by-card", withdrawStatsStatusHandleApi.FindYearlyWithdrawStatusSuccessByCardNumber)
	routerWithdraw.GET("/monthly-failed-by-card", withdrawStatsStatusHandleApi.FindMonthlyWithdrawStatusFailedByCardNumber)
	routerWithdraw.GET("/yearly-failed-by-card", withdrawStatsStatusHandleApi.FindYearlyWithdrawStatusFailedByCardNumber)

	return withdrawStatsStatusHandleApi
}

// FindMonthlyWithdrawStatusSuccess retrieves the monthly withdraw status for successful transactions.
// @Summary Get monthly withdraw status for successful transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess "Monthly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for successful transactions"
// @Router /api/withdraw-stats-status/monthly-success [get]
func (h *withdrawStatsStatusHandleApi) FindMonthlyWithdrawStatusSuccess(c echo.Context) error {
	const method = "FindMonthlyWithdrawStatusSuccess"
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

	res, err := h.client.FindMonthlyWithdrawStatusSuccess(ctx, &pb.FindMonthlyWithdrawStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("failed to retrieve monthly withdraw status success", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseWithdrawMonthStatusSuccess(res)

	logSuccess("success retrieve monthly withdraw status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusSuccess retrieves the yearly withdraw status for successful transactions.
// @Summary Get yearly withdraw status for successful transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for successful transactions"
// @Router /api/withdraw-stats-status/yearly-success [get]
func (h *withdrawStatsStatusHandleApi) FindYearlyWithdrawStatusSuccess(c echo.Context) error {
	const method = "FindYearlyWithdrawStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyWithdrawStatusSuccess(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly withdraw status success", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseWithdrawYearStatusSuccess(res)

	logSuccess("success retrieve yearly withdraw status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusFailed retrieves the monthly withdraw status for failed transactions.
// @Summary Get monthly withdraw status for failed transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed "Monthly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for failed transactions"
// @Router /api/withdraw-stats-status/monthly-failed [get]
func (h *withdrawStatsStatusHandleApi) FindMonthlyWithdrawStatusFailed(c echo.Context) error {
	const method = "FindMonthlyWithdrawStatusFailed"
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

	res, err := h.client.FindMonthlyWithdrawStatusFailed(ctx, &pb.FindMonthlyWithdrawStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("failed to retrieve monthly withdraw status failed", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusFailed(c)
	}

	so := h.mapper.ToApiResponseWithdrawMonthStatusFailed(res)

	logSuccess("success retrieve monthly withdraw status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusFailed retrieves the yearly withdraw status for failed transactions.
// @Summary Get yearly withdraw status for failed transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for failed transactions"
// @Router /api/withdraw-stats-status/yearly-failed [get]
func (h *withdrawStatsStatusHandleApi) FindYearlyWithdrawStatusFailed(c echo.Context) error {
	const method = "FindYearlyWithdrawStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyWithdrawStatusFailed(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("failed to retrieve yearly withdraw status failed", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusFailed(c)
	}

	so := h.mapper.ToApiResponseWithdrawYearStatusFailed(res)

	logSuccess("success retrieve yearly withdraw status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusSuccessByCardNumber retrieves the monthly withdraw status for successful transactions.
// @Summary Get monthly withdraw status for successful transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess "Monthly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for successful transactions"
// @Router /api/withdraw-stats-status/monthly-success-by-card [get]
func (h *withdrawStatsStatusHandleApi) FindMonthlyWithdrawStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindMonthlyWithdrawStatusSuccessByCardNumber"
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

	res, err := h.client.FindMonthlyWithdrawStatusSuccessCardNumber(ctx, &pb.FindMonthlyWithdrawStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve monthly withdraw status for successful transactions", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusSuccessCardNumber(c)
	}

	so := h.mapper.ToApiResponseWithdrawMonthStatusSuccess(res)

	logSuccess("Success retrieve monthly withdraw status for successful transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusSuccessByCardNumber retrieves the yearly withdraw status for successful transactions.
// @Summary Get yearly withdraw status for successful transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for successful transactions"
// @Router /api/withdraw-stats-status/yearly-success-by-card-number [get]
func (h *withdrawStatsStatusHandleApi) FindYearlyWithdrawStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindYearlyWithdrawStatusSuccessByCardNumber"
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

	res, err := h.client.FindYearlyWithdrawStatusSuccessCardNumber(ctx, &pb.FindYearWithdrawStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly withdraw status for successful transactions", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusSuccessCardNumber(c)
	}

	so := h.mapper.ToApiResponseWithdrawYearStatusSuccess(res)

	logSuccess("Success retrieve yearly withdraw status for successful transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusFailedByCardNumber retrieves the monthly withdraw status for failed transactions.
// @Summary Get monthly withdraw status for failed transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed "Monthly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for failed transactions"
// @Router /api/withdraw-stats-status/monthly-failed-by-card [get]
func (h *withdrawStatsStatusHandleApi) FindMonthlyWithdrawStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindMonthlyWithdrawStatusFailedByCardNumber"
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

	res, err := h.client.FindMonthlyWithdrawStatusFailedCardNumber(ctx, &pb.FindMonthlyWithdrawStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve monthly withdraw status for failed transactions", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusFailedCardNumber(c)
	}

	so := h.mapper.ToApiResponseWithdrawMonthStatusFailed(res)

	logSuccess("Success retrieve monthly withdraw status for failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusFailedByCardNumber retrieves the yearly withdraw status for failed transactions.
// @Summary Get yearly withdraw status for failed transactions
// @Tags Withdraw Stats Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for failed transactions"
// @Router /api/withdraw-stats-status/yearly-failed-by-card [get]
func (h *withdrawStatsStatusHandleApi) FindYearlyWithdrawStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindYearlyWithdrawStatusFailedByCardNumber"
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

	res, err := h.client.FindYearlyWithdrawStatusFailedCardNumber(ctx, &pb.FindYearWithdrawStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly withdraw status for failed transactions", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusFailedCardNumber(c)
	}

	so := h.mapper.ToApiResponseWithdrawYearStatusFailed(res)

	logSuccess("Success retrieve yearly withdraw status for failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *withdrawStatsStatusHandleApi) startTracingAndLogging(
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

func (s *withdrawStatsStatusHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

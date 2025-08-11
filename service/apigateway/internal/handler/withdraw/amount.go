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

type withdrawStatsAmountHandleApi struct {
	client pb.WithdrawStatsAmountServiceClient

	logger logger.LoggerInterface

	mapper apimapper.WithdrawStatsAmountResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type withdrawStatsAmountHandleDeps struct {
	client pb.WithdrawStatsAmountServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.WithdrawStatsAmountResponseMapper
}

func NewWithdrawStatsAmountHandleApi(params *withdrawStatsAmountHandleDeps) *withdrawStatsAmountHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_stats_amount_handler_requests_total",
			Help: "Total number of withdraw stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_stats_amount_handler_request_duration_seconds",
			Help:    "Duration of withdraw stats amount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	withdrawStatsAmountHandleApi := &withdrawStatsAmountHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("withdraw-stats-amount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerWithdraw := params.router.Group("/api/withdraw-stats-amount")

	routerWithdraw.GET("/monthly-amount", withdrawStatsAmountHandleApi.FindMonthlyWithdraws)
	routerWithdraw.GET("/yearly-amount", withdrawStatsAmountHandleApi.FindYearlyWithdraws)

	routerWithdraw.GET("/monthly-amount-card", withdrawStatsAmountHandleApi.FindMonthlyWithdrawsByCardNumber)
	routerWithdraw.GET("/yearly-amount-card", withdrawStatsAmountHandleApi.FindYearlyWithdrawsByCardNumber)

	return withdrawStatsAmountHandleApi
}

// FindMonthlyWithdraws retrieves the monthly withdraws for a specific year.
// @Summary Get monthly withdraws
// @Tags Withdraw Stats Amount
// @Security Bearer
// @Description Retrieve the monthly withdraws for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount "Monthly withdraws"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraws"
// @Router /api/withdraw-stats-amount/monthly [get]
func (h *withdrawStatsAmountHandleApi) FindMonthlyWithdraws(c echo.Context) error {
	const method = "FindMonthlyWithdraws"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyWithdraws(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly withdraws", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdraws(c)
	}

	so := h.mapper.ToApiResponseWithdrawMonthAmount(res)

	logSuccess("Success retrieve monthly withdraws", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdraws retrieves the yearly withdraws for a specific year.
// @Summary Get yearly withdraws
// @Tags Withdraw Stats Amount
// @Security Bearer
// @Description Retrieve the yearly withdraws for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount "Yearly withdraws"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraws"
// @Router /api/withdraw-stats-amount/yearly [get]
func (h *withdrawStatsAmountHandleApi) FindYearlyWithdraws(c echo.Context) error {
	const method = "FindYearlyWithdraws"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyWithdraws(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly withdraws", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdraws(c)
	}

	so := h.mapper.ToApiResponseWithdrawYearAmount(res)

	logSuccess("Success retrieve yearly withdraws", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawsByCardNumber retrieves the monthly withdraws for a specific card number and year.
// @Summary Get monthly withdraws by card number
// @Tags Withdraw Stats Amount
// @Security Bearer
// @Description Retrieve the monthly withdraws for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount "Monthly withdraws by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraws by card number"
// @Router /api/withdraw-stats-amount/monthly-by-card [get]
func (h *withdrawStatsAmountHandleApi) FindMonthlyWithdrawsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyWithdrawsByCardNumber"
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

	res, err := h.client.FindMonthlyWithdrawsByCardNumber(ctx, &pb.FindYearWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly withdraws by card number", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawsByCardNumber(c)
	}

	so := h.mapper.ToApiResponseWithdrawMonthAmount(res)

	logSuccess("Success retrieve monthly withdraws by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawsByCardNumber retrieves the yearly withdraws for a specific card number and year.
// @Summary Get yearly withdraws by card number
// @Tags Withdraw Stats Amount
// @Security Bearer
// @Description Retrieve the yearly withdraws for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount "Yearly withdraws by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraws by card number"
// @Router /api/withdraw-stats-amount/yearly-by-card [get]
func (h *withdrawStatsAmountHandleApi) FindYearlyWithdrawsByCardNumber(c echo.Context) error {
	const method = "FindYearlyWithdrawsByCardNumber"
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

	res, err := h.client.FindYearlyWithdrawsByCardNumber(ctx, &pb.FindYearWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly withdraws by card number", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawsByCardNumber(c)
	}

	so := h.mapper.ToApiResponseWithdrawYearAmount(res)

	logSuccess("Success retrieve yearly withdraws by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *withdrawStatsAmountHandleApi) startTracingAndLogging(
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

func (s *withdrawStatsAmountHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

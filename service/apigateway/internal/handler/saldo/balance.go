package saldohandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/saldo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoStatsBalanceHandleApi struct {
	saldo pb.SaldoStatsBalanceServiceClient

	logger logger.LoggerInterface

	mapper apimapper.SaldoStatsBalanceResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type saldoStatsBalanceHandleDeps struct {
	client pb.SaldoStatsBalanceServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.SaldoStatsBalanceResponseMapper
}

func NewSaldoStatsBalanceHandleApi(params *saldoStatsBalanceHandleDeps) *saldoStatsBalanceHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_stats_balance_handler_requests_total",
			Help: "Total number of saldo stats balance requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_stats_balance_handler_request_duration_seconds",
			Help:    "Duration of saldo stats balance requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	saldoHandler := &saldoStatsBalanceHandleApi{
		saldo:           params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("saldo-stats-balance-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerSaldo := params.router.Group("/api/saldo-stats-balance")

	routerSaldo.GET("/monthly-balances", saldoHandler.FindMonthlySaldoBalances)
	routerSaldo.GET("/yearly-balances", saldoHandler.FindYearlySaldoBalances)

	return saldoHandler
}

// FindMonthlySaldoBalances retrieves monthly saldo balances for a specific year.
// @Summary Get monthly saldo balances
// @Tags Saldo-Stats-Balance
// @Security Bearer
// @Description Retrieve monthly saldo balances for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthSaldoBalances "Monthly saldo balances"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly saldo balances"
// @Router /api/saldo-stats-balance/monthly-balances [get]
func (h *saldoStatsBalanceHandleApi) FindMonthlySaldoBalances(c echo.Context) error {
	const method = "FindMonthlySaldoBalances"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.saldo.FindMonthlySaldoBalances(ctx, &pb.FindYearlySaldo{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly saldo balances", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindMonthlySaldoBalances(c)
	}

	so := h.mapper.ToApiResponseMonthSaldoBalances(res)

	logSuccess("Successfully retrieve monthly saldo balances", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlySaldoBalances retrieves yearly saldo balances for a specific year.
// @Summary Get yearly saldo balances
// @Tags Saldo-Stats-Balance
// @Security Bearer
// @Description Retrieve yearly saldo balances for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearSaldoBalances "Yearly saldo balances"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly saldo balances"
// @Router /api/saldo-stats-balance/yearly-balances [get]
func (h *saldoStatsBalanceHandleApi) FindYearlySaldoBalances(c echo.Context) error {
	const method = "FindYearlySaldoBalances"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.saldo.FindYearlySaldoBalances(ctx, &pb.FindYearlySaldo{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly saldo balances", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindYearlySaldoBalances(c)
	}

	so := h.mapper.ToApiResponseYearSaldoBalances(res)

	logSuccess("Successfully retrieve yearly saldo balances", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *saldoStatsBalanceHandleApi) startTracingAndLogging(
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

func (s *saldoStatsBalanceHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

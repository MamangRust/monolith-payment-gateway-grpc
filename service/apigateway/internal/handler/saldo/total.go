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

type saldoTotalBalanceHandleApi struct {
	saldo pb.SaldoStatsTotalBalanceClient

	logger logger.LoggerInterface

	mapper apimapper.SaldoStatsTotalResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type saldoTotalBalanceHandleDeps struct {
	client pb.SaldoStatsTotalBalanceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.SaldoStatsTotalResponseMapper
}

func NewSaldoTotalBalanceHandleApi(params *saldoTotalBalanceHandleDeps) *saldoTotalBalanceHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_stats_total_balance_handler_requests_total",
			Help: "Total number of saldo query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_stats_total_balance_handler_request_duration_seconds",
			Help:    "Duration of saldo query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	saldoHandler := &saldoTotalBalanceHandleApi{
		saldo:           params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("saldo-stats-total-balance-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerSaldo := params.router.Group("/api/saldo-stats-total-balance")

	routerSaldo.GET("/monthly-total-balance", saldoHandler.FindMonthlyTotalSaldoBalance)
	routerSaldo.GET("/yearly-total-balance", saldoHandler.FindYearTotalSaldoBalance)

	return saldoHandler
}

// FindMonthlyTotalSaldoBalance retrieves the total saldo balance for a specific month and year.
// @Summary Get monthly total saldo balance
// @Tags Saldo-Stats-Total-Balance
// @Security Bearer
// @Description Retrieve the total saldo balance for a specific month and year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseMonthTotalSaldo "Monthly total saldo balance"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly total saldo balance"
// @Router /api/saldo-stats-total-balance/monthly-total-balance [get]
func (h *saldoTotalBalanceHandleApi) FindMonthlyTotalSaldoBalance(c echo.Context) error {
	const method = "FindMonthlyTotalSaldoBalance"
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

	res, err := h.saldo.FindMonthlyTotalSaldoBalance(ctx, &pb.FindMonthlySaldoTotalBalance{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve monthly total saldo balance", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindMonthlyTotalSaldoBalance(c)
	}

	so := h.mapper.ToApiResponseMonthTotalSaldo(res)

	logSuccess("Successfully retrieve monthly total saldo balance", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalSaldoBalance retrieves the total saldo balance for a specific year.
// @Summary Get yearly total saldo balance
// @Tags Saldo-Stats-Total-Balance
// @Security Bearer
// @Description Retrieve the total saldo balance for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearTotalSaldo "Yearly total saldo balance"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly total saldo balance"
// @Router /api/saldo-stats-total-balance/yearly-total-balance [get]
func (h *saldoTotalBalanceHandleApi) FindYearTotalSaldoBalance(c echo.Context) error {
	const method = "FindYearTotalSaldoBalance"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.saldo.FindYearTotalSaldoBalance(ctx, &pb.FindYearlySaldo{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly total saldo balance", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindYearTotalSaldoBalance(c)
	}

	so := h.mapper.ToApiResponseYearTotalSaldo(res)

	logSuccess("Successfully retrieve yearly total saldo balance", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *saldoTotalBalanceHandleApi) startTracingAndLogging(
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

func (s *saldoTotalBalanceHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

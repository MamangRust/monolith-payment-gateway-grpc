package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoStatisticsService struct {
	ctx                  context.Context
	errorhandler         errorhandler.SaldoStatisticErrorHandler
	mencache             mencache.SaldoStatisticCache
	trace                trace.Tracer
	logger               logger.LoggerInterface
	mapping              responseservice.SaldoResponseMapper
	saldoStatsRepository repository.SaldoStatisticsRepository
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewSaldoStatisticsService(ctx context.Context, errorhandler errorhandler.SaldoStatisticErrorHandler,
	mencache mencache.SaldoStatisticCache, saldoStatsRepository repository.SaldoStatisticsRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoStatisticsService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_statistics_service_request_total",
			Help: "Total number of requests to the SaldoStatisticsService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_statistics_service_request_duration_seconds",
			Help:    "Histogram of request durations for the SaldoStatisticsService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &saldoStatisticsService{
		ctx:                  ctx,
		errorhandler:         errorhandler,
		mencache:             mencache,
		trace:                otel.Tracer("saldo-statistics-service"),
		saldoStatsRepository: saldoStatsRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *saldoStatisticsService) FindMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse) {

	year := req.Year
	month := req.Month

	const method = "FindMonthlyTotalSaldoBalance"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetMonthlyTotalSaldoBalanceCache(req); found {
		logSuccess("Successfully fetched monthly total saldo balance from cache", zap.Int("year", year), zap.Int("month", month))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetMonthlyTotalSaldoBalance(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTotalSaldoBalanceError(err, method, "FAILED_FIND_MONTHLY_TOTAL_SALDO_BALANCE", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToSaldoMonthTotalBalanceResponses(res)

	s.mencache.SetMonthlyTotalSaldoCache(req, responses)

	logSuccess("Successfully fetched monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	return responses, nil
}

func (s *saldoStatisticsService) FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	const method = "FindYearTotalSaldoBalance"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetYearTotalSaldoBalanceCache(year); found {
		logSuccess("Successfully fetched yearly total saldo balance from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetYearTotalSaldoBalance(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTotalSaldoBalanceError(err, method, "FAILED_FIND_YEARLY_TOTAL_SALDO_BALANCE", span, &status, zap.Error(err))
	}
	so := s.mapping.ToSaldoYearTotalBalanceResponses(res)

	s.mencache.SetYearTotalSaldoBalanceCache(year, so)

	logSuccess("Successfully fetched yearly total saldo balance", zap.Int("year", year))

	return so, nil
}

func (s *saldoStatisticsService) FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	const method = "FindMonthlySaldoBalances"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetMonthlySaldoBalanceCache(year); found {
		logSuccess("Successfully fetched monthly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetMonthlySaldoBalances(year)

	if err != nil {
		return s.errorhandler.HandleMonthlySaldoBalancesError(err, method, "FAILED_FIND_MONTHLY_SALDO_BALANCES", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToSaldoMonthBalanceResponses(res)

	s.mencache.SetMonthlySaldoBalanceCache(year, responses)

	logSuccess("Successfully fetched monthly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoStatisticsService) FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	const method = "FindYearlySaldoBalances"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetYearlySaldoBalanceCache(year); found {
		logSuccess("Successfully fetched yearly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetYearlySaldoBalances(year)

	if err != nil {
		return s.errorhandler.HandleYearlySaldoBalancesError(err, method, "FAILED_FIND_YEARLY_SALDO_BALANCES", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToSaldoYearBalanceResponses(res)

	s.mencache.SetYearlySaldoBalanceCache(year, responses)

	logSuccess("Successfully fetched yearly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoStatisticsService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Info(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *saldoStatisticsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

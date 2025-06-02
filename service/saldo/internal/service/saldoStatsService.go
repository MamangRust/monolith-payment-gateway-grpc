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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalSaldoBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalSaldoBalance")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	if cache := s.mencache.GetMonthlyTotalSaldoBalanceCache(req); cache != nil {
		s.logger.Debug("Successfully fetched monthly total saldo balance from cache", zap.Int("year", year), zap.Int("month", month))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetMonthlyTotalSaldoBalance(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTotalSaldoBalanceError(err, "FindMonthlyTotalSaldoBalance", "FAILED_FIND_MONTHLY_TOTAL_SALDO_BALANCE", span, &status)
	}

	responses := s.mapping.ToSaldoMonthTotalBalanceResponses(res)

	s.mencache.SetMonthlyTotalSaldoCache(req, responses)

	s.logger.Debug("Successfully fetched monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	return responses, nil
}

func (s *saldoStatisticsService) FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearTotalSaldoBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearTotalSaldoBalance")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly total saldo balance", zap.Int("year", year))

	if cache := s.mencache.GetYearTotalSaldoBalanceCache(year); cache != nil {
		s.logger.Debug("Successfully fetched yearly total saldo balance from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetYearTotalSaldoBalance(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTotalSaldoBalanceError(err, "FindYearTotalSaldoBalance", "FAILED_FIND_YEARLY_TOTAL_SALDO_BALANCE", span, &status)
	}
	so := s.mapping.ToSaldoYearTotalBalanceResponses(res)

	s.mencache.SetYearTotalSaldoBalanceCache(year, so)

	s.logger.Debug("Successfully fetched yearly total saldo balance", zap.Int("year", year))

	return so, nil
}

func (s *saldoStatisticsService) FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlySaldoBalances", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlySaldoBalances")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly saldo balances", zap.Int("year", year))

	if cache := s.mencache.GetMonthlySaldoBalanceCache(year); cache != nil {
		s.logger.Debug("Successfully fetched monthly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetMonthlySaldoBalances(year)

	if err != nil {
		return s.errorhandler.HandleMonthlySaldoBalancesError(err, "FindMonthlySaldoBalances", "FAILED_FIND_MONTHLY_SALDO_BALANCES", span, &status)
	}

	responses := s.mapping.ToSaldoMonthBalanceResponses(res)

	s.mencache.SetMonthlySaldoBalanceCache(year, responses)

	s.logger.Debug("Successfully fetched monthly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoStatisticsService) FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlySaldoBalances", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlySaldoBalances")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly saldo balances", zap.Int("year", year))

	if cache := s.mencache.GetYearlySaldoBalanceCache(year); cache != nil {
		s.logger.Debug("Successfully fetched yearly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.saldoStatsRepository.GetYearlySaldoBalances(year)

	if err != nil {
		return s.errorhandler.HandleYearlySaldoBalancesError(err, "FindYearlySaldoBalances", "FAILED_FIND_YEARLY_SALDO_BALANCES", span, &status)
	}

	responses := s.mapping.ToSaldoYearBalanceResponses(res)

	s.mencache.SetYearlySaldoBalanceCache(year, responses)

	s.logger.Debug("Successfully fetched yearly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoStatisticsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

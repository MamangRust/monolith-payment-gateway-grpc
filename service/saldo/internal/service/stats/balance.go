package saldostatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/saldo"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type saldoStatsBalanceDeps struct {
	Cache        mencache.SaldoStatsBalanceCache
	ErrorHandler errorhandler.SaldoStatisticErrorHandler

	Repository repository.SaldoStatsBalanceRepository

	Logger logger.LoggerInterface

	Mapper responseservice.SaldoStatisticBalanceResponseMapper
}

type saldoStatsBalanceService struct {
	mencache mencache.SaldoStatsBalanceCache

	repository repository.SaldoStatsBalanceRepository

	errorhandler errorhandler.SaldoStatisticErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.SaldoStatisticBalanceResponseMapper

	observability observability.TraceLoggerObservability
}

func NewSaldoStatsBalanceService(params *saldoStatsBalanceDeps) SaldoStatsBalanceService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "saldo_stats_balance_service_request_total",
		Help: "The total number of requests SaldoStatsBalanceService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "saldo_stats_balance_service_request_duration_seconds",
		Help:    "The duration of requests SaldoStatsBalanceService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("saldo-stats-balance-service"), params.Logger, requestCounter, requestDuration)

	return &saldoStatsBalanceService{
		mencache:      params.Cache,
		repository:    params.Repository,
		errorhandler:  params.ErrorHandler,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlySaldoBalances retrieves saldo balances for each month in the specified year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The target year to retrieve monthly saldo balances.
//
// Returns:
//   - []*response.SaldoMonthBalanceResponse: List of saldo balances per month.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoStatsBalanceService) FindMonthlySaldoBalances(ctx context.Context, year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	const method = "FindMonthlySaldoBalances"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetMonthlySaldoBalanceCache(ctx, year); found {
		logSuccess("Successfully fetched monthly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.repository.GetMonthlySaldoBalances(ctx, year)

	if err != nil {
		return s.errorhandler.HandleMonthlySaldoBalancesError(err, method, "FAILED_FIND_MONTHLY_SALDO_BALANCES", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToSaldoMonthBalanceResponses(res)

	s.mencache.SetMonthlySaldoBalanceCache(ctx, year, responses)

	logSuccess("Successfully fetched monthly saldo balances", zap.Int("year", year))

	return responses, nil
}

// FindYearlySaldoBalances retrieves saldo balances aggregated by year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The target year to retrieve yearly saldo balances.
//
// Returns:
//   - []*response.SaldoYearBalanceResponse: List of saldo balances per year.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoStatsBalanceService) FindYearlySaldoBalances(ctx context.Context, year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	const method = "FindYearlySaldoBalances"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetYearlySaldoBalanceCache(ctx, year); found {
		logSuccess("Successfully fetched yearly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.repository.GetYearlySaldoBalances(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearlySaldoBalancesError(err, method, "FAILED_FIND_YEARLY_SALDO_BALANCES", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToSaldoYearBalanceResponses(res)

	s.mencache.SetYearlySaldoBalanceCache(ctx, year, responses)

	logSuccess("Successfully fetched yearly saldo balances", zap.Int("year", year))

	return responses, nil
}

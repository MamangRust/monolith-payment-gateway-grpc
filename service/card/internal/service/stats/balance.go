package cardstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	cardstatsmencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
)

type cardStatsBalanceService struct {
	errorHandler errorhandler.CardStatisticErrorHandler

	cache cardstatsmencache.CardStatsBalanceCache

	repository repository.CardStatsBalanceRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticBalanceResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsBalanceServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticErrorHandler

	Cache cardstatsmencache.CardStatsBalanceCache

	Repository repository.CardStatsBalanceRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticBalanceResponseMapper
}

func NewCardStatsBalanceService(params *cardStatsBalanceServiceDeps) CardStatsBalanceService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_balance_request_count",
		Help: "Number of card statistic requests CardStatsBalanceService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_balance_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsBalanceService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-balance-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsBalanceService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyBalance retrieves monthly balance statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the monthly balances are requested
//
// Returns:
//   - A slice of CardResponseMonthBalance or an error response if the operation fails.
func (s *cardStatsBalanceService) FindMonthlyBalance(ctx context.Context, year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	const method = "FindMonthlyBalance"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyBalanceCache(ctx, year); found {
		logSuccess("Monthly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyBalance(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyBalanceError(err, method, "FAILED_MONTHLY_BALANCE", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyBalances(res)

	s.cache.SetMonthlyBalanceCache(ctx, year, so)

	logSuccess("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyBalance retrieves yearly balance statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the yearly balances are requested
//
// Returns:
//   - A slice of CardResponseYearlyBalance or an error response if the operation fails.
func (s *cardStatsBalanceService) FindYearlyBalance(ctx context.Context, year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	const method = "FindYearlyBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyBalanceCache(ctx, year); found {
		logSuccess("Yearly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyBalance(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyBalanceError(err, method, "FAILED_YEARLY_BALANCE", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyBalances(res)

	s.cache.SetYearlyBalanceCache(ctx, year, so)

	logSuccess("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

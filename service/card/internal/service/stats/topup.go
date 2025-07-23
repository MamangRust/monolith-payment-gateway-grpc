package cardstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	cardstatsmencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cardStatsTopupService struct {
	errorHandler errorhandler.CardStatisticErrorHandler

	cache cardstatsmencache.CardStatsTopupCache

	repository repository.CardStatsTopupRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsTopupServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticErrorHandler

	Cache cardstatsmencache.CardStatsTopupCache

	Repository repository.CardStatsTopupRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsTopupService(params *cardStatsTopupServiceDeps) CardStatsTopupService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_topup_amount_request_count",
		Help: "Number of card statistic requests CardStatsTopupService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_topup_amount_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsTopupService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-topup-amount-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsTopupService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTopupAmount retrieves monthly top-up statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the monthly top-up data is requested
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTopupService) FindMonthlyTopupAmount(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTopupAmount"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupCache(ctx, year); found {
		logSuccess("Monthly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTopupAmount(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyTopupAmountError(err, method, "FAILED_MONTHLY_TOPUP_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTopupCache(ctx, year, so)

	logSuccess("Monthly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyTopupAmount retrieves yearly top-up statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the yearly top-up data is requested
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTopupService) FindYearlyTopupAmount(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTopupAmount"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupCache(ctx, year); found {
		logSuccess("Yearly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTopupAmount(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyTopupAmountError(err, method, "FAILED_YEARLY_TOPUP_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTopupCache(ctx, year, so)

	logSuccess("Yearly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

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

type cardStatsTransactionService struct {
	errorHandler errorhandler.CardStatisticErrorHandler

	cache cardstatsmencache.CardStatsTransactionCache

	repository repository.CardStatsTransactionRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsTransactionServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticErrorHandler

	Cache cardstatsmencache.CardStatsTransactionCache

	Repository repository.CardStatsTransactionRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsTransactionService(params *cardStatsTransactionServiceDeps) CardStatsTransactionService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_transaction_amount_request_count",
		Help: "Number of card statistic requests CardStatsTransactionService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_transaction_amount_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsTransactionService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-transaction-amount-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsTransactionService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTransactionAmount retrieves monthly transaction statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the monthly transaction data is requested
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTransactionService) FindMonthlyTransactionAmount(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransactionAmount"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransactionCache(ctx, year); found {
		logSuccess("Monthly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionAmount(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyTransactionAmountError(err, method, "FAILED_MONTHLY_TRANSACTION_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTransactionCache(ctx, year, so)

	logSuccess("Monthly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyTransactionAmount retrieves yearly transaction statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the yearly transaction data is requested
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTransactionService) FindYearlyTransactionAmount(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransactionAmount"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransactionCache(ctx, year); found {
		logSuccess("Yearly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionAmount(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyTransactionAmountError(err, method, "FAILED_YEARLY_TRANSACTION_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTransactionCache(ctx, year, so)

	logSuccess("Yearly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

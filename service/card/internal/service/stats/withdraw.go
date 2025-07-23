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

type cardStatsWithdrawService struct {
	errorHandler errorhandler.CardStatisticErrorHandler

	cache cardstatsmencache.CardStatsWithdrawCache

	repository repository.CardStatsWithdrawRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsWithdrawServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticErrorHandler

	Cache cardstatsmencache.CardStatsWithdrawCache

	Repository repository.CardStatsWithdrawRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsWithdrawService(params *cardStatsWithdrawServiceDeps) CardStatsWithdrawService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_withdraw_amount_request_count",
		Help: "Number of card statistic requests CardStatsWithdrawService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_withdraw_amount_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsWithdrawService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-withdraw-amount-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsWithdrawService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyWithdrawAmount retrieves monthly withdraw statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the monthly withdraw data is requested
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsWithdrawService) FindMonthlyWithdrawAmount(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyWithdrawAmount"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyWithdrawCache(ctx, year); found {
		logSuccess("Monthly withdraw amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyWithdrawAmount(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyWithdrawAmountError(err, method, "FAILED_MONTHLY_WITHDRAW_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyWithdrawCache(ctx, year, so)

	logSuccess("Monthly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyWithdrawAmount retrieves yearly withdraw statistics across all card numbers for a given year.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the yearly withdraw data is requested
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsWithdrawService) FindYearlyWithdrawAmount(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawAmount"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyWithdrawCache(ctx, year); found {
		logSuccess("Yearly withdraw amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyWithdrawAmount(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyWithdrawAmountError(err, method, "FAILED_YEARLY_WITHDRAW_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyWithdrawCache(ctx, year, so)

	logSuccess("Yearly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

package cardstatsbycard

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	cardstatsmencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
)

type cardStatsBalanceByCardService struct {
	errorHandler errorhandler.CardStatisticByNumberErrorHandler

	cache cardstatsmencache.CardStatsBalanceByCardCache

	repository repository.CardStatsBalanceByCardRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticBalanceResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsBalanceByCardServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticByNumberErrorHandler

	Cache cardstatsmencache.CardStatsBalanceByCardCache

	Repository repository.CardStatsBalanceByCardRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticBalanceResponseMapper
}

func NewCardStatsBalanceByCardService(params *cardStatsBalanceByCardServiceDeps) CardStatsBalanceByCardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_balance_by_card_request_count",
		Help: "Number of card statistic requests CardStatsBalanceByCardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_balance_by_card_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsBalanceByCardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-balance-by-card-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsBalanceByCardService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyBalanceByCardNumber retrieves monthly balance statistics for a specific card number and year.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the month, year, and card number
//
// Returns:
//   - A slice of CardResponseMonthBalance or an error response if the operation fails.
func (s *cardStatsBalanceByCardService) FindMonthlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	const method = "FindMonthlyBalanceByCardNumber"

	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyBalanceByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly balance card", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyBalancesByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyBalanceByCardNumberError(err, method, "FAILED_MONTHLY_BALANCE_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyBalances(res)

	s.cache.SetMonthlyBalanceByNumberCache(ctx, req, so)

	logSuccess("Successfully fetched monthly balance card", zap.Int("year", year))

	return so, nil
}

// FindYearlyBalanceByCardNumber retrieves yearly balance statistics for a specific card number and year.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the year and card number
//
// Returns:
//   - A slice of CardResponseYearlyBalance or an error response if the operation fails.
func (s *cardStatsBalanceByCardService) FindYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	const method = "FindYearlyBalanceByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", req.Year), attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyBalanceByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly balance card", zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.repository.GetYearlyBalanceByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyBalanceByCardNumberError(err, method, "FAILED_YEARLY_BALANCE_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyBalances(res)

	s.cache.SetYearlyBalanceByNumberCache(ctx, req, so)

	logSuccess("Successfully fetched yearly balance card", zap.Int("year", req.Year))

	return so, nil
}

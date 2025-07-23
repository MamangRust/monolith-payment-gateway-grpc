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

type cardStatsTransactionByCardService struct {
	errorHandler errorhandler.CardStatisticByNumberErrorHandler

	cache cardstatsmencache.CardStatsTransactionByCardCache

	repository repository.CardStatsTransactionByCardRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsTransactionByCardServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticByNumberErrorHandler

	Cache cardstatsmencache.CardStatsTransactionByCardCache

	Repository repository.CardStatsTransactionByCardRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsTransactionByCardService(params *cardStatsTransactionByCardServiceDeps) CardStatsTransactionByCardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_transaction_by_card_request_count",
		Help: "Number of card statistic requests CardStatsTransactionByCardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_transaction_by_card_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsTransactionByCardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-transaction-by-card-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsTransactionByCardService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTransactionAmountByCardNumber retrieves monthly transaction statistics for a specific card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the month, year, and card number
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTransactionByCardService) FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransactionAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransactionByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionAmountByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTransactionAmountByCardNumberError(err, method, "FAILED_MONTHLY_TRANSACTION_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTransactionByNumberCache(ctx, req, so)

	logSuccess("Successfully fetched monthly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

// FindYearlyTransactionAmountByCardNumber retrieves yearly transaction statistics for a specific card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the year and card number
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTransactionByCardService) FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransactionAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransactionByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionAmountByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTransactionAmountByCardNumberError(err, method, "FAILED_YEARLY_TRANSACTION_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTransactionByNumberCache(ctx, req, so)

	s.logger.Debug("Yearly transaction amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

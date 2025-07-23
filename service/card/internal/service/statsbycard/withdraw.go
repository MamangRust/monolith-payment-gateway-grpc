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

type cardStatsWithdrawByCardService struct {
	errorHandler errorhandler.CardStatisticByNumberErrorHandler

	cache cardstatsmencache.CardStatsWithdrawByCardCache

	repository repository.CardStatsWithdrawByCardRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsWithdrawByCardServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticByNumberErrorHandler

	Cache cardstatsmencache.CardStatsWithdrawByCardCache

	Repository repository.CardStatsWithdrawByCardRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsWithdrawByCardService(params *cardStatsWithdrawByCardServiceDeps) CardStatsWithdrawByCardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_withdraw_by_card_request_count",
		Help: "Number of card statistic requests CardStatsWithdrawByCardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_withdraw_by_card_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsWithdrawByCardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-withdraw-by-card-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsWithdrawByCardService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyWithdrawAmountByCardNumber retrieves monthly withdraw statistics for a specific card number and year.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the month, year, and card number
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsWithdrawByCardService) FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "RefreshToken"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyWithdrawByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyWithdrawAmountByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyWithdrawAmountByCardNumberError(err, method, "FAILED_MONTHLY_WITHDRAW_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyWithdrawByNumberCache(ctx, req, so)

	logSuccess("Successfully fetched monthly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

// FindYearlyWithdrawAmountByCardNumber retrieves yearly withdraw statistics for a specific card number and year.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the year and card number
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsWithdrawByCardService) FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyWithdrawByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyWithdrawAmountByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyWithdrawAmountByCardNumberError(err, method, "FAILED_YEARLY_WITHDRAW_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyWithdrawByNumberCache(ctx, req, so)

	logSuccess("Successfully fetched yearly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

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

type cardStatsTopupByCardService struct {
	errorHandler errorhandler.CardStatisticByNumberErrorHandler

	cache cardstatsmencache.CardStatsTopupByCardCache

	repository repository.CardStatsTopupByCardRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsTopupByCardServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticByNumberErrorHandler

	Cache cardstatsmencache.CardStatsTopupByCardCache

	Repository repository.CardStatsTopupByCardRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsTopupByCardService(params *cardStatsTopupByCardServiceDeps) CardStatsTopupByCardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_topup_by_card_request_count",
		Help: "Number of card statistic requests CardStatsTopupByCardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_topup_by_card_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsTopupByCardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-topup-by-card-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsTopupByCardService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTopupAmountByCardNumber retrieves monthly top-up statistics for a specific card number and year.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the month, year, and card number
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTopupByCardService) FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTopupAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTopupAmountByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTopupAmountByCardNumberError(err, "FindMonthlyTopupAmount", "FAILED_MONTHLY_TOPUP_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTopupByNumberCache(ctx, req, so)

	s.logger.Debug("Monthly topup amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyTopupAmountByCardNumber retrieves yearly top-up statistics for a specific card number and year.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: a request object containing the year and card number
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTopupByCardService) FindYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTopupAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTopupAmountByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupAmountByCardNumberError(err, method, "FAILED_YEARLY_TOPUP_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTopupByNumberCache(ctx, req, so)

	logSuccess("Successfully fetched yearly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

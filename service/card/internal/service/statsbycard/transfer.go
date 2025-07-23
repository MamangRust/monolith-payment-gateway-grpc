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

type cardStatsTransferByCardService struct {
	errorHandler errorhandler.CardStatisticByNumberErrorHandler

	cache cardstatsmencache.CardStatsTransferByCardCache

	repository repository.CardStatsTransferByCardRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsTransferByCardServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticByNumberErrorHandler

	Cache cardstatsmencache.CardStatsTransferByCardCache

	Repository repository.CardStatsTransferByCardRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsTransferByCardService(params *cardStatsTransferByCardServiceDeps) CardStatsTransferByCardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_transfer_by_card_request_count",
		Help: "Number of card statistic requests CardStatsTransferByCardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_transfer_by_card_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsTransferByCardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-transfer-by-card-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsTransferByCardService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTransferAmountBySender retrieves monthly transfer statistics for a specific sender card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: request containing year, month, and sender card number
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTransferByCardService) FindMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountBySender"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferBySenderCache(ctx, req); found {
		logSuccess("Cache hit for monthly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransferAmountBySender(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountBySenderError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_BY_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTransferBySenderCache(ctx, req, so)

	logSuccess("Successfully fetched monthly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

// FindYearlyTransferAmountBySender retrieves yearly transfer statistics for a specific sender card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: request containing year and sender card number
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTransferByCardService) FindYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountBySender"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferBySenderCache(ctx, req); found {
		logSuccess("Cache hit for yearly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTransferAmountBySender(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountBySenderError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_BY_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTransferBySenderCache(ctx, req, so)

	logSuccess("Successfully fetched yearly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransferAmountByReceiver retrieves monthly transfer statistics for a specific receiver card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: request containing year, month, and receiver card number
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTransferByCardService) FindMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountByReceiver"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferByReceiverCache(ctx, req); found {
		logSuccess("Cache hit for monthly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransferAmountByReceiver(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountByReceiverError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_BY_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTransferByReceiverCache(ctx, req, so)

	logSuccess("Successfully fetched monthly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

// FindYearlyTransferAmountByReceiver retrieves yearly transfer statistics for a specific receiver card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - req: request containing year and receiver card number
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTransferByCardService) FindYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountByReceiver"

	cardNumber := req.CardNumber
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferByReceiverCache(ctx, req); found {
		logSuccess("Cache hit for yearly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTransferAmountByReceiver(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountByReceiverError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_BY_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTransferByReceiverCache(ctx, req, so)

	logSuccess("Successfully fetched yearly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

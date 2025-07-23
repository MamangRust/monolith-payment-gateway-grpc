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

type cardStatsTransferService struct {
	errorHandler errorhandler.CardStatisticErrorHandler

	cache cardstatsmencache.CardStatsTransferCache

	repository repository.CardStatsTransferRepository

	logger logger.LoggerInterface

	mapper responseservice.CardStatisticAmountResponseMapper

	observability observability.TraceLoggerObservability
}

type cardStatsTransferServiceDeps struct {
	ErrorHandler errorhandler.CardStatisticErrorHandler

	Cache cardstatsmencache.CardStatsTransferCache

	Repository repository.CardStatsTransferRepository

	Logger logger.LoggerInterface

	Mapper responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsTransferService(params *cardStatsTransferServiceDeps) CardStatsTransferService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_stats_transfer_amount_request_count",
		Help: "Number of card statistic requests CardStatsTransferService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_stats_transfer_amount_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatsTransferService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-stats-transfer-amount-service"), params.Logger, requestCounter, requestDuration)

	return &cardStatsTransferService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTransferAmountSender retrieves total monthly transfer amounts from all cards acting as sender.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the monthly data is requested
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTransferService) FindMonthlyTransferAmountSender(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountSender"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferSenderCache(ctx, year); found {
		logSuccess("Monthly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransferAmountSender(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountSenderError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTransferSenderCache(ctx, year, so)

	logSuccess("Monthly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyTransferAmountSender retrieves total yearly transfer amounts from all cards acting as sender.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the yearly data is requested
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTransferService) FindYearlyTransferAmountSender(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountSender"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferSenderCache(ctx, year); found {
		logSuccess("Yearly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTransferAmountSender(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountSenderError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTransferSenderCache(ctx, year, so)

	logSuccess("Yearly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindMonthlyTransferAmountReceiver retrieves total monthly transfer amounts for all cards acting as receiver.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the monthly data is requested
//
// Returns:
//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
func (s *cardStatsTransferService) FindMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountReceiver"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferReceiverCache(ctx, year); found {
		logSuccess("Monthly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransferAmountReceiver(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountReceiverError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetMonthlyAmounts(res)

	s.cache.SetMonthlyTransferReceiverCache(ctx, year, so)

	logSuccess("Monthly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

// FindYearlyTransferAmountReceiver retrieves total yearly transfer amounts for all cards acting as receiver.
//
// Parameters:
//   - ctx: the context for the operation
//   - year: the year for which the yearly data is requested
//
// Returns:
//   - A slice of CardResponseYearAmount or an error response if the operation fails.
func (s *cardStatsTransferService) FindYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountReceiver"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferReceiverCache(ctx, year); found {
		logSuccess("Yearly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.repository.GetYearlyTransferAmountReceiver(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountReceiverError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToGetYearlyAmounts(res)

	s.cache.SetYearlyTransferReceiverCache(ctx, year, so)

	logSuccess("Yearly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

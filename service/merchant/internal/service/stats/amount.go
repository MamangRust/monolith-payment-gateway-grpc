package merchantstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantStatsAmountDeps struct {
	Cache mencache.MerchantStatsAmountCache

	ErrorHandler errorhandler.MerchantStatisticErrorHandler

	Repository repository.MerchantStatsAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantAmountResponseMapper
}

type merchantStatsAmountService struct {
	mencache mencache.MerchantStatsAmountCache

	repository repository.MerchantStatsAmountRepository

	errorHandler errorhandler.MerchantStatisticErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsAmountService(params *merchantStatsAmountDeps) MerchantStatsAmountService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_amount_service_request_total",
		Help: "The total number of requests MerchantStatisticService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_amount_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-amount-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsAmountService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyAmountMerchant retrieves the monthly transaction amount statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the monthly amount statistics should be retrieved.
//
// Returns:
//   - []*response.MerchantResponseMonthlyAmount: A slice of monthly amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsAmountService) FindMonthlyAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	const method = "FindMonthlyAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyAmountMerchant(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyAmountMerchantError(
			err, method, "FAILED_FIND_MONTHLY_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountMerchantCache(ctx, year, so)

	logSuccess("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

// FindYearlyAmountMerchant retrieves the yearly transaction amount statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the yearly amount statistics should be retrieved.
//
// Returns:
//   - []*response.MerchantResponseYearlyAmount: A slice of yearly amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsAmountService) FindYearlyAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	const method = "FindYearlyAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyAmountMerchant(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountMerchantError(
			err, method, "FAILED_FIND_YEARLY_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountMerchantCache(ctx, year, so)

	logSuccess("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

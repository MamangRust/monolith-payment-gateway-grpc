package merchantstatsbymerchantservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbymerchant"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbymerchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantStatsByMerchantAmountServiceDeps struct {
	Cache mencache.MerchantStatsAmountByMerchantCache

	ErrorHandler errorhandler.MerchantStatisticByMerchantErrorHandler

	Repository repository.MerchantStatsByMerchantRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantAmountResponseMapper
}

type merchantStatsByMerchantAmountService struct {
	mencache mencache.MerchantStatsAmountByMerchantCache

	repository repository.MerchantStatsByMerchantRepository

	errorHandler errorhandler.MerchantStatisticByMerchantErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsByMerchantAmountService(params *merchantStatsByMerchantAmountServiceDeps) *merchantStatsByMerchantAmountService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_amount_bymerchant_service_request_total",
		Help: "The total number of requests MerchantStatsByMerchantAmountByMerchantService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_amount_bymerchant_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatsByMerchantAmountByMerchantService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-amount-by-merchant-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsByMerchantAmountService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyAmountByMerchants retrieves monthly transaction amount statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing merchant identifier and year.
//
// Returns:
//   - []*response.MerchantResponseMonthlyAmount: A slice of monthly transaction amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsByMerchantAmountService) FindMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindMonthlyAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID))

		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyAmountByMerchants(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyAmountByMerchantsError(
			err, method, "FAILED_FIND_MONTHLY_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountByMerchantsCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

// FindYearlyAmountByMerchants retrieves yearly transaction amount statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing merchant identifier and year.
//
// Returns:
//   - []*response.MerchantResponseYearlyAmount: A slice of yearly transaction amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsByMerchantAmountService) FindYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindYearlyAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID))

		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyAmountByMerchants(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountByMerchantsError(
			err, method, "FAILED_FIND_YEARLY_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountByMerchantsCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

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

type merchantStatsByMerchantTotalAmountDeps struct {
	Cache mencache.MerchantStatsTotalAmountByMerchantCache

	ErrorHandler errorhandler.MerchantStatisticByMerchantErrorHandler

	Repository repository.MerchantStatsByMerchantRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantTotalAmountResponseMapper
}

type merchantStatsByMerchantTotalAmountService struct {
	mencache mencache.MerchantStatsTotalAmountByMerchantCache

	repository repository.MerchantStatsByMerchantRepository

	errorHandler errorhandler.MerchantStatisticByMerchantErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantTotalAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsByMerchantTotalAmountService(params *merchantStatsByMerchantTotalAmountDeps) *merchantStatsByMerchantTotalAmountService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_totalamount_bymerchant_service_request_total",
		Help: "The total number of requests MerchantStatsByMerchantTotalAmountByApiKeyService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_totalamount_bymerchant_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatsByMerchantTotalAmountByApiKeyService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-totalamount-by-apikey-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsByMerchantTotalAmountService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyTotalAmountByMerchants retrieves monthly total transaction amounts for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing merchant identifier and year.
//
// Returns:
//   - []*response.MerchantResponseMonthlyTotalAmount: A slice of monthly total amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsByMerchantTotalAmountService) FindMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindMonthlyTotalAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyTotalAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyTotalAmountByMerchants(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountByMerchantsError(
			err, method, "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountByMerchantsCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID), zap.Int("year", year))

	return so, nil
}

// FindYearlyTotalAmountByMerchants retrieves yearly total transaction amounts for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing merchant identifier and year.
//
// Returns:
//   - []*response.MerchantResponseYearlyTotalAmount: A slice of yearly total amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsByMerchantTotalAmountService) FindYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindYearlyTotalAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyTotalAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyTotalAmountByMerchants(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountByMerchantsError(
			err, method, "FAILED_FIND_YEARLY_TOTAL_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountByMerchantsCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID), zap.Int("year", year))

	return so, nil
}

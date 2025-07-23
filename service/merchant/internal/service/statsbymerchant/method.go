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

type merchantStatsByMerchantMethodDeps struct {
	Cache mencache.MerchantStatsMethodByMerchantCache

	ErrorHandler errorhandler.MerchantStatisticByMerchantErrorHandler

	Repository repository.MerchantStatsByMerchantRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantPaymentMethodResponseMapper
}

type merchantStatsByMerchantMethodService struct {
	mencache mencache.MerchantStatsMethodByMerchantCache

	repository repository.MerchantStatsByMerchantRepository

	errorHandler errorhandler.MerchantStatisticByMerchantErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantPaymentMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsByMerchantMethodService(params *merchantStatsByMerchantMethodDeps) *merchantStatsByMerchantMethodService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_method_bymerchant_service_request_total",
		Help: "The total number of requests MerchantStatsByMerchantMethodByApiKeyService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_method_bymerchant_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatsByMerchantMethodByApiKeyService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-method-by-apikey-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsByMerchantMethodService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyPaymentMethodByMerchants retrieves monthly payment method statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing merchant identifier and year.
//
// Returns:
//   - []*response.MerchantResponseMonthlyPaymentMethod: A slice of monthly payment method statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsByMerchantMethodService) FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindMonthlyPaymentMethodByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyPaymentMethodByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant_id", merchantID))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyPaymentMethodByMerchants(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodByMerchantsError(
			err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHOD_BY_MERCHANTS", span, &status,
			zap.Any("error", err),
		)
	}

	so := s.mapper.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodByMerchantsCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

// FindYearlyPaymentMethodByMerchants retrieves yearly payment method statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing merchant identifier and year.
//
// Returns:
//   - []*response.MerchantResponseYearlyPaymentMethod: A slice of yearly payment method statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsByMerchantMethodService) FindYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindYearlyPaymentMethodByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyPaymentMethodByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant_id", merchantID))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyPaymentMethodByMerchants(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodByMerchantsError(
			err, method, "FAILED_FIND_YEARLY_PAYMENT_METHOD_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodByMerchantsCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

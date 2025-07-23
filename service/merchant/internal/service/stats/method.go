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

type merchantStatsMethodDeps struct {
	Cache mencache.MerchantStatsMethodCache

	ErrorHandler errorhandler.MerchantStatisticErrorHandler

	Repository repository.MerchantStatsMethodRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantPaymentMethodResponseMapper
}

type merchantStatsMethodService struct {
	mencache mencache.MerchantStatsMethodCache

	repository repository.MerchantStatsMethodRepository

	errorHandler errorhandler.MerchantStatisticErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantPaymentMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsMethodService(params *merchantStatsMethodDeps) MerchantStatsMethodService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_method_service_request_total",
		Help: "The total number of requests MerchantStatisticService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_method_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-method-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsMethodService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyPaymentMethodsMerchant retrieves monthly payment method statistics for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the monthly payment method statistics should be retrieved.
//
// Returns:
//   - []*response.MerchantResponseMonthlyPaymentMethod: A slice of monthly payment method statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsMethodService) FindMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	const method = "FindMonthlyPaymentMethodsMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyPaymentMethodsMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyPaymentMethodsMerchant(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodsMerchantError(
			err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHODS_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodsMerchantCache(ctx, year, so)

	logSuccess("Successfully found monthly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

// FindYearlyPaymentMethodMerchant retrieves the yearly payment methods for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the yearly payment methods should be retrieved.
//
// Returns:
//   - []*response.MerchantResponseYearlyPaymentMethod: A slice of yearly payment method statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsMethodService) FindYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	const method = "FindYearlyPaymentMethodMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyPaymentMethodMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyPaymentMethodMerchant(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodMerchantError(
			err, method, "FAILED_FIND_YEARLY_PAYMENT_METHOD_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodMerchantCache(ctx, year, so)

	logSuccess("Successfully found yearly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

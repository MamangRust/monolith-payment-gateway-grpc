package merchantstatsbyapikeyservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbyapikey"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbyapikey"
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

type merchantStatsMethodByApiKeyDeps struct {
	Cache mencache.MerchantStatsMethodByApiKeyCache

	ErrorHandler errorhandler.MerchantStatisticByApikeyErrorHandler

	Repository repository.MerchantStatsMethodByApiKeyRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantPaymentMethodResponseMapper
}

type merchantStatsMethodByApiKeyService struct {
	mencache mencache.MerchantStatsMethodByApiKeyCache

	repository repository.MerchantStatsMethodByApiKeyRepository

	errorHandler errorhandler.MerchantStatisticByApikeyErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantPaymentMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsMethodByApiKeyService(params *merchantStatsMethodByApiKeyDeps) MerchantStatsByApiKeyMethodService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_method_byapikey_service_request_total",
		Help: "The total number of requests MerchantStatsMethodByApiKeyService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_method_byapikey_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatsMethodByApiKeyService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-method-by-apikey-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsMethodByApiKeyService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyPaymentMethodByApikeys retrieves monthly payment method statistics for a merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the API key and the target year.
//
// Returns:
//   - []*response.MerchantResponseMonthlyPaymentMethod: A slice of monthly payment method statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsMethodByApiKeyService) FindMonthlyPaymentMethodByApikeys(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	const method = "FindMonthlyPaymentMethodByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", api_key), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyPaymentMethodByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", api_key), zap.Int("year", year))

		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyPaymentMethodByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodByApikeysError(
			err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHOD_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodByApikeysCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

// FindYearlyPaymentMethodByApikeys retrieves yearly payment method statistics for a merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the API key and the target year.
//
// Returns:
//   - []*response.MerchantResponseYearlyPaymentMethod: A slice of yearly payment method statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsMethodByApiKeyService) FindYearlyPaymentMethodByApikeys(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	const method = "FindYearlyPaymentMethodByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", api_key), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyPaymentMethodByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyPaymentMethodByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodByApikeysError(
			err, method, "FAILED_FIND_YEARLY_PAYMENT_METHOD_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodByApikeysCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

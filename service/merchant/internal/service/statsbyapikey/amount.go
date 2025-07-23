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

type merchantStatsAmountByApiKeyDeps struct {
	Cache mencache.MerchantStatsAmountByApiKeyCache

	ErrorHandler errorhandler.MerchantStatisticByApikeyErrorHandler

	Repository repository.MerchantStatsAmountByApiKeyRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantAmountResponseMapper
}

type merchantStatsAmountByApiKeyService struct {
	mencache mencache.MerchantStatsAmountByApiKeyCache

	repository repository.MerchantStatsAmountByApiKeyRepository

	errorHandler errorhandler.MerchantStatisticByApikeyErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsAmountByApiKeyService(params *merchantStatsAmountByApiKeyDeps) MerchantStatsByApiKeyAmountService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_amount_byapikey_service_request_total",
		Help: "The total number of requests MerchantStatsAmountByApiKeyService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_amount_byapikey_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatsAmountByApiKeyService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-amount-by-apikey-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsAmountByApiKeyService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyAmountByApikeys retrieves monthly transaction amount statistics for a merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the API key and the target year.
//
// Returns:
//   - []*response.MerchantResponseMonthlyAmount: A slice of monthly transaction amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsAmountByApiKeyService) FindMonthlyAmountByApikeys(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	const method = "FindMonthlyAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", api_key), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", api_key), zap.Int("year", year))

		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyAmountByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyAmountByApikeysError(
			err, method, "FAILED_FIND_MONTHLY_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountByApikeysCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

// FindYearlyAmountByApikeys retrieves yearly transaction amount statistics for a merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the API key and the target year.
//
// Returns:
//   - []*response.MerchantResponseYearlyAmount: A slice of yearly transaction amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsAmountByApiKeyService) FindYearlyAmountByApikeys(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	const method = "FindYearlyAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", api_key), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyAmountByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountByApikeysError(
			err, method, "FAILED_FIND_YEARLY_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountByApikeysCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

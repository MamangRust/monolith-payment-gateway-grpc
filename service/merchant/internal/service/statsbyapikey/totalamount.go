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

type merchantStatsTotalAmountByApiKeyDeps struct {
	Cache mencache.MerchantStatsTotalAmountByApiKeyCache

	ErrorHandler errorhandler.MerchantStatisticByApikeyErrorHandler

	Repository repository.MerchantStatsTotalAmountByApiKeyRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantTotalAmountResponseMapper
}

type merchantStatsTotalAmountByApiKeyService struct {
	mencache mencache.MerchantStatsTotalAmountByApiKeyCache

	repository repository.MerchantStatsTotalAmountByApiKeyRepository

	errorHandler errorhandler.MerchantStatisticByApikeyErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantTotalAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsTotalAmountByApiKeyService(params *merchantStatsTotalAmountByApiKeyDeps) MerchantStatsByApiKeyTotalAmountService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_totalamount_byapikey_service_request_total",
		Help: "The total number of requests MerchantStatsTotalAmountByApiKeyService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_totalamount_byapikey_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatsTotalAmountByApiKeyService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-totalamount-by-apikey-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsTotalAmountByApiKeyService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyTotalAmountByApikeys retrieves monthly total transaction amounts for a merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the API key and the target year.
//
// Returns:
//   - []*response.MerchantResponseMonthlyTotalAmount: A slice of monthly total amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsTotalAmountByApiKeyService) FindMonthlyTotalAmountByApikeys(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	const method = "FindMonthlyTotalAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", api_key), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyTotalAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyTotalAmountByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountByApikeysError(
			err, method, "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountByApikeysCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

// FindYearlyTotalAmountByApikeys retrieves yearly total transaction amounts for a merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the API key and the target year.
//
// Returns:
//   - []*response.MerchantResponseYearlyTotalAmount: A slice of yearly total amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsTotalAmountByApiKeyService) FindYearlyTotalAmountByApikeys(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	const method = "FindYearlyTotalAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", api_key), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyTotalAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyTotalAmountByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountByApikeysError(
			err, method, "FAILED_FIND_YEARLY_TOTAL_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountByApikeysCache(ctx, req, so)

	logSuccess("Successfully fetched merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

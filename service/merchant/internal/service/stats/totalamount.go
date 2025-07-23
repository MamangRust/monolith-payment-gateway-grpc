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

type merchantStatsTotalAmountDeps struct {
	Cache mencache.MerchantStatsTotalAmountCache

	ErrorHandler errorhandler.MerchantStatisticErrorHandler

	Repository repository.MerchantStatsTotalAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.MerchantTotalAmountResponseMapper
}

type merchantStatsTotalAmountService struct {
	mencache mencache.MerchantStatsTotalAmountCache

	repository repository.MerchantStatsTotalAmountRepository

	errorHandler errorhandler.MerchantStatisticErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.MerchantTotalAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantStatsTotalAmountService(params *merchantStatsTotalAmountDeps) MerchantStatsTotalAmountService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_stats_totalamount_service_request_total",
		Help: "The total number of requests MerchantStatisticService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_stats_totalamount_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-stats-totalamount-service"), params.Logger, requestCounter, requestDuration)

	return &merchantStatsTotalAmountService{
		mencache:      params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		errorHandler:  params.ErrorHandler,
		observability: observability,
	}
}

// FindMonthlyTotalAmountMerchant retrieves the monthly total transaction amounts for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the monthly total amount statistics should be retrieved.
//
// Returns:
//   - []*response.MerchantResponseMonthlyTotalAmount: A slice of monthly total amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsTotalAmountService) FindMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTotalAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyTotalAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetMonthlyTotalAmountMerchant(ctx, year)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountMerchantError(
			err, method, "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountMerchantCache(ctx, year, so)

	logSuccess("Successfully found monthly total amount for merchant", zap.Int("year", year))

	return so, nil
}

// FindYearlyTotalAmountMerchant retrieves the yearly total transaction amounts for a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the yearly total amount statistics should be retrieved.
//
// Returns:
//   - []*response.MerchantResponseYearlyTotalAmount: A slice of yearly total amount statistics.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantStatsTotalAmountService) FindYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	const method = "FindYearlyTotalAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyTotalAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.repository.GetYearlyTotalAmountMerchant(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountMerchantError(
			err, method, "FAILED_FIND_YEARLY_TOTAL_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapper.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountMerchantCache(ctx, year, so)

	logSuccess("Successfully found yearly total amount for merchant", zap.Int("year", year))

	return so, nil
}

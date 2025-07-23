package transactionstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/stats"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionStatsAmountServiceDeps struct {
	Cache mencache.TransactionStatsAmountCache

	ErrorHandler errorhandler.TransactionStatisticErrorHandler

	Repository repository.TransactionStatsAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransactionStatsAmountResponseMapper
}

type transactionStatsAmountService struct {
	cache mencache.TransactionStatsAmountCache

	errorHandler errorhandler.TransactionStatisticErrorHandler

	repository repository.TransactionStatsAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.TransactionStatsAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransactionStatsAmountService(params *transactionStatsAmountServiceDeps) TransactionStatsAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_amount_service_request_total",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-stats-amount-service"), params.Logger, requestCounter, requestDuration)

	return &transactionStatsAmountService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyAmounts retrieves the total monthly transaction amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionMonthAmountResponse: List of monthly transaction amounts.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsAmountService) FindMonthlyAmounts(ctx context.Context, year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()
	if data, found := s.cache.GetMonthlyAmountsCache(ctx, year); found {
		s.logger.Debug("Successfully fetched monthly amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyAmounts(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyAmountsError(err, method, "FAILED_FIND_MONTHLY_AMOUNTS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionMonthlyAmountResponses(records)

	s.cache.SetMonthlyAmountsCache(ctx, year, responses)

	logSuccess("Successfully fetched monthly amounts", zap.Int("year", year))

	return responses, nil
}

// FindYearlyAmounts retrieves the total yearly transaction amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionYearlyAmountResponse: List of yearly transaction amounts.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsAmountService) FindYearlyAmounts(ctx context.Context, year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyAmountsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyAmounts(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyAmountsError(err, method, "FAILED_FIND_YEARLY_AMOUNTS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionYearlyAmountResponses(records)

	s.cache.SetYearlyAmountsCache(ctx, year, responses)

	logSuccess("Successfully fetched yearly amounts", zap.Int("year", year))

	return responses, nil
}

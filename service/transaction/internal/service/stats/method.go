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

type transactionStatsMethodServiceDeps struct {
	Cache mencache.TransactionStatsMethodCache

	ErrorHandler errorhandler.TransactionStatisticErrorHandler

	Repository repository.TransactionStatsMethodRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransactionStatsMethodResponseMapper
}

type transactionStatsMethodService struct {
	cache mencache.TransactionStatsMethodCache

	errorHandler errorhandler.TransactionStatisticErrorHandler

	repository repository.TransactionStatsMethodRepository

	logger logger.LoggerInterface

	mapper responseservice.TransactionStatsMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransactionStatsMethodService(params *transactionStatsMethodServiceDeps) TransactionStatsMethodService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_method_service_request_total",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_method_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-stats-method-service"), params.Logger, requestCounter, requestDuration)

	return &transactionStatsMethodService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyPaymentMethods retrieves monthly usage statistics for each payment method.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionMonthMethodResponse: List of monthly method statistics.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsMethodService) FindMonthlyPaymentMethods(ctx context.Context, year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyPaymentMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyPaymentMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched monthly payment methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyPaymentMethods(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodsError(err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionMonthlyMethodResponses(records)

	s.cache.SetMonthlyPaymentMethodsCache(ctx, year, responses)

	logSuccess("Successfully fetched monthly payment methods", zap.Int("year", year))

	return responses, nil
}

// FindYearlyPaymentMethods retrieves yearly usage statistics for each payment method.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionYearMethodResponse: List of yearly method statistics.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsMethodService) FindYearlyPaymentMethods(ctx context.Context, year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyPaymentMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyPaymentMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly payment methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyPaymentMethods(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodsError(err, method, "FAILED_FIND_YEARLY_PAYMENT_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionYearlyMethodResponses(records)

	s.cache.SetYearlyPaymentMethodsCache(ctx, year, responses)

	logSuccess("Successfully fetched yearly payment methods", zap.Int("year", year))

	return responses, nil
}

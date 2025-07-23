package transactionstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
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

type transactionStatsStatusServiceDeps struct {
	Cache mencache.TransactionStatsStatusCache

	ErrorHandler errorhandler.TransactionStatisticErrorHandler

	Repository repository.TransactionStatsStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransactionStatsStatusResponseMapper
}

type transactionStatsStatusService struct {
	cache mencache.TransactionStatsStatusCache

	errorHandler errorhandler.TransactionStatisticErrorHandler

	repository repository.TransactionStatsStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.TransactionStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransactionStatsStatusService(params *transactionStatsStatusServiceDeps) TransactionStatsStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_status_service_request_total",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-stats-status-service"), params.Logger, requestCounter, requestDuration)

	return &transactionStatsStatusService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthTransactionStatusSuccess retrieves monthly success statistics for all transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the year and month of the transaction.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusSuccess: List of success statistics.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsStatusService) FindMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	month := req.Month
	year := req.Year

	const method = "FindMonthTransactionStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTransactionStatusSuccessCache(ctx, req); found {
		logSuccess("Successfully fetched monthly Transaction status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthTransactionStatusSuccess(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthTransactionStatusSuccessError(err, method, "FAILED_FIND_MONTH_TRANSACTION_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesMonthStatusSuccess(records)

	s.cache.SetMonthTransactionStatusSuccessCache(ctx, req, so)

	logSuccess("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransactionStatusSuccess retrieves yearly success statistics for all transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionResponseYearStatusSuccess: List of success statistics.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsStatusService) FindYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyTransactionStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearTransactionStatusSuccessCache(ctx, year); found {
		logSuccess("Successfully fetched yearly Transaction status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransactionStatusSuccess(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTransactionStatusSuccessError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesYearStatusSuccess(records)

	s.cache.SetYearTransactionStatusSuccessCache(ctx, year, so)

	logSuccess("Successfully fetched yearly Transaction status success", zap.Int("year", year))

	return so, nil
}

// FindMonthTransactionStatusFailed retrieves monthly failed statistics for all transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the year and month of the transaction.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusFailed: List of failed statistics.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsStatusService) FindMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthTransactionStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTransactionStatusFailedCache(ctx, req); found {
		logSuccess("Successfully fetched monthly Transaction status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthTransactionStatusFailed(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthTransactionStatusFailedError(err, method, "FAILED_MONTHLY_TRANSACTION_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesMonthStatusFailed(records)

	s.cache.SetMonthTransactionStatusFailedCache(ctx, req, so)

	logSuccess("Successfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransactionStatusFailed retrieves yearly failed statistics for all transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionResponseYearStatusFailed: List of failed statistics.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsStatusService) FindYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyTransactionStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearTransactionStatusFailedCache(ctx, year); found {
		logSuccess("Successfully fetched yearly Transaction status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransactionStatusFailed(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTransactionStatusFailedError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesYearStatusFailed(records)

	s.cache.SetYearTransactionStatusFailedCache(ctx, year, so)

	logSuccess("Successfully fetched yearly Transaction status Failed", zap.Int("year", year))

	return so, nil
}

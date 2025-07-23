package transactionstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	cache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/statsbycard"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionStatsByCardStatusServiceDeps struct {
	Cache cache.TransactionStatsByCardStatusCache

	ErrorHandler errorhandler.TransactionStatisticByCardErrorHandler

	Repository repository.TransactonStatsByCardStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransactionStatsStatusResponseMapper
}

type transactionStatsByCardStatusService struct {
	cache cache.TransactionStatsByCardStatusCache

	errorHandler errorhandler.TransactionStatisticByCardErrorHandler

	repository repository.TransactonStatsByCardStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.TransactionStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransactionStatsByCardStatusService(params *transactionStatsByCardStatusServiceDeps) TransactionStatsByCardStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_bycard_status_service_request_totals",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_bycard_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-stats-bycard-status-service"), params.Logger, requestCounter, requestDuration)

	return &transactionStatsByCardStatusService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthTransactionStatusSuccessByCardNumber retrieves monthly successful transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number, year, and month.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusSuccess: List of successful transactions by month.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardStatusService) FindMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTransactionStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTransactionStatusSuccessByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly Transaction status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetMonthTransactionStatusSuccessByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthTransactionStatusSuccessByCardNumberError(err, method, "FAILED_FIND_MONTHLY_TRANSACTION_STATUS_SUCCESS_BY_CARD", span, &status, zap.Int("year", year), zap.Int("month", month), zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesMonthStatusSuccess(records)

	s.cache.SetMonthTransactionStatusSuccessByCardCache(ctx, req, so)

	logSuccess("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

// FindYearlyTransactionStatusSuccessByCardNumber retrieves yearly successful transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number and year.
//
// Returns:
//   - []*response.TransactionResponseYearStatusSuccess: List of successful transactions by year.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardStatusService) FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransactionStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearTransactionStatusSuccessByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly Transaction status success from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransactionStatusSuccessByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyTransactionStatusSuccessByCardNumberError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_SUCCESS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesYearStatusSuccess(records)

	s.cache.SetYearTransactionStatusSuccessByCardCache(ctx, req, so)

	logSuccess("Successfully fetched yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

// FindMonthTransactionStatusFailedByCardNumber retrieves monthly failed transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number, year, and month.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusFailed: List of failed transactions by month.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardStatusService) FindMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTransactionStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTransactionStatusFailedByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly Transaction status Failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetMonthTransactionStatusFailedByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthTransactionStatusFailedByCardNumberError(err, method, "FAILED_MONTHLY_TRANSACTION_STATUS_FAILED_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesMonthStatusFailed(records)

	s.cache.SetMonthTransactionStatusFailedByCardCache(ctx, req, so)

	logSuccess("Successfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

// FindYearlyTransactionStatusFailedByCardNumber retrieves yearly failed transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number and year.
//
// Returns:
//   - []*response.TransactionResponseYearStatusFailed: List of failed transactions by year.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardStatusService) FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransactionStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearTransactionStatusFailedByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly Transaction status Failed from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransactionStatusFailedByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyTransactionStatusFailedByCardNumberError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_FAILED_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransactionResponsesYearStatusFailed(records)

	s.cache.SetYearTransactionStatusFailedByCardCache(ctx, req, so)

	logSuccess("Successfully fetched yearly Transaction status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

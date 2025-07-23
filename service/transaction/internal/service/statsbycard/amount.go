package transactionstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/statsbycard"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionStatsByCardAmountServiceDeps struct {
	Cache mencache.TransactionStatsByCardAmountCache

	ErrorHandler errorhandler.TransactionStatisticByCardErrorHandler

	Repository repository.TransactonStatsByCardAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransactionStatsAmountResponseMapper
}

type transactionStatsByCardAmountService struct {
	cache mencache.TransactionStatsByCardAmountCache

	errorHandler errorhandler.TransactionStatisticByCardErrorHandler

	repository repository.TransactonStatsByCardAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.TransactionStatsAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransactionStatsByCardAmountService(params *transactionStatsByCardAmountServiceDeps) TransactionStatsByCardAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_bycard_amount_service_request_total",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_bycard_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-stats-bycard-amount-service"), params.Logger, requestCounter, requestDuration)

	return &transactionStatsByCardAmountService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyAmountsByCardNumber retrieves monthly transaction amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number, year, and month.
//
// Returns:
//   - []*response.TransactionMonthAmountResponse: List of monthly transaction amounts.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardAmountService) FindMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyAmountsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly amounts by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyAmountsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyAmountsByCardNumberError(err, method, "FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionMonthlyAmountResponses(records)

	s.cache.SetMonthlyAmountsByCardCache(ctx, req, responses)

	logSuccess("Successfully fetched monthly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

// FindYearlyAmountsByCardNumber retrieves yearly transaction amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number and year.
//
// Returns:
//   - []*response.TransactionYearlyAmountResponse: List of yearly transaction amounts.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardAmountService) FindYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyAmountsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly amounts by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyAmountsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyAmountsByCardNumberError(err, method, "FAILED_FIND_YEARLY_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionYearlyAmountResponses(records)

	s.cache.SetYearlyAmountsByCardCache(ctx, req, responses)

	logSuccess("Successfully fetched yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

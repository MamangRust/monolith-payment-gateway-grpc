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

type transactionStatsByCardMethodServiceDeps struct {
	Cache mencache.TransactionStatsByCardMethodCache

	ErrorHandler errorhandler.TransactionStatisticByCardErrorHandler

	Repository repository.TransactionStatsByCardMethodRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransactionStatsMethodResponseMapper
}

type transactionStatsByCardMethodService struct {
	cache mencache.TransactionStatsByCardMethodCache

	errorHandler errorhandler.TransactionStatisticByCardErrorHandler

	repository repository.TransactionStatsByCardMethodRepository

	logger logger.LoggerInterface

	mapper responseservice.TransactionStatsMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransactionStatsByCardMethodService(params *transactionStatsByCardMethodServiceDeps) TransactionStatsByCardMethodService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_bycard_method_service_request_total",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_bycard_method_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-stats-bycard-method-service"), params.Logger, requestCounter, requestDuration)

	return &transactionStatsByCardMethodService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyPaymentMethodsByCardNumber retrieves monthly payment method usage by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number, year, and month.
//
// Returns:
//   - []*response.TransactionMonthMethodResponse: List of monthly payment method usage.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardMethodService) FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyPaymentMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyPaymentMethodsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly payment methods by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyPaymentMethodsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodsByCardNumberError(err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHODS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionMonthlyMethodResponses(records)

	s.cache.SetMonthlyPaymentMethodsByCardCache(ctx, req, responses)

	logSuccess("Successfully fetched monthly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

// FindYearlyPaymentMethodsByCardNumber retrieves yearly payment method usage by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains the card number and year.
//
// Returns:
//   - []*response.TransactionYearMethodResponse: List of yearly payment method usage.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionStatsByCardMethodService) FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {

	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyPaymentMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyPaymentMethodsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly payment methods by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyPaymentMethodsByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodsByCardNumberError(err, method, "FAILED_FIND_YEARLY_PAYMENT_METHODS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTransactionYearlyMethodResponses(records)

	s.cache.SetYearlyPaymentMethodsByCardCache(ctx, req, responses)

	logSuccess("Successfully fetched yearly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

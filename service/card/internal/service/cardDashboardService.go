package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/dashboard"
	repository "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/dashboard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// cardDashboardDeps holds the dependencies required to initialize the cardDashboardService.
// This struct is typically used for dependency injection when constructing the service.
type cardDashboardDeps struct {
	// ErrorHandler handles errors specific to card dashboard operations.
	ErrorHandler errorhandler.CardDashboardErrorHandler

	// Cache provides a caching mechanism for dashboard-related data.
	Cache mencache.CardDashboardCache

	// CardDashboardRepository provides access to dashboard-related data from the data store.
	CardDashboardRepository repository.CardDashboardRepository

	// Logger is used to log operational and error information.
	Logger logger.LoggerInterface
}

// cardDashboardService implements the CardDashboardService interface.
// It provides business logic and access to card dashboard data with support
// for caching, logging, metrics, and tracing.
type cardDashboardService struct {
	// errorhandler handles service-specific error formatting and propagation.
	errorhandler errorhandler.CardDashboardErrorHandler

	// mencache provides caching capabilities to reduce load on the repository.
	mencache mencache.CardDashboardCache

	// cardDashboardRepository interfaces with the database or data store.
	cardDashboardRepository repository.CardDashboardRepository

	// logger is used for structured and contextual logging.
	logger logger.LoggerInterface

	observability observability.TraceLoggerObservability
}

// NewCardDashboardService initializes a new instance of cardDashboardService.
//
// It sets up Prometheus metrics for tracking request counts and durations,
// and registers them for monitoring. This constructor function requires a set
// of parameters encapsulated in cardDashboardDeps, which include the context,
// error handler, cache, repository, logger, and mapper.
//
// Parameters:
//   - params: A pointer to cardDashboardDeps containing the dependencies
//     needed to initialize the service.
//
// Returns:
//   - A pointer to a cardDashboardService struct, fully initialized and ready to handle
//     dashboard card operations.
func NewCardDashboardService(
	params *cardDashboardDeps,
) CardDashboardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_dashboard_request_count",
		Help: "Number of card dashboard requests CardDashboardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_dashboard_request_duration_seconds",
		Help:    "Duration of card dashboard requests CardDashboardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-dashboard-service"), params.Logger, requestCounter, requestDuration)

	return &cardDashboardService{
		errorhandler:            params.ErrorHandler,
		mencache:                params.Cache,
		cardDashboardRepository: params.CardDashboardRepository,
		logger:                  params.Logger,
		observability:           observability,
	}
}

// DashboardCard retrieves the total balance, topup, withdraw, transaction, and transfer amounts.
//
// It first checks if the data is available in the cache. If it is, it returns the data.
// If not, it retrieves the data from the database.
//
// If any of the retrievals fail, it returns an error with the error message and the status code.
//
// If all retrievals succeed, it sets the data in the cache and returns the result.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//
// Returns:
// - *response.DashboardCard, *response.ErrorResponse
func (s *cardDashboardService) DashboardCard(ctx context.Context) (*response.DashboardCard, *response.ErrorResponse) {
	const method = "DashboardCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetDashboardCardCache(ctx); found {
		s.logger.Debug("DashboardCard cache hit")
		return data, nil
	}

	totalBalance, err := s.cardDashboardRepository.GetTotalBalances(ctx)
	if err != nil {
		return s.errorhandler.HandleTotalBalanceError(err, "DashboardCard", "FAILED_FIND_TOTAL_BALANCE", span, &status, zap.Error(err))
	}

	totalTopup, err := s.cardDashboardRepository.GetTotalTopAmount(ctx)
	if err != nil {
		return s.errorhandler.HandleTotalTopupAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_TOPUP", span, &status, zap.Error(err))
	}

	totalWithdraw, err := s.cardDashboardRepository.GetTotalWithdrawAmount(ctx)
	if err != nil {
		return s.errorhandler.HandleTotalWithdrawAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_WITHDRAW", span, &status, zap.Error(err))
	}

	totalTransaction, err := s.cardDashboardRepository.GetTotalTransactionAmount(ctx)
	if err != nil {
		return s.errorhandler.HandleTotalTransactionAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_TRANSACTION", span, &status, zap.Error(err))
	}

	totalTransfer, err := s.cardDashboardRepository.GetTotalTransferAmount(ctx)
	if err != nil {
		return s.errorhandler.HandleTotalTransferAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_TRANSFER", span, &status, zap.Error(err))
	}

	result := &response.DashboardCard{
		TotalBalance:     totalBalance,
		TotalTopup:       totalTopup,
		TotalWithdraw:    totalWithdraw,
		TotalTransaction: totalTransaction,
		TotalTransfer:    totalTransfer,
	}

	s.mencache.SetDashboardCardCache(ctx, result)

	logSuccess("Success find dashboard card", zap.Bool("success", true))

	return result, nil
}

// DashboardCardCardNumber retrieves the total balance, topup, withdraw, transaction, and transfer amounts of a specific card.
//
// It first checks if the data is available in the cache. If it is, it returns the data.
// If not, it retrieves the data from the database.
//
// If any of the retrievals fail, it returns an error with the error message and the status code.
//
// If all retrievals succeed, it sets the data in the cache and returns the result.
//
// Parameters:
// - ctx: The context for request-scoped values, cancellation, and deadlines.
// - cardNumber: The card number of the card to retrieve data for.
//
// Returns:
// - *response.DashboardCardCardNumber, *response.ErrorResponse
func (s *cardDashboardService) DashboardCardCardNumber(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	const method = "DashboardCardCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetDashboardCardCardNumberCache(ctx, cardNumber); found {
		s.logger.Debug("DashboardCardCardNumber cache hit", zap.String("card_number", cardNumber))
		return data, nil
	}
	s.logger.Debug("DashboardCardCardNumber cache miss", zap.String("card_number", cardNumber))

	totalBalance, err := s.cardDashboardRepository.GetTotalBalanceByCardNumber(ctx, cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalBalanceCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_BALANCE_BY_CARD", span, &status, zap.Error(err))
	}

	totalTopup, err := s.cardDashboardRepository.GetTotalTopupAmountByCardNumber(ctx, cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTopupAmountCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TOPUP_BY_CARD", span, &status, zap.Error(err))
	}

	totalWithdraw, err := s.cardDashboardRepository.GetTotalWithdrawAmountByCardNumber(ctx, cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalWithdrawAmountCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_WITHDRAW_BY_CARD", span, &status, zap.Error(err))
	}

	totalTransaction, err := s.cardDashboardRepository.GetTotalTransactionAmountByCardNumber(ctx, cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTransactionAmountCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TRANSACTION_BY_CARD", span, &status, zap.Error(err))
	}

	totalTransferSent, err := s.cardDashboardRepository.GetTotalTransferAmountBySender(ctx, cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTransferAmountBySender(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TRANSFER_BY_CARD", span, &status, zap.Error(err))
	}

	totalTransferReceived, err := s.cardDashboardRepository.GetTotalTransferAmountByReceiver(ctx, cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTransferAmountByReceiver(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TRANSFER_BY_CARD", span, &status, zap.Error(err))
	}

	result := &response.DashboardCardCardNumber{
		TotalBalance:          totalBalance,
		TotalTopup:            totalTopup,
		TotalWithdraw:         totalWithdraw,
		TotalTransaction:      totalTransaction,
		TotalTransferSend:     totalTransferSent,
		TotalTransferReceiver: totalTransferReceived,
	}

	s.mencache.SetDashboardCardCardNumberCache(ctx, cardNumber, result)

	logSuccess("Success find dashboard card card number", zap.Bool("success", true))

	return result, nil
}

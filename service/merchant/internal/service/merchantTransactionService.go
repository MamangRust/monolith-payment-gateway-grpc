package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
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

// merchantTransactionDeps contains the dependencies required to
// construct a merchantTransactionService.
type merchantTransactionDeps struct {
	// ErrorHandler handles errors related to merchant transactions.
	ErrorHandler errorhandler.MerchantTransactionErrorHandler

	// Repository provides access to merchant transaction data.
	Repository repository.MerchantTransactionRepository

	// Logger is used for structured logging.
	Logger logger.LoggerInterface

	// Mapper maps internal data to response formats.
	Mapper responseservice.MerchantTransactionResponseMapper

	Cache mencache.MerchantTransactionCache
}

// merchantTransactionService handles operations related to merchant transactions.
type merchantTransactionService struct {
	// errorHandler handles errors related to merchant transactions.
	errorHandler errorhandler.MerchantTransactionErrorHandler

	// merchantTransactionRepository provides access to merchant transaction data.
	merchantTransactionRepository repository.MerchantTransactionRepository

	// logger is used for structured logging.
	logger logger.LoggerInterface

	// mencache is the cache layer for merchant transaction queries.
	mencache mencache.MerchantTransactionCache

	// mapper maps internal data to response formats.
	mapper responseservice.MerchantTransactionResponseMapper

	observability observability.TraceLoggerObservability
}

// NewMerchantTransactionService initializes a new instance of merchantTransactionService with the
// provided parameters.
//
// It sets up Prometheus metrics for tracking request counts and durations and returns a
// configured merchantTransactionService ready for handling merchant transaction-related operations.
//
// Parameters:
// - params: A pointer to merchantTransactionDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created merchantTransactionService.
func NewMerchantTransactionService(
	params *merchantTransactionDeps,
) MerchantTransactionService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_transaction_service_requests_total",
			Help: "Total number of requests to the MerchantTransactionService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_transaction_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantTransactionService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-transaction-service"), params.Logger, requestCounter, requestDuration)

	return &merchantTransactionService{
		errorHandler:                  params.ErrorHandler,
		merchantTransactionRepository: params.Repository,
		logger:                        params.Logger,
		mapper:                        params.Mapper,
		observability:                 observability,
	}
}

// FindAllTransactions retrieves all merchant transactions with optional filters and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters for filters, search, and pagination.
//
// Returns:
//   - []*response.MerchantTransactionResponse: A list of transaction records.
//   - *int: The total number of matched transactions.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantTransactionService) FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllTransactions"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCacheAllMerchantTransactions(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactions(ctx, req)

	if err != nil {
		return s.errorHandler.HandleRepositoryAllError(
			err, method, "FAILED_FIND_ALL_TRANSACTIONS", span, &status, zap.Error(err),
		)
	}

	merchantResponses := s.mapper.ToMerchantsTransactionResponse(merchants)

	s.mencache.SetCacheAllMerchantTransactions(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

// FindAllTransactionsByMerchant retrieves all transactions for a specific merchant by merchant ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing the merchant ID and other filters.
//
// Returns:
//   - []*response.MerchantTransactionResponse: A list of transaction records.
//   - *int: The total number of matched transactions.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantTransactionService) FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllTransactionsByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCacheMerchantTransactions(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactionsByMerchant(ctx, req)

	if err != nil {
		return s.errorHandler.HandleRepositoryAllError(
			err, method, "FAILED_FIND_ALL_TRANSACTIONS_BY_MERCHANT", span, &status, zap.Error(err),
		)
	}

	merchantResponses := s.mapper.ToMerchantsTransactionResponse(merchants)

	s.mencache.SetCacheMerchantTransactions(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

// FindAllTransactionsByApikey retrieves all transactions for a merchant identified by an API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing the API key and other filters.
//
// Returns:
//   - []*response.MerchantTransactionResponse: A list of transaction records.
//   - *int: The total number of matched transactions.
//   - *response.ErrorResponse: An error returned if the retrieval fails.
func (s *merchantTransactionService) FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllTransactionsByApikey"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCacheMerchantTransactionApikey(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactionsByApikey(ctx, req)

	if err != nil {
		return s.errorHandler.HandleRepositoryAllError(
			err, method, "FAILED_FIND_ALL_TRANSACTIONS_BY_APIKEY", span, &status, zap.Error(err),
		)
	}

	merchantResponses := s.mapper.ToMerchantsTransactionResponse(merchants)

	s.mencache.SetCacheMerchantTransactionApikey(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

// normalizePagination validates and normalizes pagination parameters.
// Ensures the page is set to at least 1 and pageSize to a default of 10 if
// they are not positive. Returns the normalized page and pageSize values.
//
// Parameters:
//   - page: The requested page number.
//   - pageSize: The number of items per page.
//
// Returns:
//   - The normalized page and pageSize values.
func (s *merchantTransactionService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

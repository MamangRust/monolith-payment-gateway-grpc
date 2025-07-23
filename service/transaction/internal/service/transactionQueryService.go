package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// TransactionQueryServiceDeps defines the dependencies required to initialize a transactionQueryService.
type transactionQueryServiceDeps struct {
	// Ctx is the base context used for all service operations.
	Ctx context.Context

	// Mencache is the cache layer used for query-side transaction caching.
	Cache mencache.TransactionQueryCache

	// ErrorHandler handles structured errors and tracing for the query service.
	ErrorHandler errorhandler.TransactionQueryErrorHandler

	// TransactionQueryRepository is the repository used for querying transaction data from the database.
	TransactionQueryRepository repository.TransactionQueryRepository

	// Logger is the structured logger used for logging service activities and errors.
	Logger logger.LoggerInterface

	// Mapping provides mapper logic from domain transaction data to response DTOs.
	Mapper responseservice.TransactionQueryResponseMapper
}

// transactionQueryService handles all read/query operations related to transactions.
// It supports caching, tracing, logging, and Prometheus metrics.
type transactionQueryService struct {
	// ctx is the base context for operations.
	ctx context.Context

	// mencache is the Redis cache layer for transactions.
	mencache mencache.TransactionQueryCache

	// errorhandler handles error formatting, logging, and tracing.
	errorhandler errorhandler.TransactionQueryErrorHandler

	// transactionQueryRepository provides access to transaction data from DB.
	transactionQueryRepository repository.TransactionQueryRepository

	// logger logs structured events and errors.
	logger logger.LoggerInterface

	// mapper maps transaction data to response format.
	mapper responseservice.TransactionQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewTransactionQueryService creates a new transactionQueryService with the given dependencies.
// It returns the configured service and registers the provided Prometheus metrics.
//
// Parameters:
//   - ctx: The base context used for all service operations.
//   - mencache: The cache layer used for query-side transaction caching.
//   - errorhandler: Handles structured errors and tracing for the query service.
//   - transactionQueryRepository: The repository used for querying transaction data from the database.
//   - logger: The structured logger used for logging service activities and errors.
//   - mapper: Provides mapper logic from domain transaction data to response DTOs.
//
// Returns:
//   - A pointer to a transactionQueryService with the given dependencies.
func NewTransactionQueryService(
	params *transactionQueryServiceDeps,
) TransactionQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_query_service_request_total",
			Help: "Total number of requests to the TransactionQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-query-service"), params.Logger, requestCounter, requestDuration)

	return &transactionQueryService{
		ctx:                        params.Ctx,
		mencache:                   params.Cache,
		errorhandler:               params.ErrorHandler,
		transactionQueryRepository: params.TransactionQueryRepository,
		logger:                     params.Logger,
		mapper:                     params.Mapper,
		observability:              observability,
	}
}

// FindAll retrieves a paginated list of all transactions based on the given filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination info.
//
// Returns:
//   - []*response.TransactionResponse: List of transactions.
//   - *int: Total number of transactions.
//   - *response.ErrorResponse: Error response if query fails.
func (s *transactionQueryService) FindAll(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionsCache(ctx, req); found {
		logSuccess("Successfully retrieved all transaction records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionsCache(ctx, req, responseData, totalRecords)

	logSuccess("Successfully retrieved all transaction records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

// FindAllByCardNumber retrieves all transactions associated with a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number and pagination.
//
// Returns:
//   - []*response.TransactionResponse: List of transactions for the card number.
//   - *int: Total number of transactions.
//   - *response.ErrorResponse: Error response if query fails.
func (s *transactionQueryService) FindAllByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	cardNumber := req.CardNumber

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.String("cardNumber", cardNumber))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionByCardNumberCache(ctx, req); found {
		logSuccess("Successfully retrieved all transaction records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactionByCardNumber(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_BYCARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionByCardNumberCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved all transaction records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindById retrieves a transaction by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transactionID: The ID of the transaction to retrieve.
//
// Returns:
//   - *response.TransactionResponse: The transaction data.
//   - *response.ErrorResponse: Error response if query fails.
func (s *transactionQueryService) FindById(ctx context.Context, transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransactionCache(ctx, transactionID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("transaction.id", transactionID))
		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindById(ctx, transactionID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION", span, &status, transaction_errors.ErrTransactionNotFound, zap.Error(err))
	}

	so := s.mapper.ToTransactionResponse(transaction)

	s.mencache.SetCachedTransactionCache(ctx, so)

	logSuccess("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return so, nil
}

// FindByActive retrieves active transactions with soft-delete not applied.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination info.
//
// Returns:
//   - []*response.TransactionResponseDeleteAt: List of active transactions.
//   - *int: Total number of active transactions.
//   - *response.ErrorResponse: Error response if query fails.
func (s *transactionQueryService) FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active transaction from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE", span, &status, transaction_errors.ErrFailedFindByActiveTransactions, zap.Error(err))
	}

	so := s.mapper.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active transaction", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByTrashed retrieves transactions that have been soft-deleted.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination info.
//
// Returns:
//   - []*response.TransactionResponseDeleteAt: List of trashed transactions.
//   - *int: Total number of trashed transactions.
//   - *response.ErrorResponse: Error response if query fails.
func (s *transactionQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed transaction from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED", span, &status, transaction_errors.ErrFailedFindByTrashedTransactions, zap.Error(err))
	}

	so := s.mapper.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed transaction", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindTransactionByMerchantId retrieves transactions by merchant ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: The ID of the merchant whose transactions are requested.
//
// Returns:
//   - []*response.TransactionResponse: List of transactions for the merchant.
//   - *response.ErrorResponse: Error response if query fails.
func (s *transactionQueryService) FindTransactionByMerchantId(ctx context.Context, merchantID int) ([]*response.TransactionResponse, *response.ErrorResponse) {
	const method = "FindTransactionByMerchantId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransactionByMerchantIdCache(ctx, merchantID); found {
		logSuccess("Successfully fetched transaction by merchant ID from cache", zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionQueryRepository.FindTransactionByMerchantId(ctx, merchantID)
	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, method, "FAILED_FIND_TRANSACTION_BY_MERCHANT_ID", span, &status, transaction_errors.ErrFailedFindByMerchantID, zap.Error(err))
	}

	responseData := s.mapper.ToTransactionsResponse(res)

	s.mencache.SetCachedTransactionByMerchantIdCache(ctx, merchantID, responseData)

	logSuccess("Successfully fetched transaction by merchant ID", zap.Int("merchant.id", merchantID))

	return responseData, nil
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
func (s *transactionQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

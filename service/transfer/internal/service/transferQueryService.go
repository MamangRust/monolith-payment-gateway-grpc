package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"

	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// transferQueryDeps holds the dependencies required to create a transferQueryService.
type transferQueryDeps struct {
	// Ctx is the context used across the service for deadlines, cancelation, and tracing.
	Ctx context.Context

	// ErrorHandler handles errors that occur during transfer query operations.
	ErrorHandler errorhandler.TransferQueryErrorHandler

	// Cache provides caching functionality for transfer queries.
	Cache mencache.TransferQueryCache

	// Repository is the data access layer for transfer query operations.
	Repository repository.TransferQueryRepository

	// Logger provides structured logging throughout the service.
	Logger logger.LoggerInterface

	// Mapper converts internal data models to response DTOs.
	Mapper responseservice.TransferQueryResponseMapper
}

// transferQueryService provides a service for querying transfer records.
//
// The service provides methods for retrieving transfer records by ID, finding
// transfer records by year and month, and retrieving the monthly transfer status
// for successful transactions.
type transferQueryService struct {
	// ctx is the context for the service.
	ctx context.Context

	// errorHandler is the error handler for the service.
	errorHandler errorhandler.TransferQueryErrorHandler

	// mencache is the cache for the service.
	mencache mencache.TransferQueryCache

	// transferQueryRepository is the repository for the service.
	transferQueryRepository repository.TransferQueryRepository

	// logger is the logger for the service.
	logger logger.LoggerInterface

	// mapper is the mapper for the service.
	mapper responseservice.TransferQueryResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransferQueryService(
	params *transferQueryDeps,
) TransferQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_query_service_request_total",
			Help: "Total number of requests to the TransferQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transfer-query-service"), params.Logger, requestCounter, requestDuration)

	return &transferQueryService{
		ctx:                     params.Ctx,
		errorHandler:            params.ErrorHandler,
		mencache:                params.Cache,
		transferQueryRepository: params.Repository,
		logger:                  params.Logger,
		mapper:                  params.Mapper,
		observability:           observability,
	}
}

// FindAll retrieves all transfer records based on filter and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The filter and pagination parameters.
//
// Returns:
//   - []*response.TransferResponse: List of transfer responses.
//   - *int: Total number of transfer records.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferQueryService) FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransfersCache(ctx, req); found {
		logSuccess("Successfully retrieved all transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindAll(ctx, req)

	if err != nil {
		return s.errorHandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_ALL_TRANSFERS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransfersResponse(transfers)

	s.mencache.SetCachedTransfersCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved all transfer records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindById retrieves a single transfer record by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transferId: The ID of the transfer to retrieve.
//
// Returns:
//   - *response.TransferResponse: The transfer response.
//   - *response.ErrorResponse: Error response if not found or failed.
func (s *transferQueryService) FindById(ctx context.Context, transferId int) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transfer.id", transferId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransferCache(ctx, transferId); found {
		logSuccess("Successfully fetched transfer from cache", zap.Int("transfer.id", transferId))
		return data, nil
	}

	transfer, err := s.transferQueryRepository.FindById(ctx, transferId)

	if err != nil {
		return s.errorHandler.HandleRepositorySingleError(err, method, "FAILED_TO_FIND_TRANSFER_BY_ID", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	so := s.mapper.ToTransferResponse(transfer)

	s.mencache.SetCachedTransferCache(ctx, so)

	logSuccess("Successfully fetched transfer", zap.Int("transfer.id", transferId))

	return so, nil
}

// FindByActive retrieves all active (non-deleted) transfer records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The filter and pagination parameters.
//
// Returns:
//   - []*response.TransferResponseDeleteAt: List of active transfer responses with deleted_at info.
//   - *int: Total number of active records.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferQueryService) FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransferActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorHandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_BY_ACTIVE_TRANSFERS", span, &status, transfer_errors.ErrFailedFindActiveTransfers, zap.Error(err))
	}

	so := s.mapper.ToTransfersResponseDeleteAt(transfers)

	s.mencache.SetCachedTransferActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved active transfer records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByTrashed retrieves all trashed (soft-deleted) transfer records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The filter and pagination parameters.
//
// Returns:
//   - []*response.TransferResponseDeleteAt: List of trashed transfer responses with deleted_at info.
//   - *int: Total number of trashed records.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransferTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorHandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_BY_TRASHED_TRANSFERS", span, &status, transfer_errors.ErrFailedFindTrashedTransfers, zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	}

	so := s.mapper.ToTransfersResponseDeleteAt(transfers)

	s.mencache.SetCachedTransferTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved trashed transfer records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindTransferByTransferFrom retrieves transfers by sender card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_from: The sender card number.
//
// Returns:
//   - []*response.TransferResponse: List of transfer responses.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferQueryService) FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse) {
	const method = "FindTransferByTransferFrom"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("transaction.from", transfer_from))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransferByFrom(ctx, transfer_from); found {
		logSuccess("Successfully fetched transfer from cache", zap.String("transfer_from", transfer_from))
		return data, nil
	}

	res, err := s.transferQueryRepository.FindTransferByTransferFrom(ctx, transfer_from)

	if err != nil {
		return s.errorHandler.HanldeRepositoryListError(err, method, "FAILED_TO_FIND_TRANSFER_BY_TRANSFER_FROM", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	so := s.mapper.ToTransfersResponse(res)

	s.mencache.SetCachedTransferByFrom(ctx, transfer_from, so)

	logSuccess("Successfully fetched transfer", zap.String("transfer_from", transfer_from))

	return so, nil
}

// FindTransferByTransferTo retrieves transfers by receiver card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_to: The receiver card number.
//
// Returns:
//   - []*response.TransferResponse: List of transfer responses.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferQueryService) FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse) {
	const method = "FindTransferByTransferTo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("transfer.to", transfer_to))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransferByTo(ctx, transfer_to); found {
		logSuccess("Successfully fetched transfer from cache", zap.String("transfer_to", transfer_to))
		return data, nil
	}

	res, err := s.transferQueryRepository.FindTransferByTransferTo(ctx, transfer_to)

	if err != nil {
		return s.errorHandler.HanldeRepositoryListError(err, method, "FAILED_TO_FIND_TRANSFER_BY_TRANSFER_TO", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	so := s.mapper.ToTransfersResponse(res)

	s.mencache.SetCachedTransferByTo(ctx, transfer_to, so)

	logSuccess("Successfully fetched transfer", zap.String("transfer_to", transfer_to))

	return so, nil
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
func (s *transferQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

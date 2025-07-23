package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// merchantDocumentQueryDeps contains the dependencies required to
// construct a merchantDocumentQueryService.
type merchantDocumentQueryDeps struct {
	// Cache is the cache layer for merchant document queries.
	Cache mencache.MerchantDocumentQueryCache

	// ErrorHandler handles errors for merchant document queries.
	ErrorHandler errorhandler.MerchantDocumentQueryErrorHandler

	// Repository provides access to merchant document query data.
	Repository repository.MerchantDocumentQueryRepository

	// Logger is used for structured logging.
	Logger logger.LoggerInterface

	// Mapper maps internal data to response formats.
	Mapper responseservice.MerchantDocumentQueryResponseMapper
}

// merchantDocumentQueryService handles merchant document read/query operations.
type merchantDocumentQueryService struct {
	// mencache is the cache layer for merchant document queries.
	mencache mencache.MerchantDocumentQueryCache

	// errorhandler handles errors for merchant document queries.
	errorhandler errorhandler.MerchantDocumentQueryErrorHandler

	// merchantDocumentQueryRepository provides access to merchant document query data.
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository

	// logger is used for logging within the service.
	logger logger.LoggerInterface

	// mapper maps internal data to response formats.
	mapper responseservice.MerchantDocumentQueryResponseMapper

	observability observability.TraceLoggerObservability
}

func NewMerchantDocumentQueryService(
	params *merchantDocumentQueryDeps,
) MerchantDocumentQueryService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_document_query_request_count",
		Help: "Number of merchant document query requests MerchantDocumentQueryService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_document_query_request_duration_seconds",
		Help:    "The duration of requests MerchantDocumentQueryService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-document-query-service"), params.Logger, requestCounter, requestDuration)

	return &merchantDocumentQueryService{
		mencache:                        params.Cache,
		errorhandler:                    params.ErrorHandler,
		merchantDocumentQueryRepository: params.Repository,
		logger:                          params.Logger,
		mapper:                          params.Mapper,
		observability:                   observability,
	}
}

// FindAll retrieves all merchant documents with optional filtering and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters including filters and pagination options.
//
// Returns:
//   - []*response.MerchantDocumentResponse: A list of merchant documents.
//   - *int: The total count of matched documents.
//   - *response.ErrorResponse: An error if retrieval fails.
func (s *merchantDocumentQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocuments(ctx, req); found {
		logSuccess("Successfully retrieved all merchant document records from cache", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindAllDocuments(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENTS", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapper.ToMerchantDocumentsResponse(merchantDocuments)

	s.mencache.SetCachedMerchantDocuments(ctx, req, merchantResponse, total)

	logSuccess("Successfully retrieved all merchant document records", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

// FindById retrieves a merchant document by its unique document ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - document_id: The ID of the document to retrieve.
//
// Returns:
//   - *response.MerchantDocumentResponse: The document data.
//   - *response.ErrorResponse: An error if the document is not found or retrieval fails.
func (s *merchantDocumentQueryService) FindById(ctx context.Context, merchant_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", merchant_id))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetCachedMerchantDocument(ctx, merchant_id); found {
		logSuccess("Successfully found merchant document by ID from cache", zap.Int("merchantDocument.id", merchant_id))
		return cachedMerchant, nil
	}

	merchantDocument, err := s.merchantDocumentQueryRepository.FindByIdDocument(ctx, merchant_id)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_DOCUMENT_BY_ID", span, &status, merchantdocument_errors.ErrFailedFindMerchantDocumentById, zap.Error(err))

	}

	merchantResponse := s.mapper.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.SetCachedMerchantDocument(ctx, merchant_id, merchantResponse)

	logSuccess("Successfully found merchant document by ID", zap.Int("merchantDocument.id", merchant_id))

	return merchantResponse, nil
}

// FindByActive retrieves all active (non-deleted) merchant documents.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters including filters and pagination options.
//
// Returns:
//   - []*response.MerchantDocumentResponse: A list of active merchant documents.
//   - *int: The total count of matched active documents.
//   - *response.ErrorResponse: An error if retrieval fails.
func (s *merchantDocumentQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocuments(ctx, req); found {
		logSuccess("Successfully retrieved active merchant document records from cache", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByActiveDocuments(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENTS_ACTIVE", span, &status)
	}

	merchantResponse := s.mapper.ToMerchantDocumentsResponse(merchantDocuments)

	s.mencache.SetCachedMerchantDocuments(ctx, req, merchantResponse, total)

	logSuccess("Successfully retrieved active merchant document records", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

// FindByTrashed retrieves all trashed (soft-deleted) merchant documents.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters including filters and pagination options.
//
// Returns:
//   - []*response.MerchantDocumentResponseDeleteAt: A list of trashed merchant documents.
//   - *int: The total count of matched trashed documents.
//   - *response.ErrorResponse: An error if retrieval fails.
func (s *merchantDocumentQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocumentsTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchant document records from cache", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByTrashedDocuments(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENTS_TRASHED", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapper.ToMerchantDocumentsResponseDeleteAt(merchantDocuments)

	s.mencache.SetCachedMerchantDocumentsTrashed(ctx, req, merchantResponse, total)

	logSuccess("Successfully retrieved trashed merchant document records", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
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
func (s *merchantDocumentQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

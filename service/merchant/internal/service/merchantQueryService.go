package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// merchantQueryDeps contains the dependencies required to
// construct a merchantQueryService.
type merchantQueryDeps struct {
	// Repository provides access to merchant query data.
	Repository repository.MerchantQueryRepository

	// ErrorHandler handles errors related to merchant queries.
	ErrorHandler errorhandler.MerchantQueryErrorHandler

	// Cache is the cache layer for merchant query results.
	Cache mencache.MerchantQueryCache

	// Logger is used for structured logging.
	Logger logger.LoggerInterface

	// Mapper maps internal data to response formats.
	Mapper responseservice.MerchantQueryResponseMapper
}

// merchantQueryService handles merchant read/query operations.
type merchantQueryService struct {
	// merchantQueryRepository provides access to merchant query data.
	merchantQueryRepository repository.MerchantQueryRepository

	// logger is used for structured logging.
	logger logger.LoggerInterface

	// mapper maps internal data to response formats.
	mapper responseservice.MerchantQueryResponseMapper

	// errorhandler handles errors related to merchant queries.
	errorhandler errorhandler.MerchantQueryErrorHandler

	// mencache is the cache layer for merchant query results.
	mencache mencache.MerchantQueryCache

	observability observability.TraceLoggerObservability
}

// NewMerchantQueryService initializes a new instance of merchantQueryService with the
// provided parameters.
//
// It sets up Prometheus metrics for tracking request counts and durations and returns a
// configured merchantQueryService ready for handling merchant queries.
//
// Parameters:
// - params: A pointer to merchantQueryDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created merchantQueryService.
func NewMerchantQueryService(params *merchantQueryDeps) MerchantQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_query_service_requests_total",
			Help: "Total number of requests to the MerchantQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-query-service"), params.Logger, requestCounter, requestDuration)

	return &merchantQueryService{
		merchantQueryRepository: params.Repository,
		logger:                  params.Logger,
		mapper:                  params.Mapper,
		errorhandler:            params.ErrorHandler,
		mencache:                params.Cache,
		observability:           observability,
	}
}

// FindAll retrieves a list of merchants with pagination and optional filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters for filtering and pagination.
//
// Returns:
//   - []*response.MerchantResponse: A list of merchant data.
//   - *int: The total count of matched merchants.
//   - *response.ErrorResponse: An error if the operation fails.
func (s *merchantQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchants(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANTS", span, &status, zap.Error(err))
	}

	merchantResponses := s.mapper.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchants(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

// FindById retrieves a merchant by its unique merchant ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: The ID of the merchant to be retrieved.
//
// Returns:
//   - *response.MerchantResponse: The merchant data.
//   - *response.ErrorResponse: An error if the merchant is not found or retrieval fails.
func (s *merchantQueryService) FindById(ctx context.Context, merchantID int) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetCachedMerchant(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant from cache", zap.Int("merchant.id", merchantID))
		return cachedMerchant, nil
	}

	merchant, err := s.merchantQueryRepository.FindById(ctx, merchantID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	merchantResponse := s.mapper.ToMerchantResponse(merchant)

	s.mencache.SetCachedMerchant(ctx, merchantResponse)

	logSuccess("Successfully retrieved merchant", zap.Int("merchant.id", merchantID))

	return merchantResponse, nil
}

// FindByActive retrieves all active (non-deleted) merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters for filtering and pagination.
//
// Returns:
//   - []*response.MerchantResponseDeleteAt: A list of active merchant records.
//   - *int: The total count of matched active merchants.
//   - *response.ErrorResponse: An error if the operation fails.
func (s *merchantQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_MERCHANTS", span, &status, merchant_errors.ErrFailedFindActiveMerchants,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapper.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantActive(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully fetched active merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

// FindByTrashed retrieves all soft-deleted (trashed) merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters for filtering and pagination.
//
// Returns:
//   - []*response.MerchantResponseDeleteAt: A list of trashed merchant records.
//   - *int: The total count of matched trashed merchants.
//   - *response.ErrorResponse: An error if the operation fails.
func (s *merchantQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_MERCHANTS", span, &status, merchant_errors.ErrFailedFindTrashedMerchants,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapper.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantTrashed(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully fetched trashed merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

// FindByApiKey retrieves a merchant based on the provided API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - api_key: The API key associated with the merchant.
//
// Returns:
//   - *response.MerchantResponse: The matched merchant.
//   - *response.ErrorResponse: An error if not found or retrieval fails.
func (s *merchantQueryService) FindByApiKey(ctx context.Context, apiKey string) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindByApiKey"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("api_key", apiKey))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetCachedMerchantByApiKey(ctx, apiKey); found {
		logSuccess("Successfully found merchant by API key from cache", zap.String("api_key", apiKey))
		return cachedMerchant, nil
	}

	merchant, err := s.merchantQueryRepository.FindByApiKey(ctx, apiKey)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Error(err))
	}

	merchantResponse := s.mapper.ToMerchantResponse(merchant)

	s.mencache.SetCachedMerchantByApiKey(ctx, apiKey, merchantResponse)

	logSuccess("Successfully found merchant by API key", zap.String("api_key", apiKey))

	return merchantResponse, nil
}

// FindByMerchantUserId retrieves merchants associated with a given user ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user who owns the merchants.
//
// Returns:
//   - []*response.MerchantResponse: A list of merchants owned by the user.
//   - *response.ErrorResponse: An error if retrieval fails.
func (s *merchantQueryService) FindByMerchantUserId(ctx context.Context, userID int) ([]*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindByMerchantUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", userID))

	defer func() {
		end(status)
	}()

	if cachedMerchants, found := s.mencache.GetCachedMerchantsByUserId(ctx, userID); found {
		logSuccess("Successfully found merchants by user ID from cache", zap.Int("user.id", userID), zap.Int("count", len(cachedMerchants)))

		return cachedMerchants, nil
	}

	merchants, err := s.merchantQueryRepository.FindByMerchantUserId(ctx, userID)
	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, method, "FAILED_FIND_MERCHANT_BY_USER_ID", span, &status, zap.Error(err))

	}

	merchantResponses := s.mapper.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchantsByUserId(ctx, userID, merchantResponses)

	logSuccess("Successfully found merchants by user ID", zap.Int("user.id", userID), zap.Int("count", len(merchantResponses)))

	return merchantResponses, nil
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
func (s *merchantQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

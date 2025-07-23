package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// topupQueryDeps holds the dependencies required to construct a topupQueryService.
type topupQueryDeps struct {
	// Ctx is the context used for request-scoped operations, cancellations, and timeouts.
	Ctx context.Context

	// ErrorHandler handles domain-specific errors during top-up query operations.
	ErrorHandler errorhandler.TopupQueryErrorHandler

	// Cache provides caching capabilities for frequently accessed top-up query data.
	Cache mencache.TopupQueryCache

	// Repository provides access to read-only top-up data from the database.
	Repository repository.TopupQueryRepository

	// Logger provides structured logging support for debugging and monitoring.
	Logger logger.LoggerInterface

	// Mapper transforms internal domain data to top-up response DTOs.
	Mapper responseservice.TopupQueryResponseMapper
}

// topupQueryService handles the logic for querying top-up operations.
// It integrates caching, error handling, tracing, and metrics tracking.
type topupQueryService struct {
	// ctx is the context used for request-scoped operations, cancellations, and timeouts.
	ctx context.Context

	// errorhandler handles domain-specific errors during top-up query operations.
	errorhandler errorhandler.TopupQueryErrorHandler

	// mencache provides caching capabilities for frequently accessed top-up query data.
	mencache mencache.TopupQueryCache

	// topupQueryRepository provides access to read-only top-up data from the database.
	topupQueryRepository repository.TopupQueryRepository

	// logger provides structured logging support for debugging and monitoring.
	logger logger.LoggerInterface

	// mapper transforms internal domain data to top-up response DTOs.
	mapper responseservice.TopupQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewTopupQueryService initializes a new instance of topupQueryService with the provided parameters.
//
// It sets up Prometheus metrics for counting and measuring the duration of top-up query requests.
//
// Parameters:
// - params: A pointer to topupQueryDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created topupQueryService.
func NewTopupQueryService(
	params *topupQueryDeps,
) TopupQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_query_service_request_total",
			Help: "Total number of requests to the TopupQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("topup-query-service"), params.Logger, requestCounter, requestDuration)

	return &topupQueryService{
		ctx:                  params.Ctx,
		errorhandler:         params.ErrorHandler,
		mencache:             params.Cache,
		topupQueryRepository: params.Repository,
		logger:               params.Logger,
		mapper:               params.Mapper,
		observability:        observability,
	}
}

// FindAll retrieves all topup records based on the given filter request.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination criteria.
//
// Returns:
//   - []*response.TopupResponse: List of topup responses.
//   - *int: Total number of matching records.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupQueryService) FindAll(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTopupsCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopups(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_TOPUPS", span, &status, topup_errors.ErrFailedFindAllTopups, zap.Error(err))
	}

	so := s.mapper.ToTopupResponses(topups)

	s.mencache.SetCachedTopupsCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindAllByCardNumber retrieves all topup records filtered by a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number and filter criteria.
//
// Returns:
//   - []*response.TopupResponse: List of topup responses.
//   - *int: Total number of matching records.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupQueryService) FindAllByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, *response.ErrorResponse) {

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	cardNumber := req.CardNumber

	const method = "FindAllByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.String("cardNumber", cardNumber))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCacheTopupByCardCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopupByCardNumber(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_TOPUPS_BY_CARD", span, &status, topup_errors.ErrFailedFindAllTopupsByCardNumber)
	}

	so := s.mapper.ToTopupResponses(topups)

	s.mencache.SetCacheTopupByCardCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindById retrieves a topup record by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topupID: The ID of the topup to be retrieved.
//
// Returns:
//   - *response.TopupResponse: The topup data if found.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupQueryService) FindById(ctx context.Context, topupID int) (*response.TopupResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("topup.id", topupID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTopupCache(ctx, topupID); found {
		logSuccess("Successfully retrieved topup from cache", zap.Int("topup.id", topupID))
		return data, nil
	}

	topup, err := s.topupQueryRepository.FindById(ctx, topupID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_TOPUP", span, &status, topup_errors.ErrFailedFindTopupById, zap.Error(err))
	}

	so := s.mapper.ToTopupResponse(topup)

	s.mencache.SetCachedTopupCache(ctx, so)

	logSuccess("Successfully retrieved topup", zap.Int("topup.id", topupID))

	return so, nil
}

// FindByActive retrieves all active (non-deleted) topup records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination criteria.
//
// Returns:
//   - []*response.TopupResponseDeleteAt: List of active topup records.
//   - *int: Total number of matching records.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupQueryService) FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTopupActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_TOPUPS", span, &status, topup_errors.ErrFailedFindActiveTopups)
	}

	so := s.mapper.ToTopupResponsesDeleteAt(topups)

	s.mencache.SetCachedTopupActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByTrashed retrieves all soft-deleted (trashed) topup records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination criteria.
//
// Returns:
//   - []*response.TopupResponseDeleteAt: List of trashed topup records.
//   - *int: Total number of matching records.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTopupTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_TOPUPS", span, &status, topup_errors.ErrFailedFindTrashedTopups)
	}

	so := s.mapper.ToTopupResponsesDeleteAt(topups)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
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
func (s *topupQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// WithdrawQueryServiceDeps holds the dependencies for WithdrawQueryService.
type withdrawQueryServiceDeps struct {
	// Ctx is the request context for timeout and cancellation control.
	Ctx context.Context

	// ErrorHandler handles query-related errors for withdraw operations.
	ErrorHandler errorhandler.WithdrawQueryErrorHandler

	// Cache provides cached responses for withdraw query operations.
	Cache mencache.WithdrawQueryCache

	// Repository is used to read withdraw data from the database.
	Repository repository.WithdrawQueryRepository

	// Logger is used for logging query operations.
	Logger logger.LoggerInterface
	// Mapper maps domain data to API response objects.
	Mapper responseservice.WithdrawQueryResponseMapper
}

// withdrawQueryService implements the WithdrawQueryService interface.
type withdrawQueryService struct {
	// ErrorHandler handles query-related errors for withdraw operations.
	errorhandler errorhandler.WithdrawQueryErrorHandler

	// Cache provides cached responses for withdraw query operations.
	mencache mencache.WithdrawQueryCache

	// WithdrawQueryRepository is used to read withdraw data from the database.
	withdrawQueryRepository repository.WithdrawQueryRepository

	// Logger is used for logging query operations.
	logger logger.LoggerInterface

	// Mapper maps domain data to API response objects.
	mapper responseservice.WithdrawQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewWithdrawQueryService initializes a new withdrawQueryService with the provided parameters.
//
// It sets up the prometheus metrics for counting and measuring the duration of withdraw query requests.
//
// Parameters:
// - deps: A pointer to withdrawQueryServiceDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created withdrawQueryService.
func NewWithdrawQueryService(
	deps *withdrawQueryServiceDeps,
) WithdrawQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_query_service_request_total",
			Help: "Total number of requests to the WithdrawQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("withdraw-query-service"), deps.Logger, requestCounter, requestDuration)

	return &withdrawQueryService{
		mencache:                deps.Cache,
		errorhandler:            deps.ErrorHandler,
		withdrawQueryRepository: deps.Repository,
		logger:                  deps.Logger,
		mapper:                  deps.Mapper,
		observability:           observability,
	}
}

// FindAll retrieves all withdraws based on the given request filter.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filters like pagination, date range, etc.
//
// Returns:
//   - []*response.WithdrawResponse: List of withdraws.
//   - *int: Total number of records matching the filter.
//   - *response.ErrorResponse: Error details if any.
func (s *withdrawQueryService) FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedWithdrawsCache(ctx, req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAll(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponse := s.mapper.ToWithdrawsResponse(withdraws)

	s.mencache.SetCachedWithdrawsCache(ctx, req, withdrawResponse, totalRecords)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

// FindAllByCardNumber retrieves all withdraws filtered by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number and other filters.
//
// Returns:
//   - []*response.WithdrawResponse: List of withdraws for the specified card.
//   - *int: Total number of records found.
//   - *response.ErrorResponse: Error details if any.
func (s *withdrawQueryService) FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedWithdrawByCardCache(ctx, req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}
	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAllByCardNumber(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_WITHDRAW_BY_CARD", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponse := s.mapper.ToWithdrawsResponse(withdraws)

	s.mencache.SetCachedWithdrawByCardCache(ctx, req, withdrawResponse, totalRecords)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

// FindById retrieves a single withdraw record by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - withdrawID: The ID of the withdraw to retrieve.
//
// Returns:
//   - *response.WithdrawResponse: The withdraw data if found.
//   - *response.ErrorResponse: Error details if any.
func (s *withdrawQueryService) FindById(ctx context.Context, withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("withdraw.id", withdrawID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedWithdrawCache(ctx, withdrawID); found {
		logSuccess("Successfully retrieved withdraw from cache", zap.Int("withdraw_id", withdrawID))
		return data, nil
	}

	withdraw, err := s.withdrawQueryRepository.FindById(ctx, withdrawID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_WITHDRAW", span, &status, withdraw_errors.ErrWithdrawNotFound, zap.Error(err))
	}

	withdrawResponse := s.mapper.ToWithdrawResponse(withdraw)

	s.mencache.SetCachedWithdrawCache(ctx, withdrawResponse)

	logSuccess("Successfully retrieved withdraw", zap.Int("withdraw.id", withdrawID))

	return withdrawResponse, nil
}

// FindByActive retrieves active withdraw records based on the request filter.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filters like pagination, etc.
//
// Returns:
//   - []*response.WithdrawResponseDeleteAt: List of active withdraws.
//   - *int: Total number of active records found.
//   - *response.ErrorResponse: Error details if any.
func (s *withdrawQueryService) FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponses := s.mapper.ToWithdrawsResponseDeleteAt(withdraws)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

// FindByTrashed retrieves soft-deleted withdraw records based on the request filter.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filters like pagination, etc.
//
// Returns:
//   - []*response.WithdrawResponseDeleteAt: List of trashed withdraws.
//   - *int: Total number of trashed records found.
//   - *response.ErrorResponse: Error details if any.
func (s *withdrawQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponses := s.mapper.ToWithdrawsResponseDeleteAt(withdraws)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
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
func (s *withdrawQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

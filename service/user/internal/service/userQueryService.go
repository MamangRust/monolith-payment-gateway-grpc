package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/user"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// UserQueryDeps defines the required dependencies to initialize userQueryService.
type userQueryDeps struct {
	// Ctx is the context used across service operations.
	Ctx context.Context

	// ErrorHandler handles domain-level errors for user queries.
	ErrorHandler errorhandler.UserQueryError

	// Cache provides caching layer for user query results.
	Cache mencache.UserQueryCache

	// Repository provides access to user query data in the persistence layer.
	Repository repository.UserQueryRepository

	// Logger is used to log service activities and errors.
	Logger logger.LoggerInterface

	// Mapper maps internal user entities to response models.
	Mapper responseservice.UserQueryResponseMapper
}

// userQueryService provides methods to query user data with caching, tracing, and metrics support.
type userQueryService struct {
	// ctx is the context passed throughout service operations.
	ctx context.Context

	// errorhandler handles user query-related errors.
	errorhandler errorhandler.UserQueryError

	// mencache provides caching for query results.
	mencache mencache.UserQueryCache

	// userQueryRepository is the repository used to access user data.
	userQueryRepository repository.UserQueryRepository

	// logger is used for logging service events and errors.
	logger logger.LoggerInterface

	// mapper maps internal data to response models.
	mapper responseservice.UserQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewUserQueryService initializes a new instance of userQueryService with the provided parameters.
//
// It sets up the prometheus metrics for counting and measuring the duration of user query requests.
//
// Parameters:
// - params: A pointer to userQueryDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created userQueryService.
func NewUserQueryService(
	params *userQueryDeps,
) UserQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_query_service_request_total",
			Help: "Total number of requests to the UserQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the UserQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("user-query-service"), params.Logger, requestCounter, requestDuration)

	return &userQueryService{
		ctx:                 params.Ctx,
		errorhandler:        params.ErrorHandler,
		mencache:            params.Cache,
		userQueryRepository: params.Repository,
		logger:              params.Logger,
		mapper:              params.Mapper,
		observability:       observability,
	}
}

// FindAll retrieves all users based on the given request filter.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter parameters for users.
//
// Returns:
//   - []*response.UserResponse: List of user data.
//   - *int: Total count of users.
//   - *response.ErrorResponse: Error response if query fails.
func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUsersCache(ctx, req); found {
		logSuccess("Successfully retrieved all user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_USERS", span, &status, zap.Error(err))
	}

	userResponses := s.mapper.ToUsersResponse(users)

	s.mencache.SetCachedUsersCache(ctx, req, userResponses, totalRecords)

	logSuccess("Successfully retrieved all user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return userResponses, totalRecords, nil
}

// FindByID retrieves a specific user by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the user.
//
// Returns:
//   - *response.UserResponse: The user data.
//   - *response.ErrorResponse: Error response if retrieval fails.
func (s *userQueryService) FindByID(ctx context.Context, id int) (*response.UserResponse, *response.ErrorResponse) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedUserCache(ctx, id); found {
		logSuccess("Successfully retrieved user record from cache", zap.Int("user.id", id))
		return data, nil
	}

	user, err := s.userQueryRepository.FindById(ctx, id)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	userRes := s.mapper.ToUserResponse(user)

	s.mencache.SetCachedUserCache(ctx, userRes)

	logSuccess("Successfully retrieved user record", zap.Int("user.id", id))

	return userRes, nil
}

// FindByActive retrieves all active users (not soft-deleted).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter parameters.
//
// Returns:
//   - []*response.UserResponseDeleteAt: List of active user data.
//   - *int: Total count of active users.
//   - *response.ErrorResponse: Error response if query fails.
func (s *userQueryService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_USERS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved active user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByTrashed retrieves all soft-deleted (trashed) users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter parameters.
//
// Returns:
//   - []*response.UserResponseDeleteAt: List of trashed user data.
//   - *int: Total count of trashed users.
//   - *response.ErrorResponse: Error response if query fails.
func (s *userQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserTrashedCache(ctx, req); found {
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_USERS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToUsersResponseDeleteAt(users)

	logSuccess("Successfully retrieved trashed user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

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
func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

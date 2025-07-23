package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-role/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/service"
	roleservicemapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/role"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// roleQueryDeps contains the dependencies required to create a roleQueryService.
type roleQueryDeps struct {
	// Ctx is the base context for the service.
	Ctx context.Context

	// ErrorHandler handles domain-specific errors for role queries.
	ErrorHandler errorhandler.RoleQueryErrorHandler

	// Cache provides in-memory caching for role queries.
	Cache mencache.RoleQueryCache

	// Repository provides access to the role query data layer.
	Repository repository.RoleQueryRepository

	// Logger is used for logging service activity.
	Logger logger.LoggerInterface

	// Mapper maps domain models to API response formats.
	Mapper roleservicemapper.RoleQueryResponseMapper
}

// roleQueryService handles read-only operations related to roles.
type roleQueryService struct {
	// errorhandler handles domain-specific errors for role queries.
	errorhandler errorhandler.RoleQueryErrorHandler

	// mencache provides in-memory caching for role queries.
	mencache mencache.RoleQueryCache

	// roleQuery provides access to role query operations in the repository.
	roleQuery repository.RoleQueryRepository

	// logger is used for structured logging within the service.
	logger logger.LoggerInterface

	// mapper maps internal role data to response DTOs.
	mapper roleservicemapper.RoleQueryResponseMapper

	// observability provides tracing and metrics for role query operations.
	observability observability.TraceLoggerObservability
}

// NewRoleQueryService initializes a new roleQueryService with the provided parameters.
//
// It sets up the prometheus metrics for counting and measuring the duration of role query requests.
//
// Parameters:
// - params: A pointer to roleQueryDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created roleQueryService.
func NewRoleQueryService(params *roleQueryDeps) *roleQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_query_service_request_total",
			Help: "Total number of requests to the RoleQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RoleQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("role-query-service"), params.Logger, requestCounter, requestDuration)

	return &roleQueryService{
		errorhandler:  params.ErrorHandler,
		mencache:      params.Cache,
		roleQuery:     params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindAll retrieves all role records with pagination and optional search.
//
// This method fetches role records from the database, applies pagination,
// and returns the results along with pagination metadata.
//
// Parameters:
//   - req: The request payload containing pagination and search parameters.
//
// Returns:
//   - A slice of RoleResponse domain models containing the result of the query.
//   - A pointer to an int containing the total number of records available.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryService) FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoles(ctx, req); found {
		logSuccess("Successfully retrieved all role records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindAllRoles(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(
			err, method, "FAILED_FIND_ALL_ROLE", span, &status, zap.Error(err))
	}
	so := s.mapper.ToRolesResponse(res)

	s.mencache.SetCachedRoles(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved all role records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindById retrieves a role record by its ID.
//
// Parameters:
//   - id: the role ID to retrieve.
//
// Returns:
//   - A pointer to RoleResponse containing the role record.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryService) FindById(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleById(ctx, id); found {
		logSuccess("Successfully retrieved role from cache", zap.Int("id", id))

		return data, nil
	}

	res, err := s.roleQuery.FindById(ctx, id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(
			err, method, "FAILED_FIND_ROLE_BY_ID", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	so := s.mapper.ToRoleResponse(res)

	s.mencache.SetCachedRoleById(ctx, id, so)

	logSuccess("Successfully retrieved role", zap.Int("id", id))

	return so, nil
}

// FindByUserId retrieves a role record associated with a user ID.
//
// Parameters:
//   - id: the user ID to retrieve.
//
// Returns:
//   - A slice of RoleResponse domain models containing the result of the query.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryService) FindByUserId(ctx context.Context, id int) ([]*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleByUserId(ctx, id); found {
		logSuccess("Successfully fetched role by user ID from cache", zap.Int("id", id))
		return data, nil
	}

	res, err := s.roleQuery.FindByUserId(ctx, id)

	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, method, "FAILED_FIND_ROLE_BY_USER_ID", span, &status, zap.Error(err))
	}

	so := s.mapper.ToRolesResponse(res)

	s.mencache.SetCachedRoleByUserId(ctx, id, so)

	logSuccess("Successfully fetched role by user ID", zap.Int("id", id))

	return so, nil
}

// FindByActiveRole retrieves active role records with pagination and optional search.
//
// This method fetches role records from the database, applies pagination,
// and returns the results along with pagination metadata.
//
// Parameters:
//   - req: The request payload containing pagination and search parameters.
//
// Returns:
//   - A slice of RoleResponseDeleteAt domain models containing the result of the query.
//   - A pointer to an int containing the total number of records available.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryService) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActiveRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleActive(ctx, req); found {
		logSuccess("Successfully retrieved active role records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_ACTIVE_ROLE", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	so := s.mapper.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleActive(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved active role records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByTrashedRole retrieves trashed role records with pagination and optional search.
//
// This method fetches trashed role records from the database, applies pagination,
// and returns the results along with pagination metadata.
//
// Parameters:
//   - req: The request payload containing pagination and search parameters.
//
// Returns:
//   - A slice of RoleResponseDeleteAt domain models containing the result of the query.
//   - A pointer to an int containing the total number of records available.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryService) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashedRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed role from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_TRASHED_ROLE", span, &status, role_errors.ErrRoleNotFoundRes)
	}
	so := s.mapper.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed role", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

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
func (s *roleQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

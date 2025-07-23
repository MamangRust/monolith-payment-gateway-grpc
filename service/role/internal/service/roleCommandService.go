package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-role/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	roleservicemapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/role"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// roleCommandDeps contains the required dependencies to construct a roleCommandService.
type roleCommandDeps struct {
	// Ctx is the base context for the service.
	Ctx context.Context

	// ErrorHandler handles domain-specific errors for role commands.
	ErrorHandler errorhandler.RoleCommandErrorHandler

	// Cache provides in-memory caching for role command operations.
	Cache mencache.RoleCommandCache

	// Repository is the database layer for role command operations.
	Repository repository.RoleCommandRepository

	// Logger is used for logging service operations.
	Logger logger.LoggerInterface

	// Mapper maps domain models to response DTOs.
	Mapper roleservicemapper.RoleCommandResponseMapper
}

// roleCommandService handles write operations (create, update, delete) for roles.
type roleCommandService struct {
	// errorhandler handles domain-specific errors for role commands.
	errorhandler errorhandler.RoleCommandErrorHandler

	// mencache provides cache functionality for role command data.
	mencache mencache.RoleCommandCache

	// roleCommand is the repository that handles database operations for role commands.
	roleCommand repository.RoleCommandRepository

	// logger is used for structured logging.
	logger logger.LoggerInterface

	// mapper maps internal role entities to response DTOs.
	mapper roleservicemapper.RoleCommandResponseMapper

	// observability provides tracing and metrics for role command operations.
	observability observability.TraceLoggerObservability
}

// NewRoleCommandService initializes a new RoleCommandService with the given parameters.
// The function sets up the prometheus metrics for request counters and durations,
// and returns a configured RoleCommandService ready for handling role-related commands.
//
// Parameters:
// - params: A pointer to roleCommandDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to an initialized roleCommandService.
func NewRoleCommandService(params *roleCommandDeps) RoleCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_command_service_request_total",
			Help: "Total number of requests to the RoleCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RoleCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("role-command-service"), params.Logger, requestCounter, requestDuration)

	return &roleCommandService{
		errorhandler:  params.ErrorHandler,
		mencache:      params.Cache,
		roleCommand:   params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// CreateRole creates a new role based on the given request.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing the role details.
//
// Returns:
//   - *response.RoleResponse: The created role.
//   - *response.ErrorResponse: An error response if creation failed.
func (s *roleCommandService) CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "CreateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.CreateRole(ctx, request)

	if err != nil {
		return s.errorhandler.HandleCreateRoleError(err, method, "FAILED_CREATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapper.ToRoleResponse(role)

	logSuccess("Successfully created role", zap.Int("role.id", role.ID))

	return so, nil
}

// UpdateRole updates an existing role based on the given request.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing updated role details.
//
// Returns:
//   - *response.RoleResponse: The updated role.
//   - *response.ErrorResponse: An error response if update failed.
func (s *roleCommandService) UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "UpdateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.UpdateRole(ctx, request)
	if err != nil {
		return s.errorhandler.HandleUpdateRoleError(err, method, "FAILED_UPDATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapper.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(ctx, *request.ID)

	logSuccess("Successfully updated role", zap.Int("role.id", role.ID))

	return so, nil
}

// TrashedRole soft-deletes (moves to trash) a role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to be trashed.
//
// Returns:
//   - *response.RoleResponse: The trashed role.
//   - *response.ErrorResponse: An error response if the operation failed.
func (s *roleCommandService) TrashedRole(ctx context.Context, id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.TrashedRole(ctx, id)

	if err != nil {
		return s.errorhandler.HandleTrashedRoleError(err, method, "FAILED_TRASH_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapper.ToRoleResponseDeleteAt(role)

	s.mencache.DeleteCachedRole(ctx, id)

	logSuccess("Successfully trashed role", zap.Int("role.id", role.ID))

	return so, nil
}

// RestoreRole restores a soft-deleted role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to be restored.
//
// Returns:
//   - *response.RoleResponse: The restored role.
//   - *response.ErrorResponse: An error response if the restoration failed.
func (s *roleCommandService) RestoreRole(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "RestoreRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.RestoreRole(ctx, id)

	if err != nil {
		return s.errorhandler.HandleRestoreRoleError(err, method, "FAILED_RESTORE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapper.ToRoleResponse(role)

	logSuccess("Successfully restored role", zap.Int("role.id", role.ID))

	return so, nil
}

// DeleteRolePermanent permanently deletes a role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to be permanently deleted.
//
// Returns:
//   - bool: True if deletion was successful.
//   - *response.ErrorResponse: An error response if the deletion failed.
func (s *roleCommandService) DeleteRolePermanent(ctx context.Context, id int) (bool, *response.ErrorResponse) {
	const method = "DeleteRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.roleCommand.DeleteRolePermanent(ctx, id)
	if err != nil {
		return s.errorhandler.HandleDeleteRolePermanentError(err, method, "FAILED_DELETE_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted role", zap.Int("role.id", id))

	return true, nil
}

// RestoreAllRole restores all trashed roles.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if restoration was successful.
//   - *response.ErrorResponse: An error response if the operation failed.
func (s *roleCommandService) RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.roleCommand.RestoreAllRole(ctx)
	if err != nil {
		return s.errorhandler.HandleRestoreAllRoleError(err, method, "FAILED_RESTORE_ALL_ROLE", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all roles", zap.Bool("success", true))

	return true, nil
}

// DeleteAllRolePermanent permanently deletes all trashed roles.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if deletion was successful.
//   - *response.ErrorResponse: An error response if the operation failed.
func (s *roleCommandService) DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.roleCommand.DeleteAllRolePermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllRolePermanentError(err, method, "FAILED_DELETE_ALL_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted all roles", zap.Bool("success", true))

	return true, nil
}

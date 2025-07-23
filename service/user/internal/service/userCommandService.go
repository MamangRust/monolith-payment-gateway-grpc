package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/service"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/user"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// UserCommandDeps defines the required dependencies to construct a userCommandService.
type userCommandDeps struct {
	// Ctx is the context used across the service.
	Ctx context.Context

	// ErrorHandler handles domain-level errors for user commands.
	ErrorHandler errorhandler.UserCommandError

	// Cache provides caching capabilities for user command operations.
	Cache mencache.UserCommandCache

	// UserQueryRepository provides read access to user data.
	UserQueryRepository repository.UserQueryRepository

	// UserCommandRepository provides write access to user data.
	UserCommandRepository repository.UserCommandRepository

	// RoleRepository provides access to role data, used during user-role assignment.
	RoleRepository repository.RoleRepository

	// Logger is used to log service operations and errors.
	Logger logger.LoggerInterface

	// Mapper maps internal user entities to response models.
	Mapper responseservice.UserCommandResponseMapper

	// Hashing provides password hashing utilities.
	Hashing hash.HashPassword
}

// userCommandService provides operations for creating, updating, and deleting users.
//
// It handles business logic related to user modification, including role assignment,
// password hashing, and interaction with repository, caching, logging, and metrics.
type userCommandService struct {
	// errorhandler handles user command-related errors.
	errorhandler errorhandler.UserCommandError

	// mencache provides caching for user data.
	mencache mencache.UserCommandCache

	// userQueryRepository is used to query user data for validations.
	userQueryRepository repository.UserQueryRepository

	// userCommandRepository is used to create, update, or delete user records.
	userCommandRepository repository.UserCommandRepository

	// roleRepository is used to retrieve and validate user roles.
	roleRepository repository.RoleRepository

	// logger is used for logging service activities and errors.
	logger logger.LoggerInterface

	// mapper maps user entities to API or response DTOs.
	mapper responseservice.UserCommandResponseMapper

	// hashing provides utilities for hashing and verifying passwords.
	hashing hash.HashPassword

	observability observability.TraceLoggerObservability
}

// NewUserCommandService initializes a new instance of userCommandService with the provided parameters.
// It sets up Prometheus metrics for tracking request counts and durations, and returns a configured
// userCommandService ready for handling user-related command operations.
//
// Parameters:
// - params: A pointer to userCommandDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to an initialized userCommandService.
func NewUserCommandService(
	params *userCommandDeps,
) UserCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_command_service_requests_total",
			Help: "Total number of requests to the UserCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the UserCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observabiliy := observability.NewTraceLoggerObservability(
		otel.Tracer("user-command-service"), params.Logger, requestCounter, requestDuration)

	return &userCommandService{
		mencache:              params.Cache,
		errorhandler:          params.ErrorHandler,
		userQueryRepository:   params.UserQueryRepository,
		userCommandRepository: params.UserCommandRepository,
		roleRepository:        params.RoleRepository,
		logger:                params.Logger,
		mapper:                params.Mapper,
		hashing:               params.Hashing,
		observability:         observabiliy,
	}
}

// CreateUser creates a new user with the provided request data.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The user creation request payload.
//
// Returns:
//   - *response.UserResponse: The created user response.
//   - *response.ErrorResponse: Error response if creation fails.
func (s *userCommandService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "CreateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER_BY_EMAIL", span, &status, user_errors.ErrUserEmailAlready, zap.Error(err))

	} else if existingUser != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER_BY_EMAIL", span, &status, user_errors.ErrUserEmailAlready, zap.Error(err))
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		return errorhandler.HandleErrorPasswordOperation[*response.UserResponse](s.logger, err, method, "FAILED_HASH_PASSWORD", span, &status, user_errors.ErrUserPassword, zap.Error(err))
	}

	request.Password = hash

	const defaultRoleName = "Admin Access 1"

	role, err := s.roleRepository.FindByName(ctx, defaultRoleName)

	if err != nil || role == nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ROLE", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	res, err := s.userCommandRepository.CreateUser(ctx, request)

	if err != nil {
		return s.errorhandler.HandleCreateUserError(err, method, "FAILED_CREATE_USER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToUserResponse(res)

	logSuccess("Successfully created user", zap.Int("user.id", res.ID))

	return so, nil
}

// UpdateUser updates an existing user with the provided request data.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The user update request payload.
//
// Returns:
//   - *response.UserResponse: The updated user response.
//   - *response.ErrorResponse: Error response if update fails.
func (s *userCommandService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "UpdateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindById(ctx, *request.UserID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(ctx, request.Email)

		if duplicateUser != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_EMAIL_ALREADY", span, &status, user_errors.ErrUserEmailAlready, zap.Error(err))
		}

		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			return errorhandler.HandleErrorPasswordOperation[*response.UserResponse](s.logger, err, method, "FAILED_HASH_PASSWORD", span, &status, user_errors.ErrUserPassword, zap.Error(err))
		}
		existingUser.Password = hash
	}

	const defaultRoleName = "Admin Access 1"

	role, err := s.roleRepository.FindByName(ctx, defaultRoleName)

	if err != nil || role == nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ROLE", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	res, err := s.userCommandRepository.UpdateUser(ctx, request)

	if err != nil {
		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_UPDATE_USER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToUserResponse(res)

	s.mencache.DeleteUserCache(ctx, so.ID)

	logSuccess("Successfully updated user", zap.Int("user.id", res.ID))

	return so, nil
}

// TrashedUser soft-deletes a user by marking the user as trashed.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user to be trashed.
//
// Returns:
//   - *response.UserResponseDeleteAt: Response including soft-delete timestamp.
//   - *response.ErrorResponse: Error response if trash operation fails.
func (s *userCommandService) TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.TrashedUser(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleTrashedUserError(err, method, "FAILED_TO_TRASH_USER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToUserResponseDeleteAt(res)

	s.mencache.DeleteUserCache(ctx, so.ID)

	logSuccess("Successfully trashed user", zap.Int("user.id", user_id))

	return so, nil
}

// RestoreUser restores a previously soft-deleted user.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user to be restored.
//
// Returns:
//   - *response.UserResponseDeleteAt: Response with restoration details.
//   - *response.ErrorResponse: Error response if restore fails.
func (s *userCommandService) RestoreUser(ctx context.Context, user_id int) (*response.UserResponse, *response.ErrorResponse) {
	const method = "RestoreUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.RestoreUser(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleRestoreUserError(err, "RestoreUser", "FAILED_TO_RESTORE_USER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToUserResponse(res)

	logSuccess("Successfully restored user", zap.Int("user.id", user_id))

	return so, nil
}

// DeleteUserPermanent permanently deletes a user by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user to delete permanently.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - *response.ErrorResponse: Error response if deletion fails.
func (s *userCommandService) DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userCommandRepository.DeleteUserPermanent(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleDeleteUserError(err, "DeleteUserPermanent", "FAILED_TO_DELETE_USER", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted user", zap.Int("user.id", user_id))

	return true, nil
}

// RestoreAllUser restores all soft-deleted users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether all users were successfully restored.
//   - *response.ErrorResponse: Error response if restoration fails.
func (s *userCommandService) RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userCommandRepository.RestoreAllUser(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllUserError(err, method, "FAILED_RESTORE_ALL_USER", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all users", zap.Bool("success", true))

	return true, nil
}

// DeleteAllUserPermanent permanently deletes all trashed users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether all users were successfully deleted permanently.
//   - *response.ErrorResponse: Error response if deletion fails.
func (s *userCommandService) DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userCommandRepository.DeleteAllUserPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllUserError(err, method, "FAILED_DELETE_ALL_USER", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted all users", zap.Bool("success", true))

	return true, nil
}

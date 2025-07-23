package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"

	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/service"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// registerError represents an error handler for user registration operations
type registerError struct {
	logger logger.LoggerInterface
}

// Package errorhandler provides standardized error handling for user registration operations

// NewRegisterError creates and initializes a new registerError handler instance.
//
// Parameters:
//   - logger: The logger instance that will be used for error logging and tracing (logger.LoggerInterface)
//
// Returns:
//   - *registerError: A new instance of the registration error handler ready for use
func NewRegisterError(logger logger.LoggerInterface) *registerError {
	return &registerError{
		logger: logger,
	}
}

// HandleAssignRoleError processes errors that occur during role assignment to a new user.
// This typically happens after successful user creation but before completing registration.
//
// Parameters:
//   - err: The error that occurred during role assignment (error)
//   - method: The name of the calling method (e.g., "CompleteRegistration") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "COMPLETE_REG") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (e.g., "complete_reg_error_assign_role") (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.UserResponse: Nil user response since operation failed
//   - *response.ErrorResponse: Standardized error response with user_not_found error details
func (e *registerError) HandleAssignRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

// HandleFindEmailError processes errors that occur when checking for email existence during registration.
// Used to prevent duplicate email registrations.
//
// Parameters:
//   - err: The error that occurred during email lookup (error)
//   - method: The name of the calling method (e.g., "CheckEmailAvailability") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "CHECK_EMAIL") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.UserResponse: Nil user response since operation failed
//   - *response.ErrorResponse: Standardized error response with user_not_found error details
func (e *registerError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

// HandleFindRoleError processes errors that occur when looking up role information during registration.
// Used when assigning default or requested roles to new users.
//
// Parameters:
//   - err: The error that occurred during role lookup (error)
//   - method: The name of the calling method (e.g., "AssignDefaultRole") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "ASSIGN_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.UserResponse: Nil user response since operation failed
//   - *response.ErrorResponse: Standardized error response with role_not_found error details
func (e *registerError) HandleFindRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		role_errors.ErrRoleNotFoundRes,
		fields...,
	)
}

// HandleCreateUserError processes errors that occur during core user creation in the registration flow.
// This handles failures in the primary user record creation operation.
//
// Parameters:
//   - err: The error that occurred during user creation (error)
//   - method: The name of the calling method (e.g., "CreateUserRecord") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "CREATE_USER") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.UserResponse: Nil user response since operation failed
//   - *response.ErrorResponse: Standardized error response with user_not_found error details
func (e *registerError) HandleCreateUserError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

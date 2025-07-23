package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// roleCommandError is a struct that implements the RoleCommandError interface
type roleCommandError struct {
	logger logger.LoggerInterface
}

// NewRoleCommandError initializes a new roleCommandError with the provided logger.
// It returns an instance of the roleCommandError struct.
func NewRoleCommandError(logger logger.LoggerInterface) RoleCommandErrorHandler {
	return &roleCommandError{
		logger: logger,
	}
}

// HandleCreateRoleError processes errors that occur during role creation.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during role creation (error)
//   - method: The name of the calling method (e.g., "CreateRole") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "CREATE_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.RoleResponse: Nil role response since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleCreateRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedCreateRole,
		fields...,
	)
}

// HandleUpdateRoleError processes errors that occur during role updates.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during role update (error)
//   - method: The name of the calling method (e.g., "UpdateRole") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "UPDATE_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.RoleResponse: Nil role response since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleUpdateRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedUpdateRole,
		fields...,
	)
}

// HandleTrashedRoleError processes errors that occur during role trashing.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during role trashing (error)
//   - method: The name of the calling method (e.g., "TrashRole") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "TRASH_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.RoleResponse: Nil role response since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleTrashedRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedTrashedRole,
		fields...,
	)
}

// HandleRestoreRoleError processes errors that occur during role restoration.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during role restoration (error)
//   - method: The name of the calling method (e.g., "RestoreRole") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "RESTORE_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.RoleResponse: Nil role response since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleRestoreRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedRestoreRole,
		fields...,
	)
}

// HandleDeleteRolePermanentError processes errors that occur during role permanent deletion.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during role permanent deletion (error)
//   - method: The name of the calling method (e.g., "DeleteRolePermanent") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "DELETE_PERMANENT_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: false since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleDeleteRolePermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedDeletePermanent,
		fields...,
	)
}

// HandleDeleteAllRolePermanentError processes errors that occur during all role permanent deletion.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during all role permanent deletion (error)
//   - method: The name of the calling method (e.g., "DeleteAllRolePermanent") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "DELETE_PERMANENT_ALL_ROLES") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: false since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleDeleteAllRolePermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedDeleteAll,
		fields...,
	)
}

// HandleRestoreAllRoleError processes errors that occur during all role restoration.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during all role restoration (error)
//   - method: The name of the calling method (e.g., "RestoreAllRole") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "RESTORE_ALL_ROLE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: false since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (e *roleCommandError) HandleRestoreAllRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedRestoreAll,
		fields...,
	)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RoleCommandErrorHandler interface {
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
	HandleCreateRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
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
	HandleUpdateRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
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
	HandleTrashedRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponseDeleteAt, *response.ErrorResponse)
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
	HandleRestoreRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
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
	HandleDeleteRolePermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
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
	HandleDeleteAllRolePermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
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
	HandleRestoreAllRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type RoleQueryErrorHandler interface {
	// HandleRepositoryPaginationError processes pagination errors from the repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the pagination operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A slice of RoleResponse pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	// HandleRepositoryPaginationDeletedError processes pagination errors from the repository
	// when retrieving deleted role documents. It logs the error, updates the trace span,
	// and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the pagination operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A slice of RoleResponseDeleteAt pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationDeletedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	// HandleRepositoryListError processes list errors from the repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the list operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A slice of RoleResponse pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the list failure.
	HandleRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.RoleResponse, *response.ErrorResponse)
	// HandleRepositorySingleError processes single-result errors from the repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the single-result operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - defaultErr: A pointer to an ErrorResponse that will be updated with the error details.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A RoleResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		defaultErr *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
}

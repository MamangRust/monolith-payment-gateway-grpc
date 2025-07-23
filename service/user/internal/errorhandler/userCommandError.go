package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// userCommandError handles error logging for user command operations.
type userCommandError struct {
	logger logger.LoggerInterface
}

// NewUserCommandError initializes a new userCommandError with the provided logger.
// It returns an instance of the userCommandError struct.
func NewUserCommandError(logger logger.LoggerInterface) UserCommandError {
	return &userCommandError{logger: logger}
}

// HandleRepositorySingleError processes single-result errors from the user repository.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the single-result operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A UserResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the single-result failure.
func (u *userCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

// HandleCreateUserError processes errors that occur during user creation.
// This includes errors related to the user model, repository, and database operations.
//
// Parameters:
//   - err: The error that occurred during user creation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A UserResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the user creation failure.
func (u *userCommandError) HandleCreateUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleUpdateUserError processes errors that occur during user updates.
// This includes errors related to the user model, repository, and database operations.
//
// Parameters:
//   - err: The error that occurred during user update.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A UserResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the user update failure.
func (u *userCommandError) HandleUpdateUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleTrashedUserError processes errors that occur during user soft deletion.
// This includes errors related to the user model, repository, and database operations.
//
// Parameters:
//   - err: The error that occurred during user soft deletion.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A UserResponseDeleteAt pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the user soft deletion failure.
func (u *userCommandError) HandleTrashedUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponseDeleteAt](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleRestoreUserError processes errors that occur during user data restoration.
// This includes errors related to the user model, repository, and database operations.
//
// Parameters:
//   - err: The error that occurred during user data restoration.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The OpenTelemetry span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A UserResponseDeleteAt pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the user data restoration failure.
func (u *userCommandError) HandleRestoreUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleDeleteUserError processes errors that occur during user deletion.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during user deletion.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A boolean indicating success (true) or failure (false) of the operation.
//   - A standardized ErrorResponse describing the user deletion failure.
func (u *userCommandError) HandleDeleteUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleRestoreAllUserError processes errors that occur during all user restoration.
// It logs the error details and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during all user restoration (error)
//   - method: The name of the calling method (e.g., "RestoreAllUser") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "RESTORE_ALL_USER") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: false since operation failed
//   - *response.ErrorResponse: Standardized error response with error details
func (u *userCommandError) HandleRestoreAllUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleDeleteAllUserError processes errors that occur during all user deletion.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during all user deletion.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A boolean indicating success (true) or failure (false) of the operation.
//   - A standardized ErrorResponse describing the all user deletion failure.
func (u *userCommandError) HandleDeleteAllUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

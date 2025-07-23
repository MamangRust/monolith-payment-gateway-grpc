package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	refreshtoken_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/refresh_token_errors/service"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// identityError represents an error handler for identity-related operations.
type identityError struct {
	logger logger.LoggerInterface
}

// NewIdentityError creates and initializes a new identityError handler instance.
// This handler is specifically designed to process errors related to identity verification and token management.
//
// Parameters:
//   - logger: The logger instance that will be used for error logging and tracing (logger.LoggerInterface)
//
// Returns:
//   - *identityError: A new instance of the identity error handler ready for use
func NewIdentityError(logger logger.LoggerInterface) IdentityErrorHandler {
	return &identityError{
		logger: logger,
	}
}

// HandleInvalidTokenError processes errors related to invalid tokens during identity operations.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A TokenResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleInvalidTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedInvalidToken,
		fields...,
	)
}

// HandleExpiredRefreshTokenError processes errors related to expired refresh tokens during identity operations.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A TokenResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleExpiredRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedExpire,
		fields...,
	)
}

// HandleInvalidFormatUserIDError handles errors due to invalid user ID formats.
// It logs the error, updates the trace span, and returns a zero value of the specified type T
// along with a standardized ErrorResponse.
//
// Parameters:
//   - logger: The logger instance for error logging (logger.LoggerInterface)
//   - err: The error encountered due to invalid user ID format (error)
//   - method: The name of the method where the error occurred (string)
//   - tracePrefix: Prefix used for generating trace IDs (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string to be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - A zero value of the specified type T
//   - A pointer to the error response (*response.ErrorResponse)
func HandleInvalidFormatUserIDError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorInvalidID[T](
		logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

// HandleDeleteRefreshTokenError processes errors during refresh token deletion
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A TokenResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleDeleteRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedDeleteRefreshToken,
		fields...,
	)
}

// HandleUpdateRefreshTokenError processes errors during refresh token updates.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A TokenResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleUpdateRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedUpdateRefreshToken,
		fields...,
	)
}

// HandleValidateTokenError processes token validation errors
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A UserResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleValidateTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedInvalidToken,
		fields...,
	)
}

// HandleFindByIdError processes errors during user lookup by ID
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A UserResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleFindByIdError(
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

// HandleGetMeError processes errors during user data retrieval
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A UserResponse with error details and a standardized ErrorResponse.
func (e *identityError) HandleGetMeError(
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

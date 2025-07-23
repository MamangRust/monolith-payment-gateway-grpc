package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// passwordResetError represents an error handler for password reset operations.
type passwordResetError struct {
	logger logger.LoggerInterface
}

// NewPasswordResetError initializes and returns a new instance of passwordResetError.
//
// Parameters:
//   - logger: The logger instance that will be used for error logging and tracing (logger.LoggerInterface)
//
// Returns:
//   - *passwordResetError: A new instance of the password reset error handler
func NewPasswordResetError(logger logger.LoggerInterface) *passwordResetError {
	return &passwordResetError{
		logger: logger,
	}
}

// HandleFindEmailError processes errors that occur when looking up a user's email during password reset.
//
// Parameters:
//   - err: The error that occurred during email lookup (error)
//   - method: The name of the calling method (e.g., "InitiatePasswordReset") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "INIT_PW_RESET") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

// HandleCreateResetTokenError processes errors that occur during reset token generation.
//
// Parameters:
//   - err: The error that occurred during token creation (error)
//   - method: The name of the calling method (e.g., "GenerateResetToken") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "GEN_RESET_TOKEN") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleCreateResetTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrInternalServerError,
		fields...,
	)
}

// HandleFindTokenError processes errors that occur when looking up a reset token.
//
// Parameters:
//   - err: The error that occurred during token lookup (error)
//   - method: The name of the calling method (e.g., "ValidateResetToken") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "VALIDATE_TOKEN") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleFindTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

// HandleUpdatePasswordError processes errors that occur during password updates.
//
// Parameters:
//   - err: The error that occurred during password update (error)
//   - method: The name of the calling method (e.g., "CompletePasswordReset") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "COMPLETE_PW_RESET") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleUpdatePasswordError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

// HandleDeleteTokenError processes errors that occur when deleting used reset tokens.
//
// Parameters:
//   - err: The error that occurred during token deletion (error)
//   - method: The name of the calling method (e.g., "CleanupResetToken") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "CLEANUP_TOKEN") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleDeleteTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

// HandleUpdateVerifiedError processes errors that occur when updating verification status.
//
// Parameters:
//   - err: The error that occurred during status update (error)
//   - method: The name of the calling method (e.g., "MarkAsVerified") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "MARK_VERIFIED") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleUpdateVerifiedError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

// HandleVerifyCodeError processes errors that occur during verification code validation.
//
// Parameters:
//   - err: The error that occurred during code validation (error)
//   - method: The name of the calling method (e.g., "CheckVerificationCode") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "CHECK_VERIFY_CODE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - bool: Always returns false indicating failure
//   - *response.ErrorResponse: Standardized error response with details
func (e *passwordResetError) HandleVerifyCodeError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

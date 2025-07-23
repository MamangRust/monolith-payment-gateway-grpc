package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	refreshtoken_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/refresh_token_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// tokenError handles errors related to JWT token operations
type tokenError struct {
	logger logger.LoggerInterface
}

// NewTokenError creates and initializes a new tokenError handler instance.
// This handler is specifically designed to process errors related to JWT token operations.
//
// Parameters:
//   - logger: The logger instance that will be used for error logging and tracing (logger.LoggerInterface)
//
// Returns:
//   - *tokenError: A new instance of the token error handler ready for use
func NewTokenError(logger logger.LoggerInterface) *tokenError {
	return &tokenError{
		logger: logger,
	}
}

// HandleCreateAccessTokenError processes errors that occur during access token generation.
// This includes JWT signing errors, claims validation failures, and token encoding issues.
//
// Parameters:
//   - err: The error that occurred during access token creation (error)
//   - method: The name of the calling method (e.g., "GenerateAccessToken") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "GEN_ACCESS_TOKEN") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (e.g., "gen_access_token_error_create_token") (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.TokenResponse: Nil token response since operation failed
//   - *response.ErrorResponse: Standardized error response with failed_create_access error details,
//     typically containing error code 500 (Internal Server Error)
func (e *tokenError) HandleCreateAccessTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedCreateAccess,
		fields...,
	)
}

// HandleCreateRefreshTokenError processes errors that occur during refresh token generation.
// This handles failures in token persistence, cryptographic operations, and storage errors.
//
// Parameters:
//   - err: The error that occurred during refresh token creation (error)
//   - method: The name of the calling method (e.g., "GenerateRefreshToken") (string)
//   - tracePrefix: A prefix used for generating trace IDs (e.g., "GEN_REFRESH_TOKEN") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with error status (*string)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - *response.TokenResponse: Nil token response since operation failed
//   - *response.ErrorResponse: Standardized error response with failed_create_refresh error details,
//     typically containing error code 500 (Internal Server Error)
func (e *tokenError) HandleCreateRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedCreateRefresh,
		fields...,
	)
}

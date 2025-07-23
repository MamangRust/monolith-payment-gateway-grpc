package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// passwordError represents an error handler for password operations.
type passwordError struct {
	logger logger.LoggerInterface
}

// NewPasswordError initializes a new passwordError with a logger.
func NewPasswordError(logger logger.LoggerInterface) *passwordError {
	return &passwordError{
		logger: logger,
	}
}

// HandlePasswordNotMatchError processes password mismatch errors
// Returns boolean status and standardized ErrorResponse
// Parameters:
//   - err: The error that occurred during password comparison.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
func (e *passwordError) HandlePasswordNotMatchError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorPasswordOperation[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		"not match",
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

// HandleHashPasswordError processes password hashing errors
// Returns UserResponse with error details and standardized ErrorResponse
// Parameters:
//   - err: The error that occurred during password comparison.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
func (e *passwordError) HandleHashPasswordError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorPasswordOperation[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		"hash",
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

// HandleComparePasswordError processes password comparison errors
// Returns TokenResponse with error details and standardized ErrorResponse
// Parameters:
//   - err: The error that occurred during password comparison.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
func (e *passwordError) HandleComparePasswordError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorPasswordOperation[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		"compare",
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

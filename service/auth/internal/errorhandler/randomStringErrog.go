package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type randomStringError struct {
	logger logger.LoggerInterface
}

// NewRandomStringError initializes a new randomStringError with a logger.
func NewRandomStringError(logger logger.LoggerInterface) *randomStringError {
	return &randomStringError{
		logger: logger,
	}
}

// HandleRandomStringErrorRegister processes errors that occur during random string generation
// for user registration. It leverages handleErrorGenerateRandomString to log the error and update
// the trace span with relevant error details.
//
// Parameters:
//   - err: The error that occurred during random string generation.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The OpenTelemetry span for distributed tracing.
//   - status: A pointer to a string that will be updated with the error status.
//   - fields: Additional context fields for structured logging.
//
// Returns:
//   - A UserResponse with user-related information if available.
//   - A standardized ErrorResponse containing error details.
func (r randomStringError) HandleRandomStringErrorRegister(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorGenerateRandomString[*response.UserResponse](
		r.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrInternalServerError,
		fields...,
	)
}

// HandleRandomStringErrorForgotPassword processes errors that occur during random string generation
// for forgot password operations. It leverages handleErrorGenerateRandomString to log the error and
// update the trace span with relevant error details.
//
// Parameters:
//   - err: The error that occurred during random string generation.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The OpenTelemetry span for distributed tracing.
//   - status: A pointer to a string that will be updated with the error status.
//   - fields: Additional context fields for structured logging.
//
// Returns:
//   - A boolean indicating whether the operation was successful (false) or not (true).
//   - A standardized ErrorResponse containing error details.
func (h *randomStringError) HandleRandomStringErrorForgotPassword(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorGenerateRandomString[bool](
		h.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrInternalServerError,
		fields...,
	)
}

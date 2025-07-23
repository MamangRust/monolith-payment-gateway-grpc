package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// loginError handles errors related to login operations.
type loginError struct {
	logger logger.LoggerInterface
}

// NewLoginError creates and initializes a new loginError handler instance.
// It takes a logger as input and returns a pointer to the loginError struct.
//
// Parameters:
//   - logger: The logger instance that will be used for error logging and tracing (logger.LoggerInterface)
//
// Returns:
//   - *loginError: A new instance of the login error handler ready for use
func NewLoginError(logger logger.LoggerInterface) *loginError {
	return &loginError{
		logger: logger,
	}
}

// HandleFindEmailError processes errors encountered during the email lookup
// for login operations. It logs the error, records it to the trace span,
// and returns a standardized error response.
//
// Parameters:
// - err: The error that occurred during email lookup.
// - method: The name of the method where the error occurred.
// - tracePrefix: A prefix for generating the trace ID.
// - span: The trace span used for recording the error.
// - status: A pointer to a string that will be set with the formatted status.
// - fields: Additional fields to include in the log entry.
//
// Returns:
// - A TokenResponse with error details and a standardized ErrorResponse.
func (e *loginError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
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

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// marshalError represents an error handler for marshaling errors.
type marshalError struct {
	logger logger.LoggerInterface
}

// NewMarshalError initializes a new marshalError.
//
// Parameters:
//   - logger: The logger instance that will be used for error logging and tracing (logger.LoggerInterface)
//
// Returns:
//   - *marshalError: A pointer to the initialized marshalError.
func NewMarshalError(logger logger.LoggerInterface) *marshalError {
	return &marshalError{
		logger: logger,
	}
}

// HandleMarshalRegisterError processes errors that occur during the marshaling
// of registration data. It leverages JSON marshaling error handling to standardize
// the error response format.
//
// Parameters:
//   - err: The error that occurred during marshaling.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - UserResponse: The user response containing error details.
//   - ErrorResponse: The standardized error response.
func (e *marshalError) HandleMarshalRegisterError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorJSONMarshal[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrFailedSendEmail,
		fields...,
	)
}

// HandleMarsalForgotPassword processes errors that occur during the marshaling
// of forgot password data. It leverages JSON marshaling error handling to standardize
// the error response format.
//
// Parameters:
//   - err: The error that occurred during marshaling.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - bool: A boolean indicating the success or failure of the operation.
//   - ErrorResponse: The standardized error response containing error details.
func (e *marshalError) HandleMarsalForgotPassword(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorJSONMarshal[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrFailedSendEmail,
		fields...,
	)
}

// HandleMarshalVerifyCode processes errors that occur during the marshaling
// of verification code data. It leverages JSON marshaling error handling to standardize
// the error response format.
//
// Parameters:
//   - err: The error that occurred during marshaling.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - bool: A boolean indicating the success or failure of the operation.
//   - ErrorResponse: The standardized error response containing error details.
func (e *marshalError) HandleMarshalVerifyCode(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorJSONMarshal[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrFailedSendEmail,
		fields...,
	)
}

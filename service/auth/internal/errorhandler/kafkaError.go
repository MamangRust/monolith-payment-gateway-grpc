package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// kafkaError represents an error handler for Kafka operations.
type kafkaError struct {
	logger logger.LoggerInterface
}

// NewKafkaError initializes a new KafkaErrorHandler with a logger.
// It takes a logger as input and returns a pointer to the KafkaErrorHandler struct.
func NewKafkaError(logger logger.LoggerInterface) KafkaErrorHandler {
	return &kafkaError{
		logger: logger,
	}
}

// HandleSendEmailForgotPassword processes errors that occur during the sending of forgot password emails.
// It utilizes Kafka for message handling and returns a boolean indicating success or failure, along with a standardized ErrorResponse.
// Parameters:
//   - err: The error encountered during the email sending process.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for trace logging.
//   - span: The tracing span for monitoring.
//   - status: A pointer to a string representing the status of the operation.
//   - fields: Additional logging fields for structured logging.
//
// Returns:
//   - A boolean indicating success (false) or failure (true) of the operation.
//   - A pointer to a standardized ErrorResponse containing error details.
func (e *kafkaError) HandleSendEmailForgotPassword(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorKafkaSend[bool](
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

// HandleSendEmailRegister processes errors that occur during the sending of registration emails.
// It utilizes Kafka for message handling and returns a UserResponse with error details and a standardized ErrorResponse.
// Parameters:
//   - err: The error encountered during the email sending process.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for trace logging.
//   - span: The tracing span for monitoring.
//   - status: A pointer to a string representing the status of the operation.
//   - fields: Additional logging fields for structured logging.
//
// Returns:
//   - A UserResponse containing user-related information if available.
//   - A pointer to a standardized ErrorResponse containing error details.
func (e *kafkaError) HandleSendEmailRegister(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorKafkaSend[*response.UserResponse](
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

// HandleSendEmailVerifyCode processes errors that occur during the sending of verification code emails.
// It utilizes Kafka for message handling and returns a boolean indicating success or failure, along with a standardized ErrorResponse.
// Parameters:
//   - err: The error encountered during the email sending process.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for trace logging.
//   - span: The tracing span for monitoring.
//   - status: A pointer to a string representing the status of the operation.
//   - fields: Additional logging fields for structured logging.
//
// Returns:
//   - A boolean indicating success (false) or failure (true) of the operation.
//   - A pointer to a standardized ErrorResponse containing error details.
func (e *kafkaError) HandleSendEmailVerifyCode(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorKafkaSend[bool](
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

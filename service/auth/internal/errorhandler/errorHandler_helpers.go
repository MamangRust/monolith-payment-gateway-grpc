package errorhandler

import (
	"regexp"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// handleErrorRepository specializes sharederrorhandler.HandleErrorTemplate for repository layer errors.
// It automatically sets the error message to "repository error" and follows the
// same standardized error handling pattern.
//
// Parameters:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error from repository operation
//   - method: Name of the calling method
//   - tracePrefix: Prefix for trace ID generation
//   - span: OpenTelemetry span for tracing
//   - status: Pointer to status string to be updated
//   - errorResp: Predefined error response
//   - fields: Additional contextual log fields
//
// Returns:
//   - Zero value of type T
//   - Pointer to response.ErrorResponse
func handleErrorRepository[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](
		logger, err, method, tracePrefix,
		"Repository error", span, status, errorResp, fields...,
	)
}

// handleErrorTokenTemplate specializes sharederrorhandler.HandleErrorTemplate for token-related errors.
// It automatically sets the error message to "token error" and follows the
// standardized error handling pattern for authentication/authorization failures.
//
// Parameters:
//   - logger: LoggerInterface instance
//   - err: The token-related error
//   - method: Name of the calling method
//   - tracePrefix: Trace ID prefix
//   - span: OpenTelemetry span
//   - status: Pointer to status string
//   - defaultErr: Default error response
//   - fields: Additional log fields
//
// Returns:
//   - Zero value of type T
//   - Pointer to response.ErrorResponse
func handleErrorTokenTemplate[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "token", span, status, defaultErr, fields...)
}

// HandleTokenError is a helper function used to process errors related to token operations.
// It logs the error using the provided logger, records the error to the trace span,
// and sets the status to a standardized format. It returns a zero value of the specified type T
// and a pointer to a standardized ErrorResponse.
//
// Parameters:
//
//	logger - LoggerInterface used for logging the error.
//	err - The error that occurred.
//	method - The name of the method where the error occurred.
//	tracePrefix - A prefix for generating the trace ID.
//	span - The trace span used for recording the error.
//	status - A pointer to a string that will be set with the formatted status.
//	defaultErr - The default error response to return.
//	fields - Additional fields to include in the log entry.
//
// Returns:
//
//	A zero value of type T and a pointer to an ErrorResponse.
func HandleTokenError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTokenTemplate[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

// handleErrorJSONMarshal specializes error handling for JSON marshaling failures.
//
// Parameters:
//   - logger: Logger instance
//   - err: Marshaling error
//   - method: Calling method name
//   - tracePrefix: Trace prefix
//   - span: Tracing span
//   - status: Status reference
//   - defaultErr: Default error
//   - fields: Log fields
//
// Returns:
//   - Zero value and error response
func handleErrorJSONMarshal[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "json marshal", span, status, defaultErr, fields...)
}

// handleErrorKafkaSend specializes error handling for Kafka producer failures.
//
// Parameters:
//   - logger: Logger instance
//   - err: Kafka send error
//   - method: Calling method
//   - tracePrefix: Trace prefix
//   - span: Tracing span
//   - status: Status reference
//   - defaultErr: Default error
//   - fields: Log fields
//
// Returns:
//   - Zero value and error response
func handleErrorKafkaSend[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "kafka send", span, status, defaultErr, fields...)
}

// handleErrorGenerateRandomString handles errors that occur during random string generation operations.
// It wraps the error using the standard error template with a predefined "generate random string" error message.
//
// Parameters:
//   - logger: The logger instance for recording error details (logger.LoggerInterface)
//   - err: The error that occurred during random string generation (error)
//   - method: The name of the method where the error occurred (e.g., "CreateVerificationCode") (string)
//   - tracePrefix: The prefix for generating trace IDs (e.g., "CREATE_VERIFICATION_CODE") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (e.g., "create_verification_code_error_generate_random_string") (*string)
//   - defaultErr: The predefined error response to return (*response.ErrorResponse)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - A zero value of the specified type T
//   - A pointer to the error response (*response.ErrorResponse)
func handleErrorGenerateRandomString[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "generate random string", span, status, defaultErr, fields...)
}

// handleErrorInvalidID handles errors related to invalid ID formats or values.
// It wraps the error using the standard error template with a predefined "invalid id" error message.
//
// Parameters:
//   - logger: The logger instance for recording error details (logger.LoggerInterface)
//   - err: The error that occurred due to invalid ID (error)
//   - method: The name of the method where the error occurred (e.g., "GetUserByID") (string)
//   - tracePrefix: The prefix for generating trace IDs (e.g., "GET_USER_BY_ID") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (e.g., "get_user_by_id_error_invalid_id") (*string)
//   - defaultErr: The predefined error response to return (*response.ErrorResponse)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - A zero value of the specified type T
//   - A pointer to the error response (*response.ErrorResponse)
func handleErrorInvalidID[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "invalid id", span, status, defaultErr, fields...)
}

// handleErrorPasswordOperation handles errors related to password operations (hashing, validation, etc.).
// It allows specifying a custom operation name in the error message for more context.
//
// Parameters:
//   - logger: The logger instance for recording error details (logger.LoggerInterface)
//   - err: The error that occurred during password operation (error)
//   - method: The name of the method where the error occurred (e.g., "ChangePassword") (string)
//   - tracePrefix: The prefix for generating trace IDs (e.g., "CHANGE_PASSWORD") (string)
//   - operation: The specific password operation that failed (e.g., "hashing", "validation") (string)
//   - span: The OpenTelemetry span for distributed tracing (trace.Span)
//   - status: Pointer to a string that will be updated with the error status (e.g., "change_password_error_hashing") (*string)
//   - defaultErr: The predefined error response to return (*response.ErrorResponse)
//   - fields: Additional context fields for structured logging (...zap.Field)
//
// Returns:
//   - A zero value of the specified type T
//   - A pointer to the error response (*response.ErrorResponse)
func handleErrorPasswordOperation[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix, operation string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, operation, span, status, defaultErr, fields...)
}

// HandleRepositorySingleError is the public interface for single-result repository errors.
// It provides standardized handling of database operation failures.
//
// Parameters:
//   - logger: Structured logger
//   - err: Database error
//   - method: Calling method
//   - tracePrefix: Trace prefix
//   - span: Tracing span
//   - status: Status reference
//   - defaultErr: Default error
//   - fields: Context fields
//
// Returns:
//   - Zero value and error response
func HandleRepositorySingleError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorRepository[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

// toSnakeCase converts a camelCase string to a snake_case string.
//
// Parameters:
//
//   - s: CamelCase string
//
// Returns:
//
//   - Snake case equivalent of the input string.
func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

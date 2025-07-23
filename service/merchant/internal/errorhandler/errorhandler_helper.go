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
// Args:
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

// handleErrorPagination specializes handleErrorRepository for pagination layer errors.
// It automatically sets the error message to "pagination error" and follows the
// same standardized error handling pattern.
//
// Args:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error from pagination operation
//   - method: Name of the calling method
//   - tracePrefix: Prefix for trace ID generation
//   - span: OpenTelemetry span for tracing
//   - status: Pointer to status string to be updated
//   - errorResp: Predefined error response
//   - fields: Additional contextual log fields
//
// Returns:
//   - Zero value of type T
//   - Nil pointer to int
//   - Pointer to response.ErrorResponse
func handleErrorPagination[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *int, *response.ErrorResponse) {
	result, errResp := handleErrorRepository[T](
		logger, err, method, tracePrefix, span, status, errorResp, fields...,
	)
	return result, nil, errResp
}

// handleErrorMarshal specializes error handling for marshaling failures.
// It automatically sets the error message to "Marshal error" and follows the
// same standardized error handling pattern.
//
// Args:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error from marshaling operation
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
func handleErrorMarshal[T any](
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
		"Marshal error", span, status, errorResp, fields...,
	)
}

// handleErrorSendEmail specializes error handling for email sending failures.
// It sets the error message to "Send email error" and adheres to a standardized
// error handling pattern using the sharederrorhandler.HandleErrorTemplate function.
//
// Args:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error from the email sending operation
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
func handleErrorSendEmail[T any](
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
		"Send email error", span, status, errorResp, fields...,
	)
}

// HandleRepositorySingleError is the public interface for single-result repository errors.
// It provides standardized handling of database operation failures.
//
// Args:
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

// HandleMarshalError specializes error handling for marshaling failures.
// It automatically sets the error message to "Marshal error" and follows the
// same standardized error handling pattern.
//
// Args:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error from marshaling operation
//   - method: Name of the calling method
//   - tracePrefix: Prefix for trace ID generation
//   - span: OpenTelemetry span for tracing
//   - status: Pointer to status string to be updated
//   - defaultErr: Predefined error response
//   - fields: Additional contextual log fields
//
// Returns:
//   - Zero value of type T
//   - Pointer to response.ErrorResponse
func HandleMarshalError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorMarshal[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

// HandleSendEmailError specializes error handling for email sending failures.
// It automatically sets the error message to "Send email error" and follows the
// same standardized error handling pattern.
//
// Args:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error from sending email operation
//   - method: Name of the calling method
//   - tracePrefix: Prefix for trace ID generation
//   - span: OpenTelemetry span for tracing
//   - status: Pointer to status string to be updated
//   - defaultErr: Predefined error response
//   - fields: Additional contextual log fields
//
// Returns:
//   - Zero value of type T
//   - Pointer to response.ErrorResponse
func HandleSendEmailError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorSendEmail[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

// toSnakeCase converts a camelCase string to a snake_case string.
//
// Args:
//   - s: CamelCase string
//
// Returns:
//   - Snake case equivalent of the input string.
func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

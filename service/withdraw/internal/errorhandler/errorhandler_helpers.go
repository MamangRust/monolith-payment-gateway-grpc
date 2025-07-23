package errorhandler

import (
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

// handleErrorPagination specializes handleErrorRepository for pagination layer errors.
// It automatically sets the error message to "pagination error" and follows the
// same standardized error handling pattern.
//
// Parameters:
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

// HandleErrorMarshal specializes error handling for marshaling failures.
// It automatically sets the error message to "Marshal error" and follows the
// same standardized error handling pattern.
//
// Parameters:
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
func HandleErrorMarshal[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "Marshal error", span, status, errorResp, fields...)
}

// HandleErrorKafkaSend handles errors that occur during Kafka message sending.
// It uses a standardized template for error handling, logging, and tracing.
//
// Parameters:
//   - logger: LoggerInterface instance for structured logging
//   - err: The error encountered during Kafka sending
//   - method: Name of the calling method
//   - tracePrefix: Prefix for trace ID generation
//   - span: OpenTelemetry span for distributed tracing
//   - status: Pointer to a string for updating the operation status
//   - defaultErr: Predefined error response template
//   - fields: Additional zap fields for contextual logging
//
// Returns:
//   - Zero value of type T
//   - Pointer to a standardized ErrorResponse containing error details
func HandleErrorKafkaSend[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[T](logger, err, method, tracePrefix, "Kafka send error", span, status, defaultErr, fields...)
}

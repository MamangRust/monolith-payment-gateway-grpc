package errorhandler

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func handleErrorTokenTemplate[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)
	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("Token error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Token error in %s", method))

	*status = fmt.Sprintf("token_error_%s", method)

	var zero T
	return zero, defaultErr
}

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

func handleErrorTemplate[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)
	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("Repository error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Repository error in %s", method))

	*status = fmt.Sprintf("repository_error_%s", method)

	var zero T
	return zero, errorResp
}

func handleErrorJSONMarshal[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)

	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("JSON marshal error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("JSON marshal error in %s", method))

	*status = fmt.Sprintf("marshal_json_failed_%s", method)

	var zero T
	return zero, defaultErr
}

func handleErrorKafkaSend[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)

	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("Kafka send error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Kafka send error in %s", method))

	*status = fmt.Sprintf("kafka_send_failed_%s", method)

	var zero T
	return zero, defaultErr
}

func handleErrorGenerateRandomString[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)

	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("Generate random string error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Generate random string error in %s", method))

	*status = fmt.Sprintf("generate_random_string_failed_%s", method)

	var zero T
	return zero, defaultErr
}

func handleErrorInvalidID[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)

	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("Invalid ID error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Invalid ID error in %s", method))

	*status = fmt.Sprintf("invalid_id_error_%s", method)

	var zero T
	return zero, defaultErr
}

func handleErrorPasswordOperation[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix, operation string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)

	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	msg := fmt.Sprintf("%s password error in %s", operation, method)
	logger.Error(msg, allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)

	*status = fmt.Sprintf("%s_password_error_%s", operation, method)

	var zero T
	return zero, defaultErr
}

func HandleRepositorySingleError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

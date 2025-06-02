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

func handleErrorPaginationTemplate[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *int, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)
	allFields := append(fields, zap.Error(err), zap.String("trace.id", traceID))

	logger.Error(fmt.Sprintf("Repository error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Repository error in %s", method))

	*status = fmt.Sprintf("repository_error_%s", method)

	var zero T
	return zero, nil, errorResp
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

	if status != nil {
		*status = fmt.Sprintf("repository_error_%s", method)
	}

	var zero T
	return zero, errorResp
}

func handleErrorMarshal[T any](
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

	logger.Error(fmt.Sprintf("Failed marshal error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Failed marshal error in %s", method))

	if status != nil {
		*status = fmt.Sprintf("json_marshal_error_%s", method)
	}

	var zero T
	return zero, errorResp
}

func handleErrorSendEmail[T any](
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

	logger.Error(fmt.Sprintf("Failed send email error in %s", method), allFields...)
	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, fmt.Sprintf("Failed send email error in %s", method))

	if status != nil {
		*status = fmt.Sprintf("send_email_error_%s", method)
	}

	var zero T
	return zero, errorResp
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

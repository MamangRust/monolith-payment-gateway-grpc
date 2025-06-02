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

func HandleErrorPasswordOperation[T any](
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

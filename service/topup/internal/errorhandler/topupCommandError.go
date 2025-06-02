package errorhandler

import (
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupCommandError struct {
	logger logger.LoggerInterface
}

func NewTopupCommandError(logger logger.LoggerInterface) *topupCommandError {
	return &topupCommandError{
		logger: logger,
	}
}

func (t *topupCommandError) HandleInvalidParseTimeError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	rawTime string,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {

	traceID := traceunic.GenerateTraceID("INVALID_PARSE_TIME")

	t.logger.Error("Invalid time parse error",
		append(fields,
			zap.String("trace.id", traceID),
			zap.String("raw_time", rawTime),
			zap.String("method", method),
			zap.String("trace_prefix", tracePrefix),
			zap.Error(err),
		)...,
	)

	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, "Invalid parse time")

	if status != nil {
		*status = "invalid_parse_time"
	}

	return nil, &response.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse the given time value",
		Status:  "invalid_parse_time",
	}
}

func (e *topupCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TopupResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *topupCommandError) HandleCreateTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedCreateTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleUpdateTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedUpdateTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleTrashedTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TopupResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedTrashTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleRestoreTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TopupResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedRestoreTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleDeleteTopupPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedDeleteTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleRestoreAllTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedRestoreAllTopups,
		fields...,
	)
}

func (e *topupCommandError) HandleDeleteAllTopupPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedDeleteAllTopups,
		fields...,
	)
}

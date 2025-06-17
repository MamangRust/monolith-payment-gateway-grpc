package errorhandler

import (
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
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

func (e *topupCommandError) HandleInvalidParseTimeError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	rawTime string,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	errResp := &response.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse the given time value",
		Status:  "invalid_parse_time",
	}

	return handleErrorTemplate[*response.TopupResponse](e.logger, err, method, tracePrefix, "Invalid parse time", span, status, errResp, fields...)

}

func (e *topupCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *topupCommandError) HandleCreateTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedCreateTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleUpdateTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedUpdateTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleTrashedTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedTrashTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleRestoreTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedRestoreTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleDeleteTopupPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedDeleteTopup,
		fields...,
	)
}

func (e *topupCommandError) HandleRestoreAllTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedRestoreAllTopups,
		fields...,
	)
}

func (e *topupCommandError) HandleDeleteAllTopupPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedDeleteAllTopups,
		fields...,
	)
}

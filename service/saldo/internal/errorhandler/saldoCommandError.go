package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoCommandError struct {
	logger logger.LoggerInterface
}

func NewSaldoCommandError(logger logger.LoggerInterface) *saldoCommandError {
	return &saldoCommandError{
		logger: logger,
	}
}

func (e *saldoCommandError) HandleFindCardByNumberError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		card_errors.ErrCardNotFoundRes,
		fields...,
	)
}

func (e *saldoCommandError) HandleCreateSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedCreateSaldo,
		fields...,
	)
}

func (e *saldoCommandError) HandleUpdateSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedUpdateSaldo,
		fields...,
	)
}

func (e *saldoCommandError) HandleTrashSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedTrashSaldo,
		fields...,
	)
}

func (e *saldoCommandError) HandleRestoreSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedRestoreSaldo,
		fields...,
	)
}

func (e *saldoCommandError) HandleDeleteSaldoPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedDeleteSaldoPermanent,
		fields...,
	)
}

func (e *saldoCommandError) HandleRestoreAllSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedRestoreAllSaldo,
		fields...,
	)
}

func (e *saldoCommandError) HandleDeleteAllSaldoPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedDeleteAllSaldoPermanent,
		fields...,
	)
}

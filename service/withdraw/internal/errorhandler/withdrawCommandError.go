package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawCommandError struct {
	logger logger.LoggerInterface
}

func NewWithdrawCommandError(logger logger.LoggerInterface) *withdrawCommandError {
	return &withdrawCommandError{
		logger: logger,
	}
}

func (t *withdrawCommandError) HandleInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	cardNumber string,
	fields ...zap.Field,
) (*response.WithdrawResponse, *response.ErrorResponse) {

	traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE")

	t.logger.Error("Insufficient balance",
		append(fields,
			zap.String("trace.id", traceID),
			zap.String("card_number", cardNumber),
			zap.String("method", method),
			zap.String("trace_prefix", tracePrefix),
			zap.Error(err),
		)...,
	)

	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, "Insufficient balance")

	if status != nil {
		*status = "insufficient_balance"
	}

	return nil, saldo_errors.ErrFailedInsuffientBalance
}

func (w *withdrawCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (w *withdrawCommandError) HandleCreateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedCreateWithdraw, fields...)
}

func (w *withdrawCommandError) HandleUpdateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedUpdateWithdraw, fields...)
}

func (w *withdrawCommandError) HandleTrashedWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedTrashedWithdraw, fields...)
}

func (w *withdrawCommandError) HandleRestoreWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedRestoreWithdraw, fields...)
}

func (w *withdrawCommandError) HandleDeleteWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedDeleteWithdrawPermanent, fields...)
}

func (w *withdrawCommandError) HandleRestoreAllWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedRestoreAllWithdraw, fields...)
}

func (w *withdrawCommandError) HandleDeleteAllWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedDeleteAllWithdrawPermanent, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
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
	return handleErrorTemplate[*response.WithdrawResponse](t.logger, err, method, tracePrefix, "Insufficient Balance", span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

func (w *withdrawCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (w *withdrawCommandError) HandleCreateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedCreateWithdraw, fields...)
}

func (w *withdrawCommandError) HandleUpdateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedUpdateWithdraw, fields...)
}

func (w *withdrawCommandError) HandleTrashedWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedTrashedWithdraw, fields...)
}

func (w *withdrawCommandError) HandleRestoreWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedRestoreWithdraw, fields...)
}

func (w *withdrawCommandError) HandleDeleteWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedDeleteWithdrawPermanent, fields...)
}

func (w *withdrawCommandError) HandleRestoreAllWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedRestoreAllWithdraw, fields...)
}

func (w *withdrawCommandError) HandleDeleteAllWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedDeleteAllWithdrawPermanent, fields...)
}

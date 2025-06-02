package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardStatisticError struct {
	logger logger.LoggerInterface
}

func NewCardStatisticError(logger logger.LoggerInterface) *cardStatisticError {
	return &cardStatisticError{
		logger: logger,
	}
}

func (c *cardStatisticError) HandleMonthlyBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseMonthBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyBalance, fields...)
}

func (c *cardStatisticError) HandleYearlyBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseYearlyBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyBalance, fields...)
}

func (c *cardStatisticError) HandleMonthlyTopupAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTopupAmount, fields...)
}

func (c *cardStatisticError) HandleYearlyTopupAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTopupAmount, fields...)
}

func (c *cardStatisticError) HandleMonthlyWithdrawAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyWithdrawAmount, fields...)
}

func (c *cardStatisticError) HandleYearlyWithdrawAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyWithdrawAmount, fields...)
}

func (c *cardStatisticError) HandleMonthlyTransactionAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransactionAmount, fields...)
}

func (c *cardStatisticError) HandleYearlyTransactionAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransactionAmount, fields...)
}

func (c *cardStatisticError) HandleMonthlyTransferAmountSenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountSender, fields...)
}

func (c *cardStatisticError) HandleYearlyTransferAmountSenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountSender, fields...)
}

func (c *cardStatisticError) HandleMonthlyTransferAmountReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountReceiver, fields...)
}

func (c *cardStatisticError) HandleYearlyTransferAmountReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountReceiver, fields...)
}

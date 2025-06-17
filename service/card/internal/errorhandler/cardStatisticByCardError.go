package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardStatisticByNumberError struct {
	logger logger.LoggerInterface
}

func NewCardStatisticByNumberError(logger logger.LoggerInterface) cardStatisticByNumberError {
	return cardStatisticByNumberError{
		logger: logger,
	}
}

func (c *cardStatisticByNumberError) HandleMonthlyBalanceByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyBalanceByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleYearlyBalanceByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearlyBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyBalanceByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleMonthlyTopupAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTopupAmountByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleYearlyTopupAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTopupAmountByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleMonthlyWithdrawAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyWithdrawAmountByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleYearlyWithdrawAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyWithdrawAmountByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleMonthlyTransactionAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransactionAmountByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleYearlyTransactionAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransactionAmountByCard, fields...)
}

func (c *cardStatisticByNumberError) HandleMonthlyTransferAmountBySenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountBySender, fields...)
}

func (c *cardStatisticByNumberError) HandleYearlyTransferAmountBySenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountBySender, fields...)
}

func (c *cardStatisticByNumberError) HandleMonthlyTransferAmountByReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountByReceiver, fields...)
}

func (c *cardStatisticByNumberError) HandleYearlyTransferAmountByReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountByReceiver, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatisticByCardError struct {
	logger logger.LoggerInterface
}

func NewTransactionStatisticByCardError(logger logger.LoggerInterface) *transactionStatisticByCardError {
	return &transactionStatisticByCardError{logger: logger}
}

func (e *transactionStatisticByCardError) HandleMonthTransactionStatusSuccessByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseMonthStatusSuccess](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionSuccessByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleYearlyTransactionStatusSuccessByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseYearStatusSuccess](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionSuccessByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleMonthTransactionStatusFailedByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseMonthStatusFailed](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionFailedByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleYearlyTransactionStatusFailedByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseYearStatusFailed](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionFailedByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleMonthlyPaymentMethodsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionMonthMethodResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyPaymentMethodsByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleYearlyPaymentMethodsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionYearMethodResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyPaymentMethodsByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleMonthlyAmountsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionMonthAmountResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmountsByCard, fields...)
}

func (e *transactionStatisticByCardError) HandleYearlyAmountsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionYearlyAmountResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmountsByCard, fields...)
}

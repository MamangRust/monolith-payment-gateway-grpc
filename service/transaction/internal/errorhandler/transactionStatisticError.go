package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatisticError struct {
	logger logger.LoggerInterface
}

func NewTransactionStatisticError(logger logger.LoggerInterface) *transactionStatisticError {
	return &transactionStatisticError{logger: logger}
}

func (t *transactionStatisticError) HandleMonthTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseMonthStatusSuccess](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionSuccess, fields...)
}

func (t *transactionStatisticError) HandleYearlyTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseYearStatusSuccess](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionSuccess, fields...)
}

func (t *transactionStatisticError) HandleMonthTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseMonthStatusFailed](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionFailed, fields...)
}

func (t *transactionStatisticError) HandleYearlyTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionResponseYearStatusFailed](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionFailed, fields...)
}

func (t *transactionStatisticError) HandleMonthlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionMonthMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyPaymentMethods, fields...)
}

func (t *transactionStatisticError) HandleYearlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionYearMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyPaymentMethods, fields...)
}

func (t *transactionStatisticError) HandleMonthlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionMonthAmountResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmounts, fields...)
}

func (t *transactionStatisticError) HandleYearlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransactionYearlyAmountResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmounts, fields...)
}

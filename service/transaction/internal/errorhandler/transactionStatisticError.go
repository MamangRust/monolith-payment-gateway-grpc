package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transactionStatisticError is a struct that implements the TransactionStatisticError interface
type transactionStatisticError struct {
	logger logger.LoggerInterface
}

// NewTransactionStatisticError returns a new instance of transactionStatisticError with the provided logger.
// It is used to handle errors related to transaction statistics, ensuring they are logged appropriately.
func NewTransactionStatisticError(logger logger.LoggerInterface) TransactionStatisticErrorHandler {
	return &transactionStatisticError{logger: logger}
}

// HandleMonthTransactionStatusSuccessError handles errors related to retrieving monthly successful transaction statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transaction status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionResponseMonthStatusSuccess pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleMonthTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseMonthStatusSuccess](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionSuccess, fields...)
}

// HandleYearlyTransactionStatusSuccessError handles errors related to retrieving yearly successful transaction statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transaction status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionResponseYearStatusSuccess pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleYearlyTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseYearStatusSuccess](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionSuccess, fields...)
}

// HandleMonthTransactionStatusFailedError handles errors related to retrieving monthly failed transaction statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transaction status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionResponseMonthStatusFailed pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleMonthTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseMonthStatusFailed](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionFailed, fields...)
}

// HandleYearlyTransactionStatusFailedError handles errors related to retrieving yearly failed transaction statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transaction status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionResponseYearStatusFailed pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleYearlyTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseYearStatusFailed](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionFailed, fields...)
}

// HandleMonthlyPaymentMethodsError handles errors related to retrieving monthly payment methods.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the payment methods.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionMonthMethodResponse pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleMonthlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyPaymentMethods, fields...)
}

// HandleYearlyPaymentMethodsError handles errors related to retrieving yearly payment methods.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the payment methods.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionYearMethodResponse pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleYearlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyPaymentMethods, fields...)
}

// HandleMonthlyAmountsError handles errors related to retrieving monthly transaction amounts.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the monthly amounts.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionMonthAmountResponse pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleMonthlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthAmountResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmounts, fields...)
}

// HandleYearlyAmountsError handles errors related to retrieving yearly transaction amounts.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the yearly amounts.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransactionYearlyAmountResponse pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transactionStatisticError) HandleYearlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyAmountResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmounts, fields...)
}

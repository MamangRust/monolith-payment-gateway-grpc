package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transactionStatisticByCardError handles errors related to transaction statistics by card number.
type transactionStatisticByCardError struct {
	logger logger.LoggerInterface
}

// NewTransactionStatisticByCardError returns a new instance of transactionStatisticByCardError with the given logger.
// This function returns a pointer to the transactionStatisticByCardError struct, which implements the TransactionStatisticByCardError interface.
// It is used for handling errors related to transaction statistics by card number, ensuring that they are logged appropriately.
func NewTransactionStatisticByCardError(logger logger.LoggerInterface) TransactionStatisticByCardErrorHandler {
	return &transactionStatisticByCardError{logger: logger}
}

// HandleMonthTransactionStatusSuccessByCardNumberError handles the error of retrieving monthly successful transactions by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionResponseMonthStatusSuccess containing the details of the successful operation, and an ErrorResponse indicating failure.
func (e *transactionStatisticByCardError) HandleMonthTransactionStatusSuccessByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseMonthStatusSuccess](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionSuccessByCard, fields...)
}

// HandleYearlyTransactionStatusSuccessByCardNumberError handles the error of retrieving yearly successful transactions by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionResponseYearStatusSuccess containing the details of the successful operation, and an ErrorResponse indicating failure.
func (e *transactionStatisticByCardError) HandleYearlyTransactionStatusSuccessByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseYearStatusSuccess](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionSuccessByCard, fields...)
}

// HandleMonthTransactionStatusFailedByCardNumberError handles the error of retrieving monthly failed transactions by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionResponseMonthStatusFailed containing the details of the successful operation, and an ErrorResponse indicating failure.
func (e *transactionStatisticByCardError) HandleMonthTransactionStatusFailedByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseMonthStatusFailed](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthTransactionFailedByCard, fields...)
}

// HandleYearlyTransactionStatusFailedByCardNumberError handles the error of retrieving yearly failed transactions by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionResponseYearStatusFailed containing the details of the successful operation, and an ErrorResponse indicating failure.
func (e *transactionStatisticByCardError) HandleYearlyTransactionStatusFailedByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponseYearStatusFailed](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearTransactionFailedByCard, fields...)
}

// HandleMonthlyPaymentMethodsByCardNumberError processes errors during retrieval of monthly payment methods by card number.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionMonthMethodResponse with error details and a standardized ErrorResponse.
func (e *transactionStatisticByCardError) HandleMonthlyPaymentMethodsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthMethodResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyPaymentMethodsByCard, fields...)
}

// HandleYearlyPaymentMethodsByCardNumberError processes errors during retrieval of yearly payment methods by card number.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionYearMethodResponse with error details and a standardized ErrorResponse.
func (e *transactionStatisticByCardError) HandleYearlyPaymentMethodsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearMethodResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyPaymentMethodsByCard, fields...)
}

// HandleMonthlyAmountsByCardNumberError processes errors during retrieval of monthly amounts by card number.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionMonthAmountResponse with error details and a standardized ErrorResponse.
func (e *transactionStatisticByCardError) HandleMonthlyAmountsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthAmountResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmountsByCard, fields...)
}

// HandleYearlyAmountsByCardNumberError processes errors during retrieval of yearly amounts by card number.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransactionYearlyAmountResponse with error details and a standardized ErrorResponse.
func (e *transactionStatisticByCardError) HandleYearlyAmountsByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyAmountResponse](e.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmountsByCard, fields...)
}

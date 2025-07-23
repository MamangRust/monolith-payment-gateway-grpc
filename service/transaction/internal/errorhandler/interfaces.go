package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TransactionQueryErrorHandler defines the interface for handling errors during transaction query operations.
type TransactionQueryErrorHandler interface {
	// HandleRepositoryPaginationError processes pagination errors from the repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the pagination operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A slice of TransactionResponse pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
	// when retrieving deleted transactions. It logs the error, updates the trace span,
	// and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the pagination operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A slice of TransactionResponseDeleteAt pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	// HandleRepositorySingleError processes single-result errors from the transaction repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the single-result operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A TransactionResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
	// HanldeRepositoryListError processes list errors from the transaction repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the list operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A slice of TransactionResponse pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the list failure.
	HanldeRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransactionResponse, *response.ErrorResponse)
}

// TransactionCommandErrorHandler provides methods for handling errors related to transaction commands.
type TransactionCommandErrorHandler interface {
	// HandleInvalidParseTimeError handles errors related to parsing time values.
	// It constructs an appropriate error response for cases where the provided
	// time value is invalid or cannot be parsed. The method logs the error
	// and returns a TransactionResponse and an ErrorResponse indicating a bad request.
	//
	// Parameters:
	//   - err: the error encountered during time parsing.
	//   - method: the name of the method where the error occurred.
	//   - tracePrefix: the prefix used for tracing the error.
	//   - span: the trace span for the request.
	//   - status: a pointer to a string representing the status of the operation.
	//   - rawTime: the raw time string that failed to parse.
	//   - fields: additional context fields for logging.
	//
	// Returns:
	//   - *response.TransactionResponse: the response for the transaction operation (nil in error cases).
	//   - *response.ErrorResponse: the constructed error response indicating the failure.
	HandleInvalidParseTimeError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		rawTime string,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
	// HandleInsufficientBalanceError processes errors related to insufficient balance.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - cardNumber: The card number associated with the insufficient balance.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransactionResponse with error details and a standardized ErrorResponse.
	HandleInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		cardNumber string,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
	// HandleRepositorySingleError processes single-result errors from the repository.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the single-result operation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - defaultErr: Default error message.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransactionResponse with error details and a standardized ErrorResponse.
	HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	// HandleCreateTransactionError processes errors related to creating a transaction.
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
	//   - A TransactionResponse with error details and a standardized ErrorResponse.
	HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	// HandleUpdateTransactionError processes errors related to updating a transaction.
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
	//   - A TransactionResponse with error details and a standardized ErrorResponse.
	HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	// HandleTrashedTransactionError processes errors that occur during the trashing of a transaction.
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
	//   - A TransactionResponse containing the details of the transaction trashing failure.
	//   - A standardized ErrorResponse detailing the trashing failure.
	HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	// HandleRestoreTransactionError processes errors that occur during the restore of a transaction.
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
	//   - A TransactionResponse containing the details of the transaction restore failure.
	//   - A standardized ErrorResponse detailing the restore failure.
	HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	// HandleDeleteTransactionPermanentError processes errors that occur during the permanent deletion of a transaction.
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
	//   - A boolean indicating whether the deletion failed.
	//   - A standardized ErrorResponse detailing the deletion failure.
	HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleRestoreAllTransactionError processes errors that occur during the restore of all transactions.
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
	//   - A boolean indicating whether the restoration was successful.
	//   - A standardized ErrorResponse detailing the restoration error.
	HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleDeleteAllTransactionPermanentError processes errors that occur during the permanent deletion of all transactions.
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
	//   - A boolean indicating whether the deletion failed.
	//   - A standardized ErrorResponse detailing the deletion failure.
	HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

// TransactionStatisticErrorHandler is an interface for handling errors related to transaction statistics.
type TransactionStatisticErrorHandler interface {
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
	HandleMonthTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
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
	HandleYearlyTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
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
	HandleMonthTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
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
	HandleYearlyTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)
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
	HandleMonthlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
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
	HandleYearlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
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
	HandleMonthlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
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
	HandleYearlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

// TransactionStatisticByCardErrorHandler provides methods for handling errors related to transaction statistics by card number.
type TransactionStatisticByCardErrorHandler interface {
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
	HandleMonthTransactionStatusSuccessByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
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
	HandleYearlyTransactionStatusSuccessByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,

		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
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
	HandleMonthTransactionStatusFailedByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
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
	HandleYearlyTransactionStatusFailedByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)
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
	HandleMonthlyPaymentMethodsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
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
	HandleYearlyPaymentMethodsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
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
	HandleMonthlyAmountsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
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
	HandleYearlyAmountsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

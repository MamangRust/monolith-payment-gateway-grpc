package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TransferQueryErrorHandler is an interface that defines methods for handling errors in the transfer query service.
type TransferQueryErrorHandler interface {
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
	//   - A slice of TransferResponse pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransferResponse, *int, *response.ErrorResponse)
	// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
	// when retrieving deleted transfers. It logs the error, updates the trace span,
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
	//   - A slice of TransferResponseDeleteAt pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)
	// HandleRepositorySingleError processes single-result errors from the repository.
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
	//   - A TransferResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	// HanldeRepositoryListError processes list errors from the repository.
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
	//   - A slice of TransferResponse pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the list failure.
	HanldeRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransferResponse, *response.ErrorResponse)
}

// TransferCommandErrorHandler is an interface that defines methods for handling errors in the transfer command service.
type TransferCommandErrorHandler interface {
	// HandleSenderInsufficientBalanceError processes errors related to insufficient balance
	// of the sender card. It logs the error, records it to the trace span, and returns a
	// standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - senderCardNumber: The card number associated with the insufficient balance.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransferResponse with error details and a standardized ErrorResponse.
	HandleSenderInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		senderCardNumber string,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	// HandleReceiverInsufficientBalanceError processes errors related to insufficient balance
	// of the receiver card. It logs the error, records it to the trace span, and returns a
	// standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - receiverCardNumber: The card number associated with the insufficient balance.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransferResponse with error details and a standardized ErrorResponse.
	HandleReceiverInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		receiverCardNumber string,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	// HandleRepositorySingleError processes single-result errors from the repository.
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
	//   - A TransferResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	// HandleCreateTransferError processes errors related to creating a transfer.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransferResponse with error details and a standardized ErrorResponse.
	HandleCreateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)

	// HandleUpdateTransferError processes errors related to updating a transfer.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransferResponse with error details and a standardized ErrorResponse.
	HandleUpdateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)
	// HandleTrashedTransferError processes errors related to trashing a transfer.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A TransferResponse with error details and a standardized ErrorResponse.
	HandleTrashedTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponseDeleteAt, *response.ErrorResponse)
	// HandleRestoreTransferError processes errors related to the restoration of a transfer.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the transfer restoration.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The trace span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A TransferResponse with error details and a standardized ErrorResponse.
	HandleRestoreTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)
	// HandleDeleteTransferPermanentError processes errors that occur during the permanent deletion of a transfer.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A boolean indicating whether the deletion was successful.
	//   - A standardized ErrorResponse detailing the deletion failure.
	HandleDeleteTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleRestoreAllTransferError processes errors that occur during the restoration of all transfers.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the transfer restoration.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The trace span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A boolean indicating whether the restoration was successful.
	//   - A standardized ErrorResponse detailing the restoration failure.
	HandleRestoreAllTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleDeleteAllTransferPermanentError processes errors that occur during the permanent deletion of all transfers.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A boolean indicating whether the deletion was successful.
	//   - A standardized ErrorResponse detailing the deletion failure.
	HandleDeleteAllTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

// TransferStatisticErrorHandler provides methods for handling errors related to transfer statistics.
type TransferStatisticErrorHandler interface {
	// HandleMonthTransferStatusSuccessError handles errors related to retrieving monthly successful transfer statuses.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error encountered while processing the transfer status.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The tracing span to record error details.
	//   - status: A pointer to a string representing the status, which will be updated as necessary.
	//   - fields: Additional contextual fields for logging purposes.
	//
	// Returns:
	//   - A slice of TransferResponseMonthStatusSuccess pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the failure, if any.
	HandleMonthTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	// HandleYearTransferStatusSuccessError handles errors related to retrieving yearly successful transfer statuses.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error encountered while processing the transfer status.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The tracing span to record error details.
	//   - status: A pointer to a string representing the status, which will be updated as necessary.
	//   - fields: Additional contextual fields for logging purposes.
	//
	// Returns:
	//   - A slice of TransferResponseYearStatusSuccess pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the failure, if any.
	HandleYearTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	// HandleMonthTransferStatusFailedError handles errors related to retrieving monthly failed transfer statuses.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error encountered while processing the transfer status.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The tracing span to record error details.
	//   - status: A pointer to a string representing the status, which will be updated as necessary.
	//   - fields: Additional contextual fields for logging purposes.
	//
	// Returns:
	//   - A slice of TransferResponseMonthStatusFailed pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the failure, if any.
	HandleMonthTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	// HandleYearTransferStatusFailedError handles errors related to retrieving yearly failed transfer statuses.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error encountered while processing the transfer status.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The tracing span to record error details.
	//   - status: A pointer to a string representing the status, which will be updated as necessary.
	//   - fields: Additional contextual fields for logging purposes.
	//
	// Returns:
	//   - A slice of TransferResponseYearStatusFailed pointers if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the failure, if any.
	HandleYearTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)
	// HandleMonthlyTransferAmountsError processes errors during retrieval of monthly transfer amounts.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TransferMonthAmountResponse with error details and a standardized ErrorResponse.
	HandleMonthlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	// HandleYearlyTransferAmountsError processes errors during retrieval of yearly transfer amounts.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TransferYearAmountResponse with error details and a standardized ErrorResponse.
	HandleYearlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

// TransferStatisticByCardErrorHandler provides methods for handling errors related to transfer statistics by card number.
type TransferStatisticByCardErrorHandler interface {
	// HandleMonthTransferStatusSuccessByCardNumberError handles the error of retrieving monthly successful transfers by card number.
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
	//   - A slice of TransferResponseMonthStatusSuccess containing the details of the successful operation, and an ErrorResponse indicating failure.
	HandleMonthTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	// HandleYearTransferStatusSuccessByCardNumberError handles the error of retrieving yearly successful transfers by card number.
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
	//   - A slice of TransferResponseYearStatusSuccess containing the details of the successful operation, and an ErrorResponse indicating failure.
	HandleYearTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	// HandleMonthTransferStatusFailedByCardNumberError handles the error of retrieving monthly failed transfers by card number.
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
	//   - A slice of TransferResponseMonthStatusFailed containing the details of the successful operation, and an ErrorResponse indicating failure.
	HandleMonthTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	// HandleYearTransferStatusFailedByCardNumberError handles the error of retrieving yearly failed transfers by card number.
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
	//   - A slice of TransferResponseYearStatusFailed containing the details of the successful operation, and an ErrorResponse indicating failure.
	HandleYearTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)
	// HandleMonthlyTransferAmountsBySenderError processes errors during retrieval of monthly transfer amounts by sender's card.
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
	//   - A slice of TransferMonthAmountResponse with error details and a standardized ErrorResponse.
	HandleMonthlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	// HandleYearlyTransferAmountsBySenderError processes errors during retrieval of yearly transfer amounts by sender's card.
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
	//   - A slice of TransferYearAmountResponse with error details and a standardized ErrorResponse.
	HandleYearlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
	// HandleMonthlyTransferAmountsByReceiverError processes errors during retrieval of monthly transfer amounts by receiver's card.
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
	//   - A slice of TransferMonthAmountResponse with error details and a standardized ErrorResponse.
	HandleMonthlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	// HandleYearlyTransferAmountsByReceiverError processes errors during retrieval of yearly transfer amounts by receiver's card.
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
	//   - A slice of TransferYearAmountResponse with error details and a standardized ErrorResponse.
	HandleYearlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

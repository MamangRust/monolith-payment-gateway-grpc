package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TopupQueryErrorHandler is an interface for handling errors in topup query operations.
type TopupQueryErrorHandler interface {
	// HandleRepositoryPaginationError processes pagination errors from the repository.
	// It logs the error, updates the trace span, and returns a standardized error response.
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
	//   - A slice of TopupResponse pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TopupResponse, *int, *response.ErrorResponse)
	// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
	// when retrieving deleted topup documents. It logs the error, updates the trace span,
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
	//   - A slice of TopupResponseDeleteAt pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
}

// TopupStatisticErrorHandler is an interface for handling topup statistic errors.
type TopupStatisticErrorHandler interface {
	// HandleMonthTopupStatusSuccess processes the successful retrieval of monthly topup status.
	// It logs the success information, records it to the trace span, and returns a structured response.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseMonthStatusSuccess containing the details of the successful operation,
	//     and a nil ErrorResponse indicating success.
	HandleMonthTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	// HandleYearlyTopupStatusSuccess processes the successful retrieval of yearly topup status.
	// It logs the success information, records it to the trace span, and returns a structured response.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseYearStatusSuccess containing the details of the successful operation,
	//     and a nil ErrorResponse indicating success.
	HandleYearlyTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	// HandleMonthTopupStatusFailed processes the failure to retrieve monthly topup status.
	// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - err: The error, if any, encountered during the process.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseMonthStatusFailed containing the details of the failed operation,
	//     and an ErrorResponse containing more information about the failure.
	HandleMonthTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	// HandleYearlyTopupStatusFailed processes the failure to retrieve yearly topup status.
	// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseYearStatusFailed containing the details of the failed operation,
	//     and an ErrorResponse containing more information about the failure.
	HandleYearlyTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)
	// HandleMonthlyTopupMethods processes the failure to retrieve monthly topup methods.
	// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupMonthMethodResponse containing the details of the failed operation,
	//     and an ErrorResponse containing more information about the failure.
	HandleMonthlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	// HandleYearlyTopupMethods processes the failure to retrieve yearly topup methods.
	// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupYearlyMethodResponse containing the details of the failed operation,
	//     and an ErrorResponse containing more information about the failure.
	HandleYearlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.
		TopupYearlyMethodResponse, *response.ErrorResponse)
	// HandleMonthlyTopupAmounts processes the retrieval of monthly topup amounts.
	// It logs the error information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupMonthAmountResponse containing the details of the failed operation,
	//     and an ErrorResponse containing more information about the failure.
	HandleMonthlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	// HandleYearlyTopupAmounts processes the retrieval of yearly topup amounts.
	// It logs the error information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupYearlyAmountResponse containing the details of the failed operation,
	//     and an ErrorResponse containing more information about the failure.
	HandleYearlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

// TopupStatisticByCardErrorHandler is an interface for handling errors in topup statistic by card number.
type TopupStatisticByCardErrorHandler interface {
	// HandleMonthTopupStatusSuccessByCardNumber handles the successful retrieval of monthly topup status by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating success.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseMonthStatusSuccess containing the details of the successful operation, and a nil ErrorResponse indicating success.
	HandleMonthTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	// HandleYearlyTopupStatusSuccessByCardNumber handles the successful retrieval of yearly topup status by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating success.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseYearStatusSuccess containing the details of the successful operation, and a nil ErrorResponse indicating success.
	HandleYearlyTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	// HandleMonthTopupStatusFailedByCardNumber handles the failure to retrieve monthly topup status by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseMonthStatusFailed containing the details of the failed operation, and an ErrorResponse containing more information about the failure.
	HandleMonthTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	// HandleYearlyTopupStatusFailedByCardNumber handles the failure to retrieve yearly topup status by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating failure.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the failure is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the failure.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupResponseYearStatusFailed containing the details of the failed operation, and an ErrorResponse containing more information about the failure.

	HandleYearlyTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	// HandleMonthlyTopupMethodsByCardNumber handles the successful retrieval of monthly topup methods by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating success.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupMonthMethodResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.

	HandleMonthlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	// HandleYearlyTopupMethodsByCardNumber handles the successful retrieval of yearly topup methods by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating success.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupYearlyMethodResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
	HandleYearlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	// HandleMonthlyTopupAmountsByCardNumber handles the successful retrieval of monthly topup amounts by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating success.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupMonthAmountResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
	HandleMonthlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	// HandleYearlyTopupAmountsByCardNumber handles the successful retrieval of yearly topup amounts by card number.
	// It logs the information, records it to the trace span, and returns a structured response indicating success.
	// Parameters:
	//   - err: The error, if any, encountered during the process.
	//   - method: The name of the method where the success is recorded.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the success.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of TopupYearlyAmountResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
	HandleYearlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

// TopupCommandErrorHandler is an interface for handling errors related to topup commands.
type TopupCommandErrorHandler interface {
	// HandleInvalidParseTimeError handles errors related to parsing time values.
	// It constructs an appropriate error response for cases where the provided
	// time value is invalid or cannot be parsed. The method logs the error
	// and returns a TopupResponse and an ErrorResponse indicating a bad request.
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
	//   - *response.TopupResponse: the response for the topup operation (nil in error cases).
	//   - *response.ErrorResponse: the constructed error response indicating the failure.
	HandleInvalidParseTimeError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		rawTime string,
		fields ...zap.Field,
	) (*response.TopupResponse, *response.ErrorResponse)
	// HandleCreateTopupError processes errors during the topup creation process.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the topup creation.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A TopupResponse pointer if the operation is successful, otherwise nil.
	//   - A standardized ErrorResponse describing the creation failure.
	HandleCreateTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponse, *response.ErrorResponse)
	// HandleUpdateTopupError processes errors during the topup update process.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the topup update.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A TopupResponse pointer if the operation is successful, otherwis
	HandleUpdateTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponse, *response.ErrorResponse)
	// HandleTrashedTopupError processes errors during the topup trash process.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the topup trash.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A TopupResponseDeleteAt pointer if the operation is successful, otherwise nil.
	//   - A standardized ErrorResponse describing the trash failure.
	HandleTrashedTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	// HandleRestoreTopupError processes errors during the topup restore process.
	// It logs the error, updates the trace span, and returns a standardized error response.
	//
	// Parameters:
	//   - err: The error that occurred during the topup restore.
	//   - method: The name of the method where the error originated.
	//   - tracePrefix: A prefix used for generating the trace ID.
	//   - span: The tracing span used for recording error details.
	//   - status: A pointer to a string that will be updated with the formatted status.
	//   - fields: Additional contextual fields for logging.
	//
	// Returns:
	//   - A TopupResponse pointer if the operation is successful, otherwise nil.
	//   - A standardized ErrorResponse describing the restore failure.
	HandleRestoreTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponse, *response.ErrorResponse)
	// HandleDeleteTopupPermanentError processes errors that occur during the permanent deletion of a Topup.
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
	//   - A standardized ErrorResponse detailing the deletion error.
	HandleDeleteTopupPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleRestoreAllTopupError processes errors that occur during the restoration of all Topups.
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
	//   - A boolean indicating whether the restoration was successful.
	//   - A standardized ErrorResponse detailing the restoration error.
	HandleRestoreAllTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleDeleteAllTopupPermanentError processes errors that occur during the permanent deletion of all Topups.
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
	//   - A standardized ErrorResponse detailing the deletion error.
	HandleDeleteAllTopupPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

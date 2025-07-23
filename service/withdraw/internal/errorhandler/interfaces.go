package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// WithdrawQueryErrorHandler provides methods for handling errors related to withdraw queries.
type WithdrawQueryErrorHandler interface {
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
	//   - A slice of WithdrawResponse pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)
	// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
	// when retrieving deleted withdraw documents. It logs the error, updates the trace span,
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
	//   - A slice of WithdrawResponseDeleteAt pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)
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
	//   - A WithdrawResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.WithdrawResponse, *response.ErrorResponse)
}

// WithdrawCommandErrorHandler provides methods for handling errors related to withdraw operations.
type WithdrawCommandErrorHandler interface {
	// HandleInsufficientBalanceError handles errors related to insufficient balance during a withdraw operation.
	// It logs the error, records it in the trace span, and returns a standardized error response.
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
	//   - A WithdrawResponse with error details and a standardized ErrorResponse.
	HandleInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		cardNumber string,
		fields ...zap.Field,
	) (*response.WithdrawResponse, *response.ErrorResponse)
	// HandleRepositorySingleError processes single-result errors from the repository.
	// It logs the error, records it to the trace span, and returns a standardized error response.
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
	//   - A WithdrawResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.WithdrawResponse, *response.ErrorResponse)
	// HandleCreateWithdrawError processes errors during creation of a new withdraw.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	// Args:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A WithdrawResponse with error details and a standardized ErrorResponse.
	HandleCreateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse)
	// HandleUpdateWithdrawError processes errors during updating a withdraw.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	// Args:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A WithdrawResponse with error details and a standardized ErrorResponse.
	HandleUpdateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse)
	// HandleTrashedWithdrawError processes errors during trashing a withdraw.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	// Args:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A WithdrawResponse with error details and a standardized ErrorResponse.
	HandleTrashedWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponseDeleteAt, *response.ErrorResponse)
	// HandleRestoreWithdrawError processes errors during restore of a withdraw.
	// It logs the error, records it to the trace span, and returns a standardized error response.
	// Args:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A WithdrawResponse with error details and a standardized ErrorResponse.
	HandleRestoreWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse)
	// HandleDeleteWithdrawPermanentError processes errors during the permanent deletion of a withdraw.
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
	//   - A boolean indicating whether the deletion was successful.
	//   - A standardized ErrorResponse detailing the deletion error.
	HandleDeleteWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleRestoreAllWithdrawError processes errors during the restoration of all withdraw.
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
	HandleRestoreAllWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	// HandleDeleteAllWithdrawPermanentError handles errors that occur during the permanent deletion of all withdraws.
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
	HandleDeleteAllWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

// WithdrawStatisticErrorHandler is an interface for handling errors in the withdraw statistic service.
type WithdrawStatisticErrorHandler interface {
	// HandleMonthWithdrawStatusSuccessError processes errors during the retrieval of monthly successful withdraw status.
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
	//   - A slice of WithdrawResponseMonthStatusSuccess with error details and a standardized ErrorResponse.
	HandleMonthWithdrawStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	// HandleYearWithdrawStatusSuccessError processes errors during the retrieval of yearly successful withdraw status.
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
	//   - A slice of WithdrawResponseYearStatusSuccess with error details and a standardized ErrorResponse.
	HandleYearWithdrawStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	// HandleMonthWithdrawStatusFailedError processes errors during the retrieval of monthly failed withdraw status.
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
	//   - A slice of WithdrawResponseMonthStatusFailed with error details and a standardized ErrorResponse.
	HandleMonthWithdrawStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	// HandleYearWithdrawStatusFailedError processes errors during the retrieval of yearly failed withdraw status.
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
	//   - A slice of WithdrawResponseYearStatusFailed with error details and a standardized ErrorResponse.
	HandleYearWithdrawStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)
	// HandleMonthlyWithdrawAmountsError processes errors during the retrieval of monthly withdraw amounts.
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
	//   - A slice of WithdrawMonthlyAmountResponse with error details and a standardized ErrorResponse.
	HandleMonthlyWithdrawAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	// HandleYearlyWithdrawAmountsError processes errors during the retrieval of yearly withdraw amounts.
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
	//   - A slice of WithdrawYearlyAmountResponse with error details and a standardized ErrorResponse.
	HandleYearlyWithdrawAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
}

// WithdrawStatisticByCardErrorHandler provides methods for handling errors in the withdrawal statistics by card number.
type WithdrawStatisticByCardErrorHandler interface {
	// HandleMonthWithdrawStatusSuccessByCardNumberError processes errors during retrieval of monthly successful withdraw status by card number.
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
	//   - A slice of WithdrawResponseMonthStatusSuccess with error details and a standardized ErrorResponse.
	HandleMonthWithdrawStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	// HandleYearWithdrawStatusSuccessByCardNumberError processes errors during retrieval of yearly successful withdraw status by card number.
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
	//   - A slice of WithdrawResponseYearStatusSuccess with error details and a standardized ErrorResponse.
	HandleYearWithdrawStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)

	// HandleMonthWithdrawStatusFailedByCardNumberError processes errors during retrieval of monthly failed withdraw status by card number.
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
	//   - A slice of WithdrawResponseMonthStatusFailed with error details and a standardized ErrorResponse.
	HandleMonthWithdrawStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	// HandleYearWithdrawStatusFailedByCardNumberError processes errors during retrieval of yearly failed withdraw status by card number.
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
	//   - A slice of WithdrawResponseYearStatusFailed with error details and a standardized ErrorResponse.
	HandleYearWithdrawStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)

	// HandleMonthlyWithdrawsAmountByCardNumberError processes errors during retrieval of monthly withdraw amounts by card number.
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
	//   - A slice of WithdrawMonthlyAmountResponse with error details and a standardized ErrorResponse.
	HandleMonthlyWithdrawsAmountByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	// HandleYearlyWithdrawsAmountByCardNumberError processes errors during retrieval of yearly withdraw amounts by card number.
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
	//   - A slice of WithdrawYearlyAmountResponse with error details and a standardized ErrorResponse.
	HandleYearlyWithdrawsAmountByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
}

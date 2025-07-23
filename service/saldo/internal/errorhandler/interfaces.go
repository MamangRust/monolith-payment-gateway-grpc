package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// SaldoCommandErrorHandler is an interface that defines methods for handling errors related to Saldo commands.
type SaldoCommandErrorHandler interface {
	// HandleFindCardByNumberError processes errors during card lookup by card number
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
	//   - A SaldoResponse with error details and a standardized ErrorResponse.
	HandleFindCardByNumberError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	// HandleCreateSaldoError processes errors during card creation
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
	//   - A SaldoResponse with error details and a standardized ErrorResponse.
	HandleCreateSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	// HandleUpdateSaldoError processes errors during card update
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
	//   - A SaldoResponse with error details and a standardized ErrorResponse.
	HandleUpdateSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	// HandleTrashSaldoError processes errors during Saldo soft deletion (trashing)
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
	//   - A SaldoResponse with error details and a standardized ErrorResponse.
	HandleTrashSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponseDeleteAt, *response.ErrorResponse)
	// HandleRestoreSaldoError processes errors during Saldo restore (undoing trashing)
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
	//   - A SaldoResponse with error details and a standardized ErrorResponse.
	HandleRestoreSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	// HandleDeleteSaldoPermanentError processes errors during the permanent deletion of a Saldo.
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
	HandleDeleteSaldoPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	// HandleRestoreAllSaldoError processes errors that occur during the restoration of all Saldo.
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
	HandleRestoreAllSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	// HandleDeleteAllSaldoPermanentError processes errors that occur during the permanent deletion of all Saldo.
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
	HandleDeleteAllSaldoPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

// SaldoQueryErrorHandler provides methods for handling errors related to Saldo queries.
type SaldoQueryErrorHandler interface {
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
	//   - A slice of SaldoResponse pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.SaldoResponse, *int, *response.ErrorResponse)
	// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
	// when retrieving deleted saldo documents. It logs the error, updates the trace span,
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
	//   - A slice of SaldoResponseDeleteAt pointers if successful, otherwise nil.
	//   - A pointer to an integer representing additional pagination details, otherwise nil.
	//   - A standardized ErrorResponse describing the pagination failure.
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
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
	//   - A SaldoResponse pointer if successful, otherwise nil.
	//   - A standardized ErrorResponse describing the single-result failure.
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
}

// SaldoStatisticErrorHandler provides methods for handling errors related to Saldo statistic queries.
type SaldoStatisticErrorHandler interface {
	// HandleMonthlyTotalSaldoBalanceError processes errors during the retrieval of a monthly total saldo balance.
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
	//   - A slice of SaldoMonthTotalBalanceResponse with error details and a standardized ErrorResponse.
	HandleMonthlyTotalSaldoBalanceError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse)
	// HandleYearlyTotalSaldoBalanceError processes errors during the retrieval of a yearly total saldo balance.
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
	//   - A slice of SaldoYearTotalBalanceResponse with error details and a standardized ErrorResponse.
	HandleYearlyTotalSaldoBalanceError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse)
	// HandleMonthlySaldoBalancesError processes errors during the retrieval of monthly saldo balances.
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
	//   - A slice of SaldoMonthBalanceResponse with error details and a standardized ErrorResponse.
	HandleMonthlySaldoBalancesError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse)
	// HandleYearlySaldoBalancesError processes errors during the retrieval of yearly saldo balances.
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
	//   - A slice of SaldoYearBalanceResponse with error details and a standardized ErrorResponse.
	HandleYearlySaldoBalancesError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse)
}

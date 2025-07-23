package errorhandler

import (
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transactionCommandError is a struct that implements the TransactionCommandError interface.
type transactionCommandError struct {
	logger logger.LoggerInterface
}

// NewTransactionCommandError initializes a new transactionCommandError with the provided logger.
// It returns an instance of the transactionCommandError struct.
func NewTransactionCommandError(logger logger.LoggerInterface) TransactionCommandErrorHandler {
	return &transactionCommandError{logger: logger}
}

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
func (t *transactionCommandError) HandleInvalidParseTimeError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	rawTime string,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {
	errResp := &response.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse the given time value",
		Status:  "invalid_parse_time",
	}

	return sharederrorhandler.HandleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, "Invalid parse time", span, status, errResp, fields...)
}

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
func (t *transactionCommandError) HandleInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	cardNumber string,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

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
func (t *transactionCommandError) HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

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
func (t *transactionCommandError) HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedCreateTransaction, fields...)
}

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
func (t *transactionCommandError) HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedUpdateTransaction, fields...)
}

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
func (t *transactionCommandError) HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedTrashedTransaction, fields...)
}

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
func (t *transactionCommandError) HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreTransaction, fields...)
}

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
func (t *transactionCommandError) HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteTransactionPermanent, fields...)
}

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
func (t *transactionCommandError) HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreAllTransactions, fields...)
}

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
func (t *transactionCommandError) HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteAllTransactionsPermanent, fields...)
}

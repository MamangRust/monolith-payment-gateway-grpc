package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transactionQueryError represents the error handler for transaction query operations.
type transactionQueryError struct {
	logger logger.LoggerInterface
}

// NewTransactionQueryError initializes a new transactionQueryError with the provided logger.
// It returns an instance of the transactionQueryError struct.
func NewTransactionQueryError(logger logger.LoggerInterface) TransactionQueryErrorHandler {
	return &transactionQueryError{logger: logger}
}

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
func (t *transactionQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindAllTransactions, fields...)
}

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
func (t *transactionQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransactionResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

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
func (t *transactionQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

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
func (t *transactionQueryError) HanldeRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

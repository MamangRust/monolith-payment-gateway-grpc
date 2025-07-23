package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transferQueryError is a struct that implements the TransferQueryError interface.
type transferQueryError struct {
	logger logger.LoggerInterface
}

// NewTransferQueryError initializes a new transferQueryError with the provided logger.
// It returns an instance of the transferQueryError struct.
func NewTransferQueryError(logger logger.LoggerInterface) TransferQueryErrorHandler {
	return &transferQueryError{
		logger: logger,
	}
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
//   - A slice of TransferResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (t *transferQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransferResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedFindAllTransfers, fields...)
}

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
func (t *transferQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransferResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

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
func (t *transferQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

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
func (t *transferQueryError) HanldeRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

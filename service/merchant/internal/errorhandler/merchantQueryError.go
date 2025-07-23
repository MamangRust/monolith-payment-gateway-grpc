package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantQueryError is a struct that implements the MerchantQueryErrorHandler interface
type merchantQueryError struct {
	logger logger.LoggerInterface
}

// NewMerchantQueryError initializes a new merchantQueryError with the provided logger.
// It returns an instance of the merchantQueryError struct.
func NewMerchantQueryError(logger logger.LoggerInterface) MerchantQueryErrorHandler {
	return &merchantQueryError{
		logger: logger,
	}
}

// HandleRepositoryPaginationError processes pagination errors from the repository.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Args:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A slice of MerchantResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantResponse](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...,
	)
}

// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
// when retrieving deleted merchant documents.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Args:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
//
// Returns:
//   - A slice of MerchantResponseDeleteAt pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantResponseDeleteAt](
		e.logger, err, method, tracePrefix, span, status, errResp, fields...,
	)
}

// HandleRepositoryListError processes list errors from the repository.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Args:
//   - err: The error that occurred during the list operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A slice of MerchantResponse pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the list failure.
func (e *merchantQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponse](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...,
	)
}

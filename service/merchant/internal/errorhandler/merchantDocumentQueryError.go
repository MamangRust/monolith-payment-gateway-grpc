package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantDocumentQueryError is a struct that implements the MerchantDocumentQueryErrorHandler interface
type merchantDocumentQueryError struct {
	logger logger.LoggerInterface
}

// NewMerchantDocumentQueryError returns a new instance of MerchantDocumentQueryError with the given logger.
func NewMerchantDocumentQueryError(logger logger.LoggerInterface) MerchantDocumentQueryErrorHandler {
	return &merchantDocumentQueryError{
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
//   - A slice of MerchantDocumentResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantDocumentQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantDocumentResponse](e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...)
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
//
// Returns:
//   - A slice of MerchantDocumentResponseDeleteAt pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantDocumentQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantDocumentResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...)
}

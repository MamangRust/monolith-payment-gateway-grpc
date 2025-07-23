package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantTransactionError is a struct that implements the MerchantTransactionErrorHandler interface
type merchantTransactionError struct {
	logger logger.LoggerInterface
}

// NewMerchantTransactionError initializes a new merchantTransactionError with the provided logger.
// It returns an instance of the merchantTransactionError struct.
func NewMerchantTransactionError(logger logger.LoggerInterface) MerchantTransactionErrorHandler {
	return &merchantTransactionError{
		logger: logger,
	}
}

// HandleRepositoryAllError processes pagination errors from the repository when retrieving all transactions.
// It logs the error, records it to the trace span, and returns a standardized error response.
//
// Args:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantTransactionResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantTransactionError) HandleRepositoryAllError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantTransactionResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		merchant_errors.ErrFailedFindAllTransactions,
		fields...,
	)
}

// HandleRepositoryByMerchantError processes pagination errors from the repository when retrieving
// transactions by merchant ID.
// It logs the error, records it to the trace span, and returns a standardized error response.
//
// Args:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantTransactionResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantTransactionError) HandleRepositoryByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantTransactionResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		merchant_errors.ErrFailedFindAllTransactionsByMerchant,
		fields...,
	)
}

// HandleRepositoryByApiKeyError processes pagination errors from the repository when retrieving
// transactions by API key.
// It logs the error, records it to the trace span, and returns a standardized error response.
//
// Args:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantTransactionResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *merchantTransactionError) HandleRepositoryByApiKeyError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantTransactionResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		merchant_errors.ErrFailedFindAllTransactionsByApikey,
		fields...,
	)
}

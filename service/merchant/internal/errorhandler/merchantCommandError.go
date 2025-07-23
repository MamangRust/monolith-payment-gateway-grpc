package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantCommandError is a struct that implements the MerchantQueryErrorHandler interface.
type merchantCommandError struct {
	logger logger.LoggerInterface
}

// NewMerchantCommandError initializes a new merchantCommandError with the provided logger.
// It returns an instance of the merchantCommandError struct.

func NewMerchantCommandError(logger logger.LoggerInterface) MerchantCommandErrorHandler {
	return &merchantCommandError{
		logger: logger,
	}
}

// HandleCreateMerchantError processes errors that occur during the creation of a merchant.
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
//   - A MerchantResponse containing the details of the merchant creation failure.
//   - A standardized ErrorResponse detailing the creation failure.
func (e *merchantCommandError) HandleCreateMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedCreateMerchant,
		fields...,
	)
}

// HandleUpdateMerchantError processes errors that occur during the update of a merchant.
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
//   - A MerchantResponse containing the details of the merchant update failure.
//   - A standardized ErrorResponse detailing the update failure.
func (e *merchantCommandError) HandleUpdateMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedUpdateMerchant,
		fields...,
	)
}

// HandleUpdateMerchantStatusError processes errors that occur during the update of a merchant's status.
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
//   - A MerchantResponse containing the details of the merchant status update failure.
//   - A standardized ErrorResponse detailing the update failure.
func (e *merchantCommandError) HandleUpdateMerchantStatusError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedUpdateMerchant,
		fields...,
	)
}

// HandleTrashedMerchantError processes errors that occur during the trashing of a merchant.
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
//   - A MerchantResponse containing the details of the merchant trashing failure.
//   - A standardized ErrorResponse detailing the trashing failure.
func (e *merchantCommandError) HandleTrashedMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedTrashMerchant,
		fields...,
	)
}

// HandleRestoreMerchantError processes errors that occur during the restore of a merchant.
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
//   - A MerchantResponse containing the details of the merchant restore failure.
//   - A standardized ErrorResponse detailing the restore failure.
func (e *merchantCommandError) HandleRestoreMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedRestoreMerchant,
		fields...,
	)
}

// HandleDeleteMerchantPermanentError processes errors that occur during the deletion of a merchant.
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
//   - A boolean indicating whether the deletion failed.
//   - A standardized ErrorResponse detailing the deletion failure.
func (e *merchantCommandError) HandleDeleteMerchantPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedDeleteMerchant,
		fields...,
	)
}

// HandleRestoreAllMerchantError processes errors that occur during the restore of all merchants.
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
//   - A boolean indicating whether the restore failed.
//   - A standardized ErrorResponse detailing the restore failure.
func (e *merchantCommandError) HandleRestoreAllMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedRestoreAllMerchants,
		fields...,
	)
}

// HandleDeleteAllMerchantPermanentError processes errors that occur during the permanent deletion of all merchants.
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
//   - A boolean indicating whether the deletion failed.
//   - A standardized ErrorResponse detailing the deletion failure.
func (e *merchantCommandError) HandleDeleteAllMerchantPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedDeleteAllMerchants,
		fields...,
	)
}

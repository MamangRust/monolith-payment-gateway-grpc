package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantDocumentCommandError is a struct that implements the MerchantDocumentCommandErrorHandler interface.
type merchantDocumentCommandError struct {
	logger logger.LoggerInterface
}

// NewMerchantDocumentCommandError initializes a new merchantDocumentCommandError with the provided logger.
// It returns an instance of the merchantDocumentCommandError struct.
func NewMerchantDocumentCommandError(logger logger.LoggerInterface) MerchantDocumentCommandErrorHandler {
	return &merchantDocumentCommandError{
		logger: logger,
	}
}

// HandleCreateMerchantDocumentError processes errors that occur during the creation of a merchant document.
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
//   - A MerchantDocumentResponse if successful or nil.
//   - A standardized ErrorResponse detailing the creation failure.
func (e *merchantDocumentCommandError) HandleCreateMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedCreateMerchantDocument,
		fields...,
	)
}

// HandleUpdateMerchantDocumentError processes errors that occur during the updating of a merchant document.
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
//   - A MerchantDocumentResponse if successful or nil.
//   - A standardized ErrorResponse detailing the update failure.
func (e *merchantDocumentCommandError) HandleUpdateMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedUpdateMerchantDocument,
		fields...,
	)
}

// HandleUpdateMerchantDocumentStatusError processes errors that occur during the updating of a merchant document status.
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
//   - A MerchantDocumentResponse if successful or nil.
//   - A standardized ErrorResponse detailing the update failure.
func (e *merchantDocumentCommandError) HandleUpdateMerchantDocumentStatusError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedUpdateMerchantDocument,
		fields...,
	)
}

// HandleTrashedMerchantDocumentError processes errors that occur during the trashing of a merchant document.
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
//   - A MerchantDocumentResponse if successful or nil.
//   - A standardized ErrorResponse detailing the trashing failure.
func (e *merchantDocumentCommandError) HandleTrashedMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedTrashMerchantDocument,
		fields...,
	)
}

// HandleRestoreMerchantDocumentError processes errors that occur during the restoration of a merchant document.
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
//   - A MerchantDocumentResponse if successful or nil.
//   - A standardized ErrorResponse detailing the restoration failure.
func (e *merchantDocumentCommandError) HandleRestoreMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedRestoreMerchantDocument,
		fields...,
	)
}

// HandleDeleteMerchantDocumentPermanentError processes errors that occur during the permanent deletion of a merchant document.
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
func (e *merchantDocumentCommandError) HandleDeleteMerchantDocumentPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedDeleteMerchantDocument,
		fields...,
	)
}

// HandleRestoreAllMerchantDocumentError processes errors that occur during the restoration of all merchant documents.
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
//   - A boolean indicating whether the restoration failed.
//   - A standardized ErrorResponse detailing the restoration failure.
func (e *merchantDocumentCommandError) HandleRestoreAllMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedRestoreAllMerchantDocuments,
		fields...,
	)
}

// HandleDeleteAllMerchantDocumentPermanentError processes errors that occur during the permanent deletion of all merchant documents.
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
func (e *merchantDocumentCommandError) HandleDeleteAllMerchantDocumentPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedDeleteAllMerchantDocuments,
		fields...,
	)
}

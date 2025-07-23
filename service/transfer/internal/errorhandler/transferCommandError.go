package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transferCommandError is a struct that implements the TransferCommandError interface.
type transferCommandError struct {
	logger logger.LoggerInterface
}

// NewTransferCommandError initializes a new transferCommandError with the provided logger.
// It returns an instance of the transferCommandError struct.
func NewTransferCommandError(logger logger.LoggerInterface) TransferCommandErrorHandler {
	return &transferCommandError{
		logger: logger,
	}
}

// HandleSenderInsufficientBalanceError processes errors related to insufficient balance
// of the sender card. It logs the error, records it to the trace span, and returns a
// standardized error response.
//
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - senderCardNumber: The card number associated with the insufficient balance.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A TransferResponse with error details and a standardized ErrorResponse.
func (t *transferCommandError) HandleSenderInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	senderCardNumber string,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[*response.TransferResponse](t.logger, err, method, tracePrefix, "InsufficientBalance", span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

// HandleReceiverInsufficientBalanceError processes errors related to insufficient balance
// of the receiver card. It logs the error, records it to the trace span, and returns a
// standardized error response.
//
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - receiverCardNumber: The card number associated with the insufficient balance.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A TransferResponse with error details and a standardized ErrorResponse.
func (t *transferCommandError) HandleReceiverInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	receiverCardNumber string,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[*response.TransferResponse](t.logger, err, method, tracePrefix, "InsufficientBalance", span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
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
func (t *transferCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

// HandleCreateTransferError processes errors related to creating a transfer.
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
//   - A TransferResponse with error details and a standardized ErrorResponse.
func (t *transferCommandError) HandleCreateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedCreateTransfer, fields...)
}

// HandleUpdateTransferError processes errors related to updating a transfer.
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
//   - A TransferResponse with error details and a standardized ErrorResponse.
func (t *transferCommandError) HandleUpdateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedUpdateTransfer, fields...)
}

// HandleTrashedTransferError processes errors related to trashing a transfer.
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
//   - A TransferResponse with error details and a standardized ErrorResponse.
func (t *transferCommandError) HandleTrashedTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedTrashedTransfer, fields...)
}

// HandleRestoreTransferError processes errors related to the restoration of a transfer.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the transfer restoration.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The trace span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A TransferResponse with error details and a standardized ErrorResponse.
func (t *transferCommandError) HandleRestoreTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedRestoreTransfer, fields...)
}

// HandleDeleteTransferPermanentError processes errors that occur during the permanent deletion of a transfer.
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
//   - A standardized ErrorResponse detailing the deletion failure.
func (t *transferCommandError) HandleDeleteTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedDeleteTransferPermanent, fields...)
}

// HandleRestoreAllTransferError processes errors that occur during the restoration of all transfers.
// It logs the error, records it to the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the transfer restoration.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The trace span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A boolean indicating whether the restoration was successful.
//   - A standardized ErrorResponse detailing the restoration failure.
func (t *transferCommandError) HandleRestoreAllTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedRestoreAllTransfers, fields...)
}

// HandleDeleteAllTransferPermanentError processes errors that occur during the permanent deletion of all transfers.
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
//   - A standardized ErrorResponse detailing the deletion failure.
func (t *transferCommandError) HandleDeleteAllTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedDeleteAllTransfersPermanent, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transferStatisticByCardError is a struct that implements the TransferStatisticByCardError interface.
type transferStatisticByCardError struct {
	logger logger.LoggerInterface
}

// NewTransferStatisticByCardError initializes a new transferStatisticByCardError with the provided logger.
// It returns an instance of the transferStatisticByCardError struct.
func NewTransferStatisticByCardError(logger logger.LoggerInterface) TransferStatisticByCardErrorHandler {
	return &transferStatisticByCardError{
		logger: logger,
	}
}

// HandleMonthTransferStatusSuccessByCardNumberError handles the error of retrieving monthly successful transfers by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransferResponseMonthStatusSuccess containing the details of the successful operation, and an ErrorResponse indicating failure.
func (t *transferStatisticByCardError) HandleMonthTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseMonthStatusSuccess](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthTransferStatusSuccessByCard,
		fields...,
	)
}

// HandleYearTransferStatusSuccessByCardNumberError handles the error of retrieving yearly successful transfers by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransferResponseYearStatusSuccess containing the details of the successful operation, and an ErrorResponse indicating failure.
func (t *transferStatisticByCardError) HandleYearTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseYearStatusSuccess](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusSuccessByCard,
		fields...,
	)
}

// HandleMonthTransferStatusFailedByCardNumberError handles the error of retrieving monthly failed transfers by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransferResponseMonthStatusFailed containing the details of the successful operation, and an ErrorResponse indicating failure.
func (t *transferStatisticByCardError) HandleMonthTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseMonthStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthTransferStatusFailedByCard,
		fields...,
	)
}

// HandleYearTransferStatusFailedByCardNumberError handles the error of retrieving yearly failed transfers by card number.
// It logs the error, records it to the trace span, and returns a structured ErrorResponse indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the error is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TransferResponseYearStatusFailed containing the details of the successful operation, and an ErrorResponse indicating failure.
func (t *transferStatisticByCardError) HandleYearTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseYearStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusFailedByCard,
		fields...,
	)
}

// HandleMonthlyTransferAmountsBySenderError processes errors during retrieval of monthly transfer amounts by sender's card.
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
//   - A slice of TransferMonthAmountResponse with error details and a standardized ErrorResponse.
func (t *transferStatisticByCardError) HandleMonthlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferMonthAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthlyTransferAmountsBySenderCard,
		fields...,
	)
}

// HandleYearlyTransferAmountsBySenderError processes errors during retrieval of yearly transfer amounts by sender's card.
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
//   - A slice of TransferYearAmountResponse with error details and a standardized ErrorResponse.
func (t *transferStatisticByCardError) HandleYearlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferYearAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearlyTransferAmountsBySenderCard,
		fields...,
	)
}

// HandleMonthlyTransferAmountsByReceiverError processes errors during retrieval of monthly transfer amounts by receiver's card.
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
//   - A slice of TransferMonthAmountResponse with error details and a standardized ErrorResponse.
func (t *transferStatisticByCardError) HandleMonthlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferMonthAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthlyTransferAmountsByReceiverCard,
		fields...,
	)
}

// HandleYearlyTransferAmountsByReceiverError processes errors during retrieval of yearly transfer amounts by receiver's card.
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
//   - A slice of TransferYearAmountResponse with error details and a standardized ErrorResponse.
func (t *transferStatisticByCardError) HandleYearlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferYearAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearlyTransferAmountsByReceiverCard,
		fields...,
	)
}

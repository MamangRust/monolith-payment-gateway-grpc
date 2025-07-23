package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transferStatisticError is a struct that implements the TransferStatisticError interface.
type transferStatisticError struct {
	logger logger.LoggerInterface
}

// NewTransferStatisticError initializes a new transferStatisticError with the provided logger.
// It returns an instance of the transferStatisticError struct.
func NewTransferStatisticError(logger logger.LoggerInterface) TransferStatisticErrorHandler {
	return &transferStatisticError{
		logger: logger,
	}
}

// HandleMonthTransferStatusSuccessError handles errors related to retrieving monthly successful transfer statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transfer status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransferResponseMonthStatusSuccess pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transferStatisticError) HandleMonthTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseMonthStatusSuccess](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedFindMonthTransferStatusSuccess, fields...)
}

// HandleYearTransferStatusSuccessError handles errors related to retrieving yearly successful transfer statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transfer status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransferResponseYearStatusSuccess pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transferStatisticError) HandleYearTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseYearStatusSuccess](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusSuccess,
		fields...,
	)
}

// HandleMonthTransferStatusFailedError handles errors related to retrieving monthly failed transfer statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transfer status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransferResponseMonthStatusFailed pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transferStatisticError) HandleMonthTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseMonthStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthTransferStatusFailed,
		fields...,
	)
}

// HandleYearTransferStatusFailedError handles errors related to retrieving yearly failed transfer statuses.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error encountered while processing the transfer status.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The tracing span to record error details.
//   - status: A pointer to a string representing the status, which will be updated as necessary.
//   - fields: Additional contextual fields for logging purposes.
//
// Returns:
//   - A slice of TransferResponseYearStatusFailed pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the failure, if any.
func (t *transferStatisticError) HandleYearTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponseYearStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusFailed,
		fields...,
	)
}

// HandleMonthlyTransferAmountsError processes errors during retrieval of monthly transfer amounts.
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
//   - A slice of TransferMonthAmountResponse with error details and a standardized ErrorResponse.
func (t *transferStatisticError) HandleMonthlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferMonthAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthlyTransferAmounts,
		fields...,
	)
}

// HandleYearlyTransferAmountsError processes errors during retrieval of yearly transfer amounts.
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
//   - A slice of TransferYearAmountResponse with error details and a standardized ErrorResponse.
func (t *transferStatisticError) HandleYearlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferYearAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearlyTransferAmounts,
		fields...,
	)
}

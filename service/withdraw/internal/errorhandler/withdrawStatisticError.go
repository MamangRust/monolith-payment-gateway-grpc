package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// withdrawStatisticError handles error logging for withdraw statistic operations.
type withdrawStatisticError struct {
	logger logger.LoggerInterface
}

// NewWithdrawStatisticError initializes a new withdrawStatisticError with the provided logger.
// It returns an instance of the withdrawStatisticError struct.
func NewWithdrawStatisticError(logger logger.LoggerInterface) WithdrawStatisticErrorHandler {
	return &withdrawStatisticError{
		logger: logger,
	}
}

// HandleMonthWithdrawStatusSuccessError processes errors during the retrieval of monthly successful withdraw status.
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
//   - A slice of WithdrawResponseMonthStatusSuccess with error details and a standardized ErrorResponse.
func (w *withdrawStatisticError) HandleMonthWithdrawStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseMonthStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess, fields...)
}

// HandleYearWithdrawStatusSuccessError processes errors during the retrieval of yearly successful withdraw status.
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
//   - A slice of WithdrawResponseYearStatusSuccess with error details and a standardized ErrorResponse.
func (w *withdrawStatisticError) HandleYearWithdrawStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseYearStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess, fields...)
}

// HandleMonthWithdrawStatusFailedError processes errors during the retrieval of monthly failed withdraw status.
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
//   - A slice of WithdrawResponseMonthStatusFailed with error details and a standardized ErrorResponse.
func (w *withdrawStatisticError) HandleMonthWithdrawStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseMonthStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed, fields...)
}

// HandleYearWithdrawStatusFailedError processes errors during the retrieval of yearly failed withdraw status.
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
//   - A slice of WithdrawResponseYearStatusFailed with error details and a standardized ErrorResponse.
func (w *withdrawStatisticError) HandleYearWithdrawStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseYearStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusFailed, fields...)
}

// HandleMonthlyWithdrawAmountsError processes errors during the retrieval of monthly withdraw amounts.
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
//   - A slice of WithdrawMonthlyAmountResponse with error details and a standardized ErrorResponse.
func (w *withdrawStatisticError) HandleMonthlyWithdrawAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawMonthlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthlyWithdraws, fields...)
}

// HandleYearlyWithdrawAmountsError processes errors during the retrieval of yearly withdraw amounts.
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
//   - A slice of WithdrawYearlyAmountResponse with error details and a standardized ErrorResponse.
func (w *withdrawStatisticError) HandleYearlyWithdrawAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawYearlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearlyWithdraws, fields...)
}

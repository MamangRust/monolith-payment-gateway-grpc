package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// withdrawStatisticByCardError handles error logging for withdraw statistic by card operations.
type withdrawStatisticByCardError struct {
	logger logger.LoggerInterface
}

// NewWithdrawStatisticByCardError initializes a new withdrawStatisticByCardError with the provided logger.
// It returns an instance of the withdrawStatisticByCardError struct.
func NewWithdrawStatisticByCardError(logger logger.LoggerInterface) WithdrawStatisticByCardErrorHandler {
	return &withdrawStatisticByCardError{
		logger: logger,
	}
}

// HandleMonthWithdrawStatusSuccessByCardNumberError processes errors during retrieval of monthly successful withdraw status by card number.
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
func (w *withdrawStatisticByCardError) HandleMonthWithdrawStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseMonthStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccessByCard, fields...)
}

// HandleYearWithdrawStatusSuccessByCardNumberError processes errors during retrieval of yearly successful withdraw status by card number.
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
func (w *withdrawStatisticByCardError) HandleYearWithdrawStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseYearStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccessByCard, fields...)
}

// HandleMonthWithdrawStatusFailedByCardNumberError processes errors during retrieval of monthly failed withdraw status by card number.
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
func (w *withdrawStatisticByCardError) HandleMonthWithdrawStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseMonthStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailedByCard, fields...)
}

// HandleYearWithdrawStatusFailedByCardNumberError processes errors during retrieval of yearly failed withdraw status by card number.
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
func (w *withdrawStatisticByCardError) HandleYearWithdrawStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseYearStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusFailedByCard, fields...)
}

// HandleMonthlyWithdrawsAmountByCardNumberError processes errors during retrieval of monthly withdraw amounts by card number.
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
func (w *withdrawStatisticByCardError) HandleMonthlyWithdrawsAmountByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawMonthlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthlyWithdrawsByCardNumber, fields...)
}

// HandleYearlyWithdrawsAmountByCardNumberError processes errors during retrieval of yearly withdraw amounts by card number.
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
func (w *withdrawStatisticByCardError) HandleYearlyWithdrawsAmountByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawYearlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearlyWithdrawsByCardNumber, fields...)
}

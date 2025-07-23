package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// topupStatisticByCardError is a struct that implements the TopupStatisticByCardError interface.
type topupStatisticByCardError struct {
	logger logger.LoggerInterface
}

// NewTopupStatisticByCardError returns a new instance of TopupStatisticByCardError with the given logger.
func NewTopupStatisticByCardError(logger logger.LoggerInterface) TopupStatisticByCardErrorHandler {
	return &topupStatisticByCardError{
		logger: logger,
	}
}

// HandleMonthTopupStatusSuccessByCardNumber handles the successful retrieval of monthly topup status by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating success.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseMonthStatusSuccess containing the details of the successful operation, and a nil ErrorResponse indicating success.
func (e *topupStatisticByCardError) HandleMonthTopupStatusSuccessByCardNumber(err error,
	method,
	tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseMonthStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusSuccessByCard, fields...)
}

// HandleYearlyTopupStatusSuccessByCardNumber handles the successful retrieval of yearly topup status by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating success.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseYearStatusSuccess containing the details of the successful operation, and a nil ErrorResponse indicating success.
func (e *topupStatisticByCardError) HandleYearlyTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseYearStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusSuccessByCard, fields...)
}

// HandleMonthTopupStatusFailedByCardNumber handles the failure to retrieve monthly topup status by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseMonthStatusFailed containing the details of the failed operation, and an ErrorResponse containing more information about the failure.
func (e *topupStatisticByCardError) HandleMonthTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseMonthStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusFailedByCard, fields...)
}

// HandleYearlyTopupStatusFailedByCardNumber handles the failure to retrieve yearly topup status by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseYearStatusFailed containing the details of the failed operation, and an ErrorResponse containing more information about the failure.
func (e *topupStatisticByCardError) HandleYearlyTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseYearStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusFailedByCard, fields...)
}

// HandleMonthlyTopupMethodsByCardNumber handles the successful retrieval of monthly topup methods by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating success.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupMonthMethodResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
func (e *topupStatisticByCardError) HandleMonthlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupMonthMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupMethodsByCard, fields...)
}

// HandleYearlyTopupMethodsByCardNumber handles the successful retrieval of yearly topup methods by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating success.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupYearlyMethodResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
func (e *topupStatisticByCardError) HandleYearlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupYearlyMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupMethodsByCard, fields...)
}

// HandleMonthlyTopupAmountsByCardNumber handles the successful retrieval of monthly topup amounts by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating success.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupMonthAmountResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
func (e *topupStatisticByCardError) HandleMonthlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupMonthAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupAmountsByCard, fields...)
}

// HandleYearlyTopupAmountsByCardNumber handles the successful retrieval of yearly topup amounts by card number.
// It logs the information, records it to the trace span, and returns a structured response indicating success.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupYearlyAmountResponse containing the details of the successful operation, and a nil ErrorResponse indicating success.
func (e *topupStatisticByCardError) HandleYearlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupYearlyAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupAmountsByCard, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// topupStatisticError is a struct that implements the TopupStatisticError interface.
type topupStatisticError struct {
	logger logger.LoggerInterface
}

// NewTopupStatisticError initializes a new instance of topupStatisticError with the provided logger.
// This function returns a pointer to the topupStatisticError struct, which implements the TopupStatisticError interface.
// It is used for handling errors related to top-up statistics, ensuring that they are logged appropriately.
func NewTopupStatisticError(logger logger.LoggerInterface) TopupStatisticErrorHandler {
	return &topupStatisticError{
		logger: logger,
	}
}

// HandleMonthTopupStatusSuccess processes the successful retrieval of monthly topup status.
// It logs the success information, records it to the trace span, and returns a structured response.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseMonthStatusSuccess containing the details of the successful operation,
//     and a nil ErrorResponse indicating success.
func (e *topupStatisticError) HandleMonthTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseMonthStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusSuccess, fields...)
}

// HandleYearlyTopupStatusSuccess processes the successful retrieval of yearly topup status.
// It logs the success information, records it to the trace span, and returns a structured response.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the success is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the success.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseYearStatusSuccess containing the details of the successful operation,
//     and a nil ErrorResponse indicating success.
func (e *topupStatisticError) HandleYearlyTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseYearStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusSuccess, fields...)
}

// HandleMonthTopupStatusFailed processes the failure to retrieve monthly topup status.
// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - err: The error, if any, encountered during the process.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseMonthStatusFailed containing the details of the failed operation,
//     and an ErrorResponse containing more information about the failure.
func (e *topupStatisticError) HandleMonthTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseMonthStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusFailed, fields...)
}

// HandleYearlyTopupStatusFailed processes the failure to retrieve yearly topup status.
// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupResponseYearStatusFailed containing the details of the failed operation,
//     and an ErrorResponse containing more information about the failure.
func (e *topupStatisticError) HandleYearlyTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseYearStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusFailed, fields...)
}

// HandleMonthlyTopupMethods processes the failure to retrieve monthly topup methods.
// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupMonthMethodResponse containing the details of the failed operation,
//     and an ErrorResponse containing more information about the failure.
func (e *topupStatisticError) HandleMonthlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupMonthMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupMethods, fields...)
}

// HandleYearlyTopupMethods processes the failure to retrieve yearly topup methods.
// It logs the failure information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupYearlyMethodResponse containing the details of the failed operation,
//     and an ErrorResponse containing more information about the failure.
func (e *topupStatisticError) HandleYearlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupYearlyMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupMethods, fields...)
}

// HandleMonthlyTopupAmounts processes the retrieval of monthly topup amounts.
// It logs the error information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupMonthAmountResponse containing the details of the failed operation,
//     and an ErrorResponse containing more information about the failure.
func (e *topupStatisticError) HandleMonthlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupMonthAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupAmounts, fields...)
}

// HandleYearlyTopupAmounts processes the retrieval of yearly topup amounts.
// It logs the error information, records it to the trace span, and returns a structured response indicating failure.
// Parameters:
//   - err: The error, if any, encountered during the process.
//   - method: The name of the method where the failure is recorded.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the failure.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of TopupYearlyAmountResponse containing the details of the failed operation,
//     and an ErrorResponse containing more information about the failure.
func (e *topupStatisticError) HandleYearlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupYearlyAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupAmounts, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// cardStatisticError is a struct that implements the CardStatisticError interface
type cardStatisticError struct {
	logger logger.LoggerInterface
}

// NewCardStatisticError initializes a new cardStatisticError with the provided logger.
// It returns an instance of the cardStatisticError struct.
func NewCardStatisticError(logger logger.LoggerInterface) CardStatisticErrorHandler {
	return &cardStatisticError{
		logger: logger,
	}
}

// HandleMonthlyBalanceError processes errors during the retrieval of a monthly balance.
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
//   - A CardResponse with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleMonthlyBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyBalance, fields...)
}

// HandleYearlyBalanceError processes errors during the retrieval of yearly balance.
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
//   - A CardResponseYearlyBalance with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleYearlyBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearlyBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyBalance, fields...)
}

// HandleMonthlyTopupAmountError processes errors during the retrieval of monthly topup amount.
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
//   - A slice of CardResponseMonthAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleMonthlyTopupAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTopupAmount, fields...)
}

// HandleYearlyTopupAmountError processes errors during the retrieval of yearly topup amount.
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
//   - A slice of CardResponseYearAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleYearlyTopupAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTopupAmount, fields...)
}

// HandleMonthlyWithdrawAmountError processes errors during the retrieval of monthly withdraw amount.
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
//   - A slice of CardResponseMonthAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleMonthlyWithdrawAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyWithdrawAmount, fields...)
}

// HandleYearlyWithdrawAmountError processes errors during the retrieval of yearly withdraw amounts.
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
//   - A slice of CardResponseYearAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleYearlyWithdrawAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyWithdrawAmount, fields...)
}

// HandleMonthlyTransactionAmountError processes errors during retrieval of monthly transaction amount.
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
//   - A slice of CardResponseMonthAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleMonthlyTransactionAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransactionAmount, fields...)
}

// HandleYearlyTransactionAmountError processes errors during retrieval of yearly transaction amount.
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
//   - A slice of CardResponseYearAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleYearlyTransactionAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransactionAmount, fields...)
}

// HandleMonthlyTransferAmountSenderError processes errors during retrieval of monthly transfer amounts by sender.
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
//   - A slice of CardResponseMonthAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleMonthlyTransferAmountSenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountSender, fields...)
}

// HandleYearlyTransferAmountSenderError processes errors during retrieval of yearly transfer amounts by sender.
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
//   - A slice of CardResponseYearAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleYearlyTransferAmountSenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountSender, fields...)
}

// HandleMonthlyTransferAmountReceiverError processes errors during retrieval of monthly transfer amounts by receiver.
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
//   - A slice of CardResponseMonthAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleMonthlyTransferAmountReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountReceiver, fields...)
}

// HandleYearlyTransferAmountReceiverError processes errors during the retrieval of yearly transfer amounts by receiver.
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
//   - A slice of CardResponseYearAmount with error details and a standardized ErrorResponse.
func (c *cardStatisticError) HandleYearlyTransferAmountReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountReceiver, fields...)
}

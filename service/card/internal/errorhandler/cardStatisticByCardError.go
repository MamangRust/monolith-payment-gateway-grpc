package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// cardStatisticByNumberError is a struct that implements the CardStatisticByNumberError interface
type cardStatisticByNumberError struct {
	logger logger.LoggerInterface
}

// NewCardStatisticByNumberError initializes a new cardStatisticByNumberError with the provided logger.
// It returns an instance of the cardStatisticByNumberError struct.
func NewCardStatisticByNumberError(logger logger.LoggerInterface) CardStatisticByNumberErrorHandler {
	return &cardStatisticByNumberError{
		logger: logger,
	}
}

// HandleMonthlyBalanceByCardNumberError processes errors during the retrieval of a monthly balance by card number.
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
func (c *cardStatisticByNumberError) HandleMonthlyBalanceByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyBalanceByCard, fields...)
}

// HandleYearlyBalanceByCardNumberError processes errors during retrieval of yearly balance by card number.
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
//   - A slice of CardResponseYearlyBalance with error details and a standardized ErrorResponse.
func (c *cardStatisticByNumberError) HandleYearlyBalanceByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearlyBalance](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyBalanceByCard, fields...)
}

// HandleMonthlyTopupAmountByCardNumberError processes errors during retrieval of monthly topup amount by card number.
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
func (c *cardStatisticByNumberError) HandleMonthlyTopupAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTopupAmountByCard, fields...)
}

// HandleYearlyTopupAmountByCardNumberError processes errors during retrieval of yearly topup amount by card number.
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

func (c *cardStatisticByNumberError) HandleYearlyTopupAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTopupAmountByCard, fields...)
}

// HandleMonthlyWithdrawAmountByCardNumberError processes errors during retrieval of monthly withdraw amount by card number.
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
func (c *cardStatisticByNumberError) HandleMonthlyWithdrawAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyWithdrawAmountByCard, fields...)
}

// HandleYearlyWithdrawAmountByCardNumberError processes errors during retrieval of yearly withdraw amount by card number.
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
func (c *cardStatisticByNumberError) HandleYearlyWithdrawAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyWithdrawAmountByCard, fields...)
}

// HandleMonthlyTransactionAmountByCardNumberError processes errors during retrieval of monthly transaction amount by card number.
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
func (c *cardStatisticByNumberError) HandleMonthlyTransactionAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransactionAmountByCard, fields...)
}

// HandleYearlyTransactionAmountByCardNumberError processes errors during retrieval of yearly transaction amount by card number.
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
func (c *cardStatisticByNumberError) HandleYearlyTransactionAmountByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransactionAmountByCard, fields...)
}

// HandleMonthlyTransferAmountBySenderError processes errors during retrieval of monthly transfer amounts by sender.
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
func (c *cardStatisticByNumberError) HandleMonthlyTransferAmountBySenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountBySender, fields...)
}

// HandleYearlyTransferAmountBySenderError processes errors during retrieval of yearly transfer amounts by sender.
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
func (c *cardStatisticByNumberError) HandleYearlyTransferAmountBySenderError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountBySender, fields...)
}

// HandleMonthlyTransferAmountByReceiverError processes errors during retrieval of monthly transfer amounts by receiver.
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
func (c *cardStatisticByNumberError) HandleMonthlyTransferAmountByReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseMonthAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindMonthlyTransferAmountByReceiver, fields...)
}

// HandleYearlyTransferAmountByReceiverError processes errors during retrieval of yearly transfer amounts by receiver.
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
func (c *cardStatisticByNumberError) HandleYearlyTransferAmountByReceiverError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CardResponseYearAmount](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindYearlyTransferAmountByReceiver, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// cardDasboardError is a struct that implements the CardDashboardError interface
type cardDasboardError struct {
	logger logger.LoggerInterface
}

// NewCardDashboardError initializes a new cardDasboardError with the provided logger.
// It returns a pointer to the cardDasboardError struct.
func NewCardDashboardError(logger logger.LoggerInterface) CardDashboardErrorHandler {
	return &cardDasboardError{
		logger: logger,
	}
}

// HandleTotalBalanceError processes errors during retrieval of total balance.
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
//   - A DashboardCard with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalBalances, fields...)
}

// HandleTotalTopupAmountError processes errors during retrieval of total topup amount.
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
//   - A DashboardCard with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTopupAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTopAmount, fields...)
}

// HandleTotalWithdrawAmountError processes errors during retrieval of total withdraw amount.
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
//   - A DashboardCard with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalWithdrawAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalWithdrawAmount, fields...)
}

// HandleTotalTransactionAmountError processes errors during retrieval of total transaction amount.
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
//   - A DashboardCard with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTransactionAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransactionAmount, fields...)
}

// HandleTotalTransferAmountError processes errors during retrieval of total transfer amount.
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
//   - A DashboardCard with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTransferAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransferAmount, fields...)
}

// HandleTotalBalanceCardNumberError processes errors during retrieval of total balance by card number.
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
//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalBalanceCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalBalanceByCard, fields...)
}

// HandleTotalTopupAmountCardNumberError processes errors during retrieval of total topup amount by card number.
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
//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTopupAmountCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTopupAmountByCard, fields...)
}

// HandleTotalWithdrawAmountCardNumberError processes errors during retrieval of total withdraw amount by card number.
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
//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalWithdrawAmountCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalWithdrawAmountByCard, fields...)
}

// HandleTotalTransactionAmountCardNumberError processes errors during retrieval of total transaction amount by card number.
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
//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTransactionAmountCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransactionAmountByCard, fields...)
}

// HandleTotalTransferAmountBySender processes errors during retrieval of total transfer amount by sender.
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
//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTransferAmountBySender(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransferAmountBySender, fields...)
}

// HandleTotalTransferAmountByReceiver processes errors during retrieval of total transfer amount by receiver.
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
//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
func (c *cardDasboardError) HandleTotalTransferAmountByReceiver(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransferAmountByReceiver, fields...)
}

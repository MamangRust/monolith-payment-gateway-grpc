package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// CardQueryErrorHandler is a struct that implements the CardQueryError interface
//
//go:generate mockgen -source=interfaces.go -destination=mocks/errorhandler.go
type CardQueryErrorHandler interface {
	// HandleFindAllError processes errors when fetching all cards.
	// It logs the error, records it to the trace span, and returns a paginated CardResponse with error details.
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of CardResponse pointers with paginated card details.
	//   - A pointer to an integer indicating the total count of cards.
	//   - A standardized ErrorResponse detailing the error.
	HandleFindAllError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponse, *int, *response.ErrorResponse)

	// HandleFindByActiveError processes errors when fetching all active cards.
	// It logs the error, records it to the trace span, and returns a paginated CardResponse with error details.
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of CardResponse pointers with paginated card details.
	//   - A pointer to an integer indicating the total count of cards.
	//   - A standardized ErrorResponse detailing the error.
	HandleFindByActiveError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	// HandleFindByTrashedError processes errors when fetching all trashed cards.
	// It logs the error, records it to the trace span, and returns a paginated CardResponse with error details.
	// Parameters:
	//   - err: The error that occurred.
	//   - method: The name of the method where the error occurred.
	//   - tracePrefix: A prefix for generating the trace ID.
	//   - span: The trace span used for recording the error.
	//   - status: A pointer to a string that will be set with the formatted status.
	//   - fields: Additional fields to include in the log entry.
	//
	// Returns:
	//   - A slice of CardResponse pointers with paginated card details.
	//   - A pointer to an integer indicating the total count of cards.
	//   - A standardized ErrorResponse detailing the error.
	HandleFindByTrashedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	// HandleFindByIdError processes errors during card lookup by ID
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
	HandleFindByIdError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	// HandleFindByUserIdError processes errors during card lookup by user ID
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
	HandleFindByUserIdError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	// HandleFindByCardNumberError processes errors during card lookup by card number
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
	HandleFindByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
}

type CardDashboardErrorHandler interface {
	// HandleTotalBalanceError processes errors during retrieval of total balance.
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
	//   - A DashboardCard with error details and a standardized ErrorResponse.
	HandleTotalBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	// HandleTotalTopupAmountError processes errors during retrieval of total topup amount.
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
	//   - A DashboardCard with error details and a standardized ErrorResponse.
	HandleTotalTopupAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	// HandleTotalWithdrawAmountError processes errors during retrieval of total withdraw amount.
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
	//   - A DashboardCard with error details and a standardized ErrorResponse.
	HandleTotalWithdrawAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	// HandleTotalTransactionAmountError processes errors during retrieval of total transaction amount.
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
	//   - A DashboardCard with error details and a standardized ErrorResponse.
	HandleTotalTransactionAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	// HandleTotalTransferAmountError processes errors during retrieval of total transfer amount.
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
	//   - A DashboardCard with error details and a standardized ErrorResponse.
	HandleTotalTransferAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	// HandleTotalBalanceCardNumberError processes errors during retrieval of total balance by card number.
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
	//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
	HandleTotalBalanceCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	// HandleTotalTopupAmountCardNumberError processes errors during retrieval of total topup amount by card number.
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
	//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
	HandleTotalTopupAmountCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	// HandleTotalWithdrawAmountCardNumberError processes errors during retrieval of total withdraw amount by card number.
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
	//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
	HandleTotalWithdrawAmountCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	// HandleTotalTransactionAmountCardNumberError processes errors during retrieval of total transaction amount by card number.
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
	//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
	HandleTotalTransactionAmountCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	// HandleTotalTransferAmountBySender processes errors during retrieval of total transfer amount by sender.
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
	//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
	HandleTotalTransferAmountBySender(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	// HandleTotalTransferAmountByReceiver processes errors during retrieval of total transfer amount by receiver.
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
	//   - A DashboardCardCardNumber with error details and a standardized ErrorResponse.
	HandleTotalTransferAmountByReceiver(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
}

type CardStatisticErrorHandler interface {
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
	HandleMonthlyBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)

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
	HandleYearlyBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)

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
	HandleMonthlyTopupAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

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
	HandleYearlyTopupAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

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
	HandleMonthlyWithdrawAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

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
	HandleYearlyWithdrawAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyTransactionAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

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
	HandleYearlyTransactionAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyTransferAmountSenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	HandleYearlyTransferAmountSenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyTransferAmountReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	HandleYearlyTransferAmountReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardStatisticByNumberErrorHandler interface {
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
	HandleMonthlyBalanceByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)
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
	HandleYearlyBalanceByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)
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
	HandleMonthlyTopupAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	HandleYearlyTopupAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyWithdrawAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	HandleYearlyWithdrawAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyTransactionAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	HandleYearlyTransactionAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyTransferAmountBySenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	HandleYearlyTransferAmountBySenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
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
	HandleMonthlyTransferAmountByReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
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
	//   - A slice of Card
	HandleYearlyTransferAmountByReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardCommandErrorHandler interface {
	// HandleFindByIdUserError processes errors during user lookup by ID
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
	HandleFindByIdUserError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	// HandleCreateCardError processes errors that occur during card creation.
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
	//   - A CardResponse, which is nil since operation failed.
	//   - A standardized ErrorResponse detailing the card creation failure.
	HandleCreateCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	// HandleUpdateCardError processes errors that occur during card updates.
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
	//   - A CardResponse, which is nil since operation failed.
	//   - A standardized ErrorResponse detailing the card update failure.
	HandleUpdateCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	// HandleTrashedCardError processes errors that occur during card trashing.
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
	//   - A CardResponse, which is nil since operation failed.
	//   - A standardized ErrorResponse detailing the card trashing failure.
	HandleTrashedCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponseDeleteAt, *response.ErrorResponse)
	// HandleRestoreCardError processes errors that occur during card restoration.
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
	//   - A CardResponse, which is nil since operation failed.
	//   - A standardized ErrorResponse detailing the card restoration failure.
	HandleRestoreCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	// HandleDeleteCardPermanentError processes errors that occur during card deletion.
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
	//   - A boolean indicating whether the error is fatal.
	//   - A standardized ErrorResponse detailing the card deletion failure.
	HandleDeleteCardPermanentError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	// HandleRestoreAllCardError processes errors that occur during card restoration.
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
	//   - A boolean indicating whether the error is fatal.
	//   - A standardized ErrorResponse detailing the card restoration failure.
	HandleRestoreAllCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	// HandleDeleteAllCardPermanentError processes errors that occur during the permanent deletion of all cards.
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
	//   - A boolean indicating whether the error is fatal.
	//   - A standardized ErrorResponse detailing the card deletion failure.

	HandleDeleteAllCardPermanentError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

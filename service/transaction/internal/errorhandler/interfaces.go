package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TransactionQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponse, *int, *response.ErrorResponse)

	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)

	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)

	HanldeRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransactionResponse, *response.ErrorResponse)
}

type TransactionCommandErrorHandler interface {
	HandleInvalidParseTimeError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		rawTime string,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
	HandleInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		cardNumber string,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
	HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)

	HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)

	HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)

	HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)

	HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)

	HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)

	HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)

	HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

type TransactionStatisticErrorHandler interface {
	HandleMonthTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	HandleYearlyTransactionStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	HandleMonthTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	HandleYearlyTransactionStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	HandleMonthlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	HandleYearlyPaymentMethodsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)

	HandleMonthlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	HandleYearlyAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

type TransactionStatisticByCardErrorHandler interface {
	HandleMonthTransactionStatusSuccessByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)

	HandleYearlyTransactionStatusSuccessByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,

		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)

	HandleMonthTransactionStatusFailedByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)

	HandleYearlyTransactionStatusFailedByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	HandleMonthlyPaymentMethodsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)

	HandleYearlyPaymentMethodsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)

	HandleMonthlyAmountsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)

	HandleYearlyAmountsByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

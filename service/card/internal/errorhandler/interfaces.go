package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CardQueryErrorHandler interface {
	HandleFindAllError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponse, *int, *response.ErrorResponse)
	HandleFindByActiveError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	HandleFindByTrashedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	HandleFindByIdError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	HandleFindByUserIdError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	HandleFindByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
}

type CardDashboardErrorHandler interface {
	HandleTotalBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	HandleTotalTopupAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	HandleTotalWithdrawAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	HandleTotalTransactionAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	HandleTotalTransferAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCard, *response.ErrorResponse)
	HandleTotalBalanceCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	HandleTotalTopupAmountCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	HandleTotalWithdrawAmountCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	HandleTotalTransactionAmountCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	HandleTotalTransferAmountBySender(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
	HandleTotalTransferAmountByReceiver(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.DashboardCardCardNumber, *response.ErrorResponse)
}

type CardStatisticErrorHandler interface {
	HandleMonthlyBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)

	HandleYearlyBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)

	HandleMonthlyTopupAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTopupAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyWithdrawAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyWithdrawAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyTransactionAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTransactionAmountError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyTransferAmountSenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTransferAmountSenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyTransferAmountReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTransferAmountReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardStatisticByNumberErrorHandler interface {
	HandleMonthlyBalanceByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)

	HandleYearlyBalanceByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)

	HandleMonthlyTopupAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTopupAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyWithdrawAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyWithdrawAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyTransactionAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTransactionAmountByCardNumberError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyTransferAmountBySenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTransferAmountBySenderError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	HandleMonthlyTransferAmountByReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	HandleYearlyTransferAmountByReceiverError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardCommandErrorHandler interface {
	HandleFindByIdUserError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)
	HandleCreateCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)

	HandleUpdateCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)

	HandleTrashedCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)

	HandleRestoreCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (*response.CardResponse, *response.ErrorResponse)

	HandleDeleteCardPermanentError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleRestoreAllCardError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleDeleteAllCardPermanentError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

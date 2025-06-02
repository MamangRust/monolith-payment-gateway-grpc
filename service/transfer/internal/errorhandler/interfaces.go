package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TransferQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransferResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	HanldeRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransferResponse, *response.ErrorResponse)
}

type TransferCommandErrorHandler interface {
	HandleSenderInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		senderCardNumber string,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	HandleReceiverInsufficientBalanceError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		receiverCardNumber string,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransferResponse, *response.ErrorResponse)
	HandleCreateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)
	HandleUpdateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)
	HandleTrashedTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)
	HandleRestoreTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse)
	HandleDeleteTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

type TransferStatisticErrorHandler interface {
	HandleMonthTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	HandleYearTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	HandleMonthTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	HandleYearTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	HandleMonthlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	HandleYearlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

type TransferStatisticByCardErrorHandler interface {
	HandleMonthTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	HandleYearTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	HandleMonthTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	HandleYearTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	HandleMonthlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	HandleYearlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)

	HandleMonthlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	HandleYearlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferStatisticByCardError struct {
	logger logger.LoggerInterface
}

func NewTransferStatisticByCardError(logger logger.LoggerInterface) *transferStatisticByCardError {
	return &transferStatisticByCardError{
		logger: logger,
	}
}

func (t *transferStatisticByCardError) HandleMonthTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseMonthStatusSuccess](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthTransferStatusSuccessByCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleYearTransferStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseYearStatusSuccess](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusSuccessByCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleMonthTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseMonthStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthTransferStatusFailedByCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleYearTransferStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseYearStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusFailedByCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleMonthlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferMonthAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthlyTransferAmountsBySenderCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleYearlyTransferAmountsBySenderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferYearAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearlyTransferAmountsBySenderCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleMonthlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferMonthAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthlyTransferAmountsByReceiverCard,
		fields...,
	)
}

func (t *transferStatisticByCardError) HandleYearlyTransferAmountsByReceiverError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferYearAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearlyTransferAmountsByReceiverCard,
		fields...,
	)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferStatisticError struct {
	logger logger.LoggerInterface
}

func NewTransferStatisticError(logger logger.LoggerInterface) *transferStatisticError {
	return &transferStatisticError{
		logger: logger,
	}
}

func (t *transferStatisticError) HandleMonthTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseMonthStatusSuccess](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedFindMonthTransferStatusSuccess, fields...)
}

func (t *transferStatisticError) HandleYearTransferStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseYearStatusSuccess](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusSuccess,
		fields...,
	)
}

func (t *transferStatisticError) HandleMonthTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseMonthStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthTransferStatusFailed,
		fields...,
	)
}

func (t *transferStatisticError) HandleYearTransferStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferResponseYearStatusFailed](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearTransferStatusFailed,
		fields...,
	)
}

func (t *transferStatisticError) HandleMonthlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferMonthAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindMonthlyTransferAmounts,
		fields...,
	)
}

func (t *transferStatisticError) HandleYearlyTransferAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TransferYearAmountResponse](
		t.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		transfer_errors.ErrFailedFindYearlyTransferAmounts,
		fields...,
	)
}

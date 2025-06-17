package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawStatisticByCardError struct {
	logger logger.LoggerInterface
}

func NewWithdrawStatisticByCardError(logger logger.LoggerInterface) *withdrawStatisticByCardError {
	return &withdrawStatisticByCardError{
		logger: logger,
	}
}

func (w *withdrawStatisticByCardError) HandleMonthWithdrawStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseMonthStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccessByCard, fields...)
}

func (w *withdrawStatisticByCardError) HandleYearWithdrawStatusSuccessByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseYearStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccessByCard, fields...)
}

func (w *withdrawStatisticByCardError) HandleMonthWithdrawStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseMonthStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailedByCard, fields...)
}

func (w *withdrawStatisticByCardError) HandleYearWithdrawStatusFailedByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawResponseYearStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusFailedByCard, fields...)
}

func (w *withdrawStatisticByCardError) HandleMonthlyWithdrawsAmountByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawMonthlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthlyWithdrawsByCardNumber, fields...)
}

func (w *withdrawStatisticByCardError) HandleYearlyWithdrawsAmountByCardNumberError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.WithdrawYearlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearlyWithdrawsByCardNumber, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawStatisticError struct {
	logger logger.LoggerInterface
}

func NewWithdrawStatisticError(logger logger.LoggerInterface) *withdrawStatisticError {
	return &withdrawStatisticError{
		logger: logger,
	}
}

func (w *withdrawStatisticError) HandleMonthWithdrawStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.WithdrawResponseMonthStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess, fields...)
}

func (w *withdrawStatisticError) HandleYearWithdrawStatusSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.WithdrawResponseYearStatusSuccess](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess, fields...)
}

func (w *withdrawStatisticError) HandleMonthWithdrawStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.WithdrawResponseMonthStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed, fields...)
}

func (w *withdrawStatisticError) HandleYearWithdrawStatusFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.WithdrawResponseYearStatusFailed](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearWithdrawStatusFailed, fields...)
}

func (w *withdrawStatisticError) HandleMonthlyWithdrawAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.WithdrawMonthlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindMonthlyWithdraws, fields...)
}

func (w *withdrawStatisticError) HandleYearlyWithdrawAmountsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.WithdrawYearlyAmountResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedFindYearlyWithdraws, fields...)
}

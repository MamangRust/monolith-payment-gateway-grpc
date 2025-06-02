package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoStatisticError struct {
	logger logger.LoggerInterface
}

func NewSaldoStatisticError(logger logger.LoggerInterface) *saldoStatisticError {
	return &saldoStatisticError{
		logger: logger,
	}
}

func (e *saldoStatisticError) HandleMonthlyTotalSaldoBalanceError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.SaldoMonthTotalBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindMonthlyTotalSaldoBalance,
		fields...,
	)
}

func (e *saldoStatisticError) HandleYearlyTotalSaldoBalanceError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.SaldoYearTotalBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindYearTotalSaldoBalance,
		fields...,
	)
}

func (e *saldoStatisticError) HandleMonthlySaldoBalancesError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.SaldoMonthBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindMonthlySaldoBalances,
		fields...,
	)
}

func (e *saldoStatisticError) HandleYearlySaldoBalancesError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.SaldoYearBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindYearlySaldoBalances,
		fields...,
	)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// saldoStatisticError is a struct that implements the SaldoStatisticError interface
type saldoStatisticError struct {
	logger logger.LoggerInterface
}

// NewSaldoStatisticError initializes a new SaldoStatisticError with the provided logger.
// It returns an instance of the saldoStatisticError struct.
func NewSaldoStatisticError(logger logger.LoggerInterface) SaldoStatisticErrorHandler {
	return &saldoStatisticError{
		logger: logger,
	}
}

// HandleMonthlyTotalSaldoBalanceError processes errors during the retrieval of a monthly total saldo balance.
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
//   - A slice of SaldoMonthTotalBalanceResponse with error details and a standardized ErrorResponse.
func (e *saldoStatisticError) HandleMonthlyTotalSaldoBalanceError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.SaldoMonthTotalBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindMonthlyTotalSaldoBalance,
		fields...,
	)
}

// HandleYearlyTotalSaldoBalanceError processes errors during the retrieval of a yearly total saldo balance.
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
//   - A slice of SaldoYearTotalBalanceResponse with error details and a standardized ErrorResponse.
func (e *saldoStatisticError) HandleYearlyTotalSaldoBalanceError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.SaldoYearTotalBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindYearTotalSaldoBalance,
		fields...,
	)
}

// HandleMonthlySaldoBalancesError processes errors during the retrieval of monthly saldo balances.
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
//   - A slice of SaldoMonthBalanceResponse with error details and a standardized ErrorResponse.
func (e *saldoStatisticError) HandleMonthlySaldoBalancesError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.SaldoMonthBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindMonthlySaldoBalances,
		fields...,
	)
}

// HandleYearlySaldoBalancesError processes errors during the retrieval of yearly saldo balances.
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
//   - A slice of SaldoYearBalanceResponse with error details and a standardized ErrorResponse.
func (e *saldoStatisticError) HandleYearlySaldoBalancesError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.SaldoYearBalanceResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedFindYearlySaldoBalances,
		fields...,
	)
}

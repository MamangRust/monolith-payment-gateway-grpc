package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SaldoCommandErrorHandler interface {
	HandleFindCardByNumberError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	HandleCreateSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	HandleUpdateSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	HandleTrashSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	HandleRestoreSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
	HandleDeleteSaldoPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleRestoreAllSaldoError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleDeleteAllSaldoPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type SaldoQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.SaldoResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.SaldoResponse, *response.ErrorResponse)
}

type SaldoStatisticErrorHandler interface {
	HandleMonthlyTotalSaldoBalanceError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse)
	HandleYearlyTotalSaldoBalanceError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse)
	HandleMonthlySaldoBalancesError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse)
	HandleYearlySaldoBalancesError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse)
}

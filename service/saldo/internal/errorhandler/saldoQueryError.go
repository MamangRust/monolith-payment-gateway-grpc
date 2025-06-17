package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoQueryError struct {
	logger logger.LoggerInterface
}

func NewSaldoQueryError(logger logger.LoggerInterface) *saldoQueryError {
	return &saldoQueryError{
		logger: logger,
	}
}

func (e *saldoQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.SaldoResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.SaldoResponse](e.logger, err, method, tracePrefix, span, status, saldo_errors.ErrFailedFindAllSaldos, fields...)
}

func (e *saldoQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.SaldoResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *saldoQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

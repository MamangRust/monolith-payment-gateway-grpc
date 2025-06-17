package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawQueryError struct {
	logger logger.LoggerInterface
}

func NewWithdrawQueryError(logger logger.LoggerInterface) *withdrawQueryError {
	return &withdrawQueryError{
		logger: logger,
	}
}

func (w *withdrawQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (w *withdrawQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.WithdrawResponseDeleteAt](w.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (w *withdrawQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupQueryError struct {
	logger logger.LoggerInterface
}

func NewTopupQueryError(logger logger.LoggerInterface) *topupQueryError {
	return &topupQueryError{
		logger: logger,
	}
}

func (e *topupQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TopupResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TopupResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *topupQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TopupResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *topupQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

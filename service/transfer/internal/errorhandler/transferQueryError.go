package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferQueryError struct {
	logger logger.LoggerInterface
}

func NewTransferQueryError(logger logger.LoggerInterface) *transferQueryError {
	return &transferQueryError{
		logger: logger,
	}
}

func (t *transferQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransferResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedFindAllTransfers, fields...)
}

func (t *transferQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransferResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transferQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transferQueryError) HanldeRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

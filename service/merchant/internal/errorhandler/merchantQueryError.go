package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantQueryError(logger logger.LoggerInterface) *merchantQueryError {
	return &merchantQueryError{
		logger: logger,
	}
}

func (e *merchantQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.MerchantResponse](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...,
	)
}

func (e *merchantQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.MerchantResponseDeleteAt](
		e.logger, err, method, tracePrefix, span, status, errResp, fields...,
	)
}

func (e *merchantQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponse](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...,
	)
}

func (e *merchantQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantResponse](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrMerchantNotFoundRes, fields...,
	)
}

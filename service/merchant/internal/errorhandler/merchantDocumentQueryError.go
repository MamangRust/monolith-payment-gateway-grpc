package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDocumentQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantDocumentQueryError(logger logger.LoggerInterface) *merchantDocumentQueryError {
	return &merchantDocumentQueryError{
		logger: logger,
	}
}

func (e *merchantDocumentQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.MerchantDocumentResponse](e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...)
}

func (e *merchantDocumentQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.MerchantDocumentResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...)
}

func (e *merchantDocumentQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantDocumentResponse](e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindAllMerchants, fields...)
}

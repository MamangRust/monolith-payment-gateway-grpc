package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantTransactionError struct {
	logger logger.LoggerInterface
}

func NewMerchantTransactionError(logger logger.LoggerInterface) *merchantTransactionError {
	return &merchantTransactionError{
		logger: logger,
	}
}

func (e *merchantTransactionError) HandleRepositoryAllError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantTransactionResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		merchant_errors.ErrFailedFindAllTransactions,
		fields...,
	)
}

func (e *merchantTransactionError) HandleRepositoryByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantTransactionResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		merchant_errors.ErrFailedFindAllTransactionsByMerchant,
		fields...,
	)
}

func (e *merchantTransactionError) HandleRepositoryByApiKeyError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantTransactionResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		merchant_errors.ErrFailedFindAllTransactionsByApikey,
		fields...,
	)
}

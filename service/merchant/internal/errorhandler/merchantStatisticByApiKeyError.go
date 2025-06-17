package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisticByApiKeyError struct {
	logger logger.LoggerInterface
}

func NewMerchantStatisticByApiKeyError(logger logger.LoggerInterface) *merchantStatisticByApiKeyError {
	return &merchantStatisticByApiKeyError{
		logger: logger,
	}
}

func (e *merchantStatisticByApiKeyError) HandleMonthlyPaymentMethodByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyPaymentMethodByApikeys, fields...,
	)
}

func (e *merchantStatisticByApiKeyError) HandleYearlyPaymentMethodByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyPaymentMethodByApikeys, fields...,
	)
}

func (e *merchantStatisticByApiKeyError) HandleMonthlyAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyAmountByApikeys, fields...,
	)
}

func (e *merchantStatisticByApiKeyError) HandleYearlyAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyAmountByApikeys, fields...,
	)
}

func (e *merchantStatisticByApiKeyError) HandleMonthlyTotalAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyTotalAmountByApikeys, fields...,
	)
}

func (e *merchantStatisticByApiKeyError) HandleYearlyTotalAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyTotalAmountByApikeys, fields...,
	)
}

package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisticByMerchantError struct {
	logger logger.LoggerInterface
}

func NewMerchantStatisticByMerchantError(logger logger.LoggerInterface) *merchantStatisticByMerchantError {
	return &merchantStatisticByMerchantError{
		logger: logger,
	}
}

func (e *merchantStatisticByMerchantError) HandleMonthlyPaymentMethodByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyPaymentMethodByMerchants, fields...,
	)
}

func (e *merchantStatisticByMerchantError) HandleYearlyPaymentMethodByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyPaymentMethodByMerchants, fields...,
	)
}

func (e *merchantStatisticByMerchantError) HandleMonthlyAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyAmountByMerchants, fields...,
	)
}

func (e *merchantStatisticByMerchantError) HandleYearlyAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyAmountByMerchants, fields...,
	)
}

func (e *merchantStatisticByMerchantError) HandleMonthlyTotalAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyTotalAmountByMerchants, fields...,
	)
}

func (e *merchantStatisticByMerchantError) HandleYearlyTotalAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyTotalAmountByMerchants, fields...,
	)
}

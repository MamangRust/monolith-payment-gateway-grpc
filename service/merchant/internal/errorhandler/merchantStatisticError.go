package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisticError struct {
	logger logger.LoggerInterface
}

func NewMerchantStatisticError(logger logger.LoggerInterface) *merchantStatisticError {
	return &merchantStatisticError{
		logger: logger,
	}
}

func (e *merchantStatisticError) HandleMonthlyPaymentMethodsMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponseMonthlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyPaymentMethodsMerchant, fields...,
	)
}

func (e *merchantStatisticError) HandleYearlyPaymentMethodMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponseYearlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyPaymentMethodMerchant, fields...,
	)
}

func (e *merchantStatisticError) HandleMonthlyAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	statuus *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponseMonthlyAmount](
		e.logger, err, method, tracePrefix, span, nil, merchant_errors.ErrFailedFindMonthlyAmountMerchant, fields...,
	)
}

func (e *merchantStatisticError) HandleYearlyAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponseYearlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyAmountMerchant, fields...,
	)
}

func (e *merchantStatisticError) HandleMonthlyTotalAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponseMonthlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyTotalAmountMerchant, fields...,
	)
}

func (e *merchantStatisticError) HandleYearlyTotalAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.MerchantResponseYearlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyTotalAmountMerchant, fields...,
	)
}

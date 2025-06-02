package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantCommandError(logger logger.LoggerInterface) *merchantCommandError {
	return &merchantCommandError{
		logger: logger,
	}
}

func (e *merchantCommandError) HandleCreateMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedCreateMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleUpdateMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedUpdateMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleUpdateMerchantStatusError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedUpdateMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleTrashedMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedTrashMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleRestoreMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedRestoreMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleDeleteMerchantPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedDeleteMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleRestoreAllMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedRestoreAllMerchants,
		fields...,
	)
}

func (e *merchantCommandError) HandleDeleteAllMerchantPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedDeleteAllMerchants,
		fields...,
	)
}

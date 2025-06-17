package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStatisticError struct {
	logger logger.LoggerInterface
}

func NewTopupStatisticError(logger logger.LoggerInterface) *topupStatisticError {
	return &topupStatisticError{
		logger: logger,
	}
}

func (e *topupStatisticError) HandleMonthTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseMonthStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusSuccess, fields...)
}

func (e *topupStatisticError) HandleYearlyTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseYearStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusSuccess, fields...)
}

func (e *topupStatisticError) HandleMonthTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseMonthStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusFailed, fields...)
}

func (e *topupStatisticError) HandleYearlyTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupResponseYearStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusFailed, fields...)
}

func (e *topupStatisticError) HandleMonthlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupMonthMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupMethods, fields...)
}

func (e *topupStatisticError) HandleYearlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupYearlyMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupMethods, fields...)
}

func (e *topupStatisticError) HandleMonthlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupMonthAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupAmounts, fields...)
}

func (e *topupStatisticError) HandleYearlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TopupYearlyAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupAmounts, fields...)
}

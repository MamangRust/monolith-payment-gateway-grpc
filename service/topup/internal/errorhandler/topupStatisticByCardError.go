package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStatisticByCardError struct {
	logger logger.LoggerInterface
}

func NewTopupStatisticByCardError(logger logger.LoggerInterface) *topupStatisticByCardError {
	return &topupStatisticByCardError{
		logger: logger,
	}
}

func (e *topupStatisticByCardError) HandleMonthTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupResponseMonthStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusSuccessByCard, fields...)
}
func (e *topupStatisticByCardError) HandleYearlyTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupResponseYearStatusSuccess](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusSuccessByCard, fields...)
}

func (e *topupStatisticByCardError) HandleMonthTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupResponseMonthStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthTopupStatusFailedByCard, fields...)
}

func (e *topupStatisticByCardError) HandleYearlyTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupResponseYearStatusFailed](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupStatusFailedByCard, fields...)
}

func (e *topupStatisticByCardError) HandleMonthlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupMonthMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupMethodsByCard, fields...)
}

func (e *topupStatisticByCardError) HandleYearlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupYearlyMethodResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupMethodsByCard, fields...)
}

func (e *topupStatisticByCardError) HandleMonthlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupMonthAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindMonthlyTopupAmountsByCard, fields...)
}

func (e *topupStatisticByCardError) HandleYearlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	return handleErrorTemplate[[]*response.TopupYearlyAmountResponse](e.logger, err, method, tracePrefix, span, status, topup_errors.ErrFailedFindYearlyTopupAmountsByCard, fields...)
}

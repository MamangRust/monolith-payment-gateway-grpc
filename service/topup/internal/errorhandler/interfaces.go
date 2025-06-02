package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TopupQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TopupResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TopupResponse, *response.ErrorResponse)
}

type TopupStatisticErrorHandler interface {
	HandleMonthTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	HandleYearlyTopupStatusSuccess(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	HandleMonthTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	HandleYearlyTopupStatusFailed(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	HandleMonthlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	HandleYearlyTopupMethods(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	HandleMonthlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	HandleYearlyTopupAmounts(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

type TopupStatisticByCardErrorHandler interface {
	HandleMonthTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	HandleYearlyTopupStatusSuccessByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	HandleMonthTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	HandleYearlyTopupStatusFailedByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	HandleMonthlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	HandleYearlyTopupMethodsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	HandleMonthlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	HandleYearlyTopupAmountsByCardNumber(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

type TopupCommandErrorHandler interface {
	HandleInvalidParseTimeError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		rawTime string,
		fields ...zap.Field,
	) (*response.TopupResponse, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TopupResponse, *response.ErrorResponse)

	HandleCreateTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponse, *response.ErrorResponse)
	HandleUpdateTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponse, *response.ErrorResponse)
	HandleTrashedTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteTopupPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)

	HandleRestoreAllTopupError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllTopupPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

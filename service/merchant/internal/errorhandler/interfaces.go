package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MerchantQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponse, *response.ErrorResponse)
}

type MerchantDocumentQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantDocumentResponse, *response.ErrorResponse)
}

type MerchantTransactionErrorHandler interface {
	HandleRepositoryAllError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
	HandleRepositoryByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
	HandleRepositoryByApiKeyError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
}

type MerchantCommandErrorHandler interface {
	HandleCreateMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleUpdateMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleUpdateMerchantStatusError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleTrashedMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleRestoreMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleDeleteMerchantPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleRestoreAllMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleDeleteAllMerchantPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type MerchantDocumentCommandErrorHandler interface {
	HandleCreateMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleUpdateMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleUpdateMerchantDocumentStatusError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleTrashedMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleRestoreMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleDeleteMerchantDocumentPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleRestoreAllMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleDeleteAllMerchantDocumentPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type MerchantStatisticErrorHandler interface {
	HandleMonthlyPaymentMethodsMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)

	HandleYearlyPaymentMethodMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)

	HandleMonthlyAmountMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)

	HandleYearlyAmountMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)

	HandleMonthlyTotalAmountMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)

	HandleYearlyTotalAmountMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}

type MerchantStatisticByMerchantErrorHandler interface {
	HandleMonthlyPaymentMethodByMerchantsError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)

	HandleYearlyPaymentMethodByMerchantsError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)

	HandleMonthlyAmountByMerchantsError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)

	HandleYearlyAmountByMerchantsError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)

	HandleMonthlyTotalAmountByMerchantsError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)

	HandleYearlyTotalAmountByMerchantsError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}

type MerchantStatisticByApikeyErrorHandler interface {
	HandleMonthlyPaymentMethodByApikeysError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)

	HandleYearlyPaymentMethodByApikeysError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)

	HandleMonthlyAmountByApikeysError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)

	HandleYearlyAmountByApikeysError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)

	HandleMonthlyTotalAmountByApikeysError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)

	HandleYearlyTotalAmountByApikeysError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}

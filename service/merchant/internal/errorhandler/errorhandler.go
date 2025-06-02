package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	MerchantQueryError               MerchantQueryErrorHandler
	MerchantCommandError             MerchantCommandErrorHandler
	MerchantDocumentQueryError       MerchantDocumentQueryErrorHandler
	MerchantDocumentCommandError     MerchantDocumentCommandErrorHandler
	MerchantTransactionError         MerchantTransactionErrorHandler
	MerchantStatisticError           MerchantStatisticErrorHandler
	MerchantStatisticByMerchantError MerchantStatisticByMerchantErrorHandler
	MerchantStatisticByApiKeyError   MerchantStatisticByApikeyErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantQueryError: NewMerchantQueryError(logger),
		MerchantCommandError: NewMerchantCommandError(
			logger,
		),
		MerchantDocumentQueryError: NewMerchantDocumentQueryError(
			logger,
		),
		MerchantDocumentCommandError: NewMerchantDocumentCommandError(
			logger),
		MerchantTransactionError:         NewMerchantTransactionError(logger),
		MerchantStatisticError:           NewMerchantStatisticError(logger),
		MerchantStatisticByMerchantError: NewMerchantStatisticByMerchantError(logger),
		MerchantStatisticByApiKeyError:   NewMerchantStatisticByApiKeyError(logger),
	}
}

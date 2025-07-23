package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that handles errors in the application.
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

// NewErrorHandler returns a new instance of ErrorHandler with all its fields initialized
// to their respective error handlers with the given logger.
//
// The returned ErrorHandler instance is ready to be used for error handling in the application.
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

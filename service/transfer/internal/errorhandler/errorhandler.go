package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	TransferQueryError           TransferQueryErrorHandler
	TransferCommandError         TransferCommandErrorHandler
	TransferStatisticError       TransferStatisticErrorHandler
	TransferStatisticByCardError TransferStatisticByCardErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler for the transfer service.
// It takes a logger as input and returns a pointer to the ErrorHandler struct,
// with all its fields initialized to their respective error handlers.
// The returned ErrorHandler instance is ready to be used for error handling
// in transfer operations, ensuring that errors are appropriately logged and managed.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		TransferQueryError:           NewTransferQueryError(logger),
		TransferCommandError:         NewTransferCommandError(logger),
		TransferStatisticError:       NewTransferStatisticError(logger),
		TransferStatisticByCardError: NewTransferStatisticByCardError(logger),
	}
}

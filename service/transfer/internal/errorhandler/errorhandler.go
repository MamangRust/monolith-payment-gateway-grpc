package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	TransferQueryError           TransferQueryErrorHandler
	TransferCommandError         TransferCommandErrorHandler
	TransferStatisticError       TransferStatisticErrorHandler
	TransferStatisticByCardError TransferStatisticByCardErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		TransferQueryError:           NewTransferQueryError(logger),
		TransferCommandError:         NewTransferCommandError(logger),
		TransferStatisticError:       NewTransferStatisticError(logger),
		TransferStatisticByCardError: NewTransferStatisticByCardError(logger),
	}
}

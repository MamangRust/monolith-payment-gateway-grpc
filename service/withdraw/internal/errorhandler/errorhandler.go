package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	WithdrawQueryError           WithdrawQueryErrorHandler
	WithdrawCommandError         WithdrawCommandErrorHandler
	WithdrawStatisticError       WithdrawStatisticErrorHandler
	WithdrawStatisticByCardError WithdrawStatisticByCardErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		WithdrawQueryError:           NewWithdrawQueryError(logger),
		WithdrawCommandError:         NewWithdrawCommandError(logger),
		WithdrawStatisticError:       NewWithdrawStatisticError(logger),
		WithdrawStatisticByCardError: NewWithdrawStatisticByCardError(logger),
	}
}

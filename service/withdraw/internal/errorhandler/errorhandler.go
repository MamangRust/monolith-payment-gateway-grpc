package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that holds the error handlers for the withdraw service.
type ErrorHandler struct {
	WithdrawQueryError           WithdrawQueryErrorHandler
	WithdrawCommandError         WithdrawCommandErrorHandler
	WithdrawStatisticError       WithdrawStatisticErrorHandler
	WithdrawStatisticByCardError WithdrawStatisticByCardErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler with all the other error handlers.
// It takes a logger as input and returns a pointer to the ErrorHandler struct.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		WithdrawQueryError:           NewWithdrawQueryError(logger),
		WithdrawCommandError:         NewWithdrawCommandError(logger),
		WithdrawStatisticError:       NewWithdrawStatisticError(logger),
		WithdrawStatisticByCardError: NewWithdrawStatisticByCardError(logger),
	}
}

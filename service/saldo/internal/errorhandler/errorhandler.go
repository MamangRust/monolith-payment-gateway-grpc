package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that implements the SaldoQueryError, SaldoCommandError, and SaldoStatisticError interfaces
type ErrorHandler struct {
	SaldoQueryError     SaldoQueryErrorHandler
	SaldoCommandError   SaldoCommandErrorHandler
	SaldoStatisticError SaldoStatisticErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler with all the other error handlers.
// It takes a logger as input and returns a pointer to the ErrorHandler struct.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		SaldoQueryError:     NewSaldoQueryError(logger),
		SaldoCommandError:   NewSaldoCommandError(logger),
		SaldoStatisticError: NewSaldoStatisticError(logger),
	}
}

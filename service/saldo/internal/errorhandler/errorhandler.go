package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	SaldoQueryError     SaldoQueryErrorHandler
	SaldoCommandError   SaldoCommandErrorHandler
	SaldoStatisticError SaldoStatisticErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		SaldoQueryError:     NewSaldoQueryError(logger),
		SaldoCommandError:   NewSaldoCommandError(logger),
		SaldoStatisticError: NewSaldoStatisticError(logger),
	}
}

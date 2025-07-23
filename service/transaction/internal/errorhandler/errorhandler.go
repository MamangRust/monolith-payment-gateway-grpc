package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that holds all the error handlers for the transaction service.
type ErrorHandler struct {
	TransactionQueryError      TransactionQueryErrorHandler
	TransactonCommandError     TransactionCommandErrorHandler
	TransactionStatisticError  TransactionStatisticErrorHandler
	TransactionStatisticByCard TransactionStatisticByCardErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler with all the other error handlers.
// It takes a logger as input and returns a pointer to the ErrorHandler struct.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		TransactionQueryError:      NewTransactionQueryError(logger),
		TransactonCommandError:     NewTransactionCommandError(logger),
		TransactionStatisticError:  NewTransactionStatisticError(logger),
		TransactionStatisticByCard: NewTransactionStatisticByCardError(logger),
	}
}

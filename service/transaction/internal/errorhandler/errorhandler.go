package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	TransactionQueryError      TransactionQueryErrorHandler
	TransactonCommandError     TransactionCommandErrorHandler
	TransactionStatisticError  TransactionStatisticErrorHandler
	TransactionStatisticByCard TransactionStatisticByCardErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		TransactionQueryError:      NewTransactionQueryError(logger),
		TransactonCommandError:     NewTransactionCommandError(logger),
		TransactionStatisticError:  NewTransactionStatisticError(logger),
		TransactionStatisticByCard: NewTransactionStatisticByCardError(logger),
	}
}

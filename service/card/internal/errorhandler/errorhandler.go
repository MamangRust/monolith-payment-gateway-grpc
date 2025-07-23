package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that contains all the error handlers for the card service
type ErrorHandler struct {
	CardQueryError           CardQueryErrorHandler
	CardCommandError         CardCommandErrorHandler
	CardDashboardError       CardDashboardErrorHandler
	CardStatisticError       CardStatisticErrorHandler
	CardStatisticByCardError CardStatisticByNumberErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler with all the other error handlers.
// It takes a logger as input and returns a pointer to the ErrorHandler struct.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		CardQueryError:           NewCardQueryError(logger),
		CardCommandError:         NewCardCommandError(logger),
		CardDashboardError:       NewCardDashboardError(logger),
		CardStatisticError:       NewCardStatisticError(logger),
		CardStatisticByCardError: NewCardStatisticByNumberError(logger),
	}
}

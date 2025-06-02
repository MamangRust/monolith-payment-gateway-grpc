package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	CardQueryError           CardQueryErrorHandler
	CardCommandError         CardCommandErrorHandler
	CardDashboardError       CardDashboardErrorHandler
	CardStatisticError       CardStatisticErrorHandler
	CardStatisticByCardError cardStatisticByNumberError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		CardQueryError:           NewCardQueryError(logger),
		CardCommandError:         NewCardCommandError(logger),
		CardDashboardError:       NewCardDashboardError(logger),
		CardStatisticError:       NewCardStatisticError(logger),
		CardStatisticByCardError: NewCardStatisticByNumberError(logger),
	}
}

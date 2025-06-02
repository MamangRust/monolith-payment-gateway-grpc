package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	TopupQueryError      TopupQueryErrorHandler
	TopupCommandError    TopupCommandErrorHandler
	TopupStatisticError  TopupStatisticErrorHandler
	TopupStatisticByCard TopupStatisticByCardErrorHandler
}

func NewErrorHandler(
	logger logger.LoggerInterface,
) *ErrorHandler {
	return &ErrorHandler{
		TopupQueryError:      NewTopupQueryError(logger),
		TopupCommandError:    NewTopupCommandError(logger),
		TopupStatisticError:  NewTopupStatisticError(logger),
		TopupStatisticByCard: NewTopupStatisticByCardError(logger),
	}
}

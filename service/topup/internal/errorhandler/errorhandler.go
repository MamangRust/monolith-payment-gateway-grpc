package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that holds the error handlers for the top-up service.
type ErrorHandler struct {
	TopupQueryError      TopupQueryErrorHandler
	TopupCommandError    TopupCommandErrorHandler
	TopupStatisticError  TopupStatisticErrorHandler
	TopupStatisticByCard TopupStatisticByCardErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler for the top-up service.
// It takes a logger as input and returns a pointer to the ErrorHandler struct,
// with all its fields initialized to their respective error handlers.
// The returned ErrorHandler instance is ready to be used for error handling
// in top-up operations, ensuring that errors are appropriately logged and managed.
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

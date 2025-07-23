package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that contains the error handlers for user service operations.
type ErrorHandler struct {
	UserQueryError   UserQueryError
	UserCommandError UserCommandError
}

// NewErrorHandler initializes a new ErrorHandler for the user service.
// It takes a logger as input and returns a pointer to the ErrorHandler struct,
// with all its fields initialized to their respective error handlers.
// The returned ErrorHandler instance is ready to be used for error handling
// in user-related operations, ensuring that errors are appropriately logged and managed.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		UserQueryError:   NewUserQueryError(logger),
		UserCommandError: NewUserCommandError(logger),
	}
}

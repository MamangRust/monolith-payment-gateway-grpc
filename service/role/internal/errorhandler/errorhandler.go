package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a struct that holds the error handlers for the role service.
type ErrorHandler struct {
	RoleQueryError   RoleQueryErrorHandler
	RoleCommandError RoleCommandErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler for the role service.
// It takes a logger as input and returns a pointer to the ErrorHandler struct
// with all its fields initialized to their respective error handlers.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		RoleQueryError:   NewRoleQueryError(logger),
		RoleCommandError: NewRoleCommandError(logger),
	}
}

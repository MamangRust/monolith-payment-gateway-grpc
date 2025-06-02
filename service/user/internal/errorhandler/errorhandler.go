package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	UserQueryError   UserQueryError
	UserCommandError UserCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		UserQueryError:   NewUserQueryError(logger),
		UserCommandError: NewUserCommandError(logger),
	}
}

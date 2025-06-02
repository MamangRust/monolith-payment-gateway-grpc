package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

type ErrorHandler struct {
	RoleQueryError   RoleQueryErrorHandler
	RoleCommandError RoleCommandErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) ErrorHandler {
	return ErrorHandler{
		RoleQueryError:   NewRoleQueryError(logger),
		RoleCommandError: NewRoleCommandError(logger),
	}
}

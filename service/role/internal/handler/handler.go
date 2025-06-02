package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/service"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role RoleHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Role: NewRoleHandleGrpc(deps.Service, deps.Logger),
	}
}

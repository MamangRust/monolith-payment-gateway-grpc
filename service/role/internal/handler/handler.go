package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-role/internal/service"
)

type Deps struct {
	Service service.Service
}

type Handler struct {
	Role RoleHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Role: NewRoleHandleGrpc(deps.Service),
	}
}

package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
)

type Deps struct {
	Service service.Service
}

type Handler struct {
	Card CardHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Card: NewCardHandleGrpc(deps.Service),
	}
}

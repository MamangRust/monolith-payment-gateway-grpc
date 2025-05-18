package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
)

type Deps struct {
	Service service.Service
}

type Handler struct {
	Saldo SaldoHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Saldo: NewSaldoHandleGrpc(deps.Service),
	}
}

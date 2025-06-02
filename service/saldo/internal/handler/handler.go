package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Saldo SaldoHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Saldo: NewSaldoHandleGrpc(deps.Service, deps.Logger),
	}
}

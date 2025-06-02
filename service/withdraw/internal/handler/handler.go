package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Withdraw WithdrawHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Withdraw: NewWithdrawHandleGrpc(deps.Service, deps.Logger),
	}
}

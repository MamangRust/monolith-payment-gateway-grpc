package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Card CardHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Card: NewCardHandleGrpc(deps.Service, deps.Logger),
	}
}

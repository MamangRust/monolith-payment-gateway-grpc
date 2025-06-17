package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Auth AuthHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Auth: NewAuthHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}

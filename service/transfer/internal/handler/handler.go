package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Transfer TransferHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Transfer: NewTransferHandleGrpc(deps.Service, deps.Logger),
	}
}

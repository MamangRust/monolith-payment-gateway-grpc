package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Transaction TransactionHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Transaction: NewTransactionHandleGrpc(deps.Service, deps.Logger),
	}
}

package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
)

type Deps struct {
	Service *service.Service
	Mapper  *protomapper.ProtoMapper
	Logger  logger.LoggerInterface
}

type Handler struct {
	Topup TopupHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Topup: NewTopupHandleGrpc(deps.Service, deps.Logger),
	}
}

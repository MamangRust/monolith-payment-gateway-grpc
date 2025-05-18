package handler

import (
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
)

type Deps struct {
	Service service.Service
	Mapper  protomapper.ProtoMapper
}

type Handler struct {
	Topup TopupHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Topup: NewTopupHandleGrpc(deps.Service),
	}
}

package handler

import (
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
)

type Deps struct {
	Service service.Service
	Mapper  protomapper.ProtoMapper
}

type Handler struct {
	Withdraw WithdrawHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Withdraw: NewWithdrawHandleGrpc(deps.Service),
	}
}

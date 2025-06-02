package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
)

type Deps struct {
	Logger  logger.LoggerInterface
	Service service.Service
}

type Handler struct {
	Merchant         MerchantHandleGrpc
	MerchantDocument MerchantDocumentHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	merchantProto := protomapper.NewMerchantProtoMapper()
	merchantDocumentProto := protomapper.NewMerchantDocumentProtoMapper()

	return &Handler{
		Merchant:         NewMerchantHandleGrpc(deps.Service, merchantProto),
		MerchantDocument: NewMerchantDocumentHandleGrpc(deps.Service, merchantDocumentProto, deps.Logger),
	}
}

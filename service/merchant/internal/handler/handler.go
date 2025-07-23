package handler

import (
	merchantstatshandler "github.com/MamangRust/monolith-payment-gateway-merchant/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomappermerchant "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchant"
	protomappermerchantdocument "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchantdocument"
)

// Deps holds the dependencies required by the handler.
type Deps struct {
	// Logger is used for logging information and errors.
	Logger logger.LoggerInterface
	// Service provides access to the business logic and operations.
	Service service.Service
}

type Handler interface {
	MerchantQueryHandleGrpc
	MerchantCommandHandleGrpc
	MerchantDocumentQueryHandleGrpc
	MerchantDocumentCommandHandleGrpc
	MerchantTransactionHandleGrpc
	merchantstatshandler.HandleStats
}

// Handler contains the gRPC handlers for merchant and merchant document operations.
type handler struct {
	MerchantQueryHandleGrpc
	MerchantCommandHandleGrpc
	MerchantDocumentQueryHandleGrpc
	MerchantDocumentCommandHandleGrpc
	MerchantTransactionHandleGrpc
	merchantstatshandler.HandleStats
}

// NewHandler creates a new Handler instance.
//
// It takes a pointer to a Deps struct as argument, which contains the dependencies
// required to set up the handler.
//
// The handler contains the gRPC handlers for merchant and merchant document operations.
//
// The returned handler is ready to be used.
func NewHandler(deps *Deps) Handler {
	merchantProto := protomappermerchant.NewMerchantProtoMapper()
	merchantDocumentProto := protomappermerchantdocument.NewMerchantDocumentProtoMapper()

	return &handler{
		MerchantQueryHandleGrpc:           NewMerchantQueryHandleGrpc(deps.Service.MerchantQueryService(), deps.Logger, merchantProto.MerchantQueryProtoMapper),
		MerchantCommandHandleGrpc:         NewMerchantCommandHandleGrpc(deps.Service.MerchantCommandService(), deps.Logger, merchantProto.MerchantCommandProtoMapper),
		MerchantDocumentQueryHandleGrpc:   NewMerchantDocumentQueryHandleGrpc(deps.Service.MerchantDocumentQueryService(), deps.Logger, merchantDocumentProto.MerchantDocumentQueryProtoMapper),
		MerchantDocumentCommandHandleGrpc: NewMerchantDocumentCommandHandleGrpc(deps.Service.MerchantDocumentCommandService(), deps.Logger, merchantDocumentProto.MerchantDocumentCommandProtoMapper),
		MerchantTransactionHandleGrpc: NewMerchantTransactionHandleGrpc(
			deps.Service.MerchantTransactionService(), deps.Logger, merchantProto.MerchantTransactionProtoMapper),
		HandleStats: merchantstatshandler.NewMerchantStatsHandler(&merchantstatshandler.DepsStats{
			Service:           deps.Service,
			Logger:            deps.Logger,
			MapperAmount:      merchantProto.MerchantStatsAmountProtoMapper,
			MapperMethod:      merchantProto.MerchantStatsMethodProtoMapper,
			MapperTotalAmount: merchantProto.MerchantStatsTotalAmountProtoMapper,
		}),
	}
}

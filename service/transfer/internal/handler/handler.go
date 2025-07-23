package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transfer"
	transferstatshandler "github.com/MamangRust/monolith-payment-gateway-transfer/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
)

// Deps represents the dependencies required to initialize a Handler.
type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler interface {
	TransferQueryHandleGrpc
	TransferCommandHandleGrpc
	transferstatshandler.HandleStats
}

type handler struct {
	TransferQueryHandleGrpc
	TransferCommandHandleGrpc
	transferstatshandler.HandleStats
}

// NewHandler initializes and returns a new Handler instance for transfer operations.
//
// It takes a Deps struct as a parameter, which contains the necessary
// dependencies like the service and logger required by the handler.
//
// The returned Handler struct includes the gRPC handler for transfer operations.

func NewHandler(deps *Deps) Handler {
	protomapper := protomapper.NewTransferProtoMapper()

	return &handler{
		TransferQueryHandleGrpc:   NewTransferQueryHandler(deps.Service, deps.Logger, protomapper.TransferQueryProtoMapper),
		TransferCommandHandleGrpc: NewTransferCommandHandler(deps.Service, deps.Logger, protomapper.TransferCommandProtoMapper),
		HandleStats: transferstatshandler.NewTransferStatsHandleGrpc(&transferstatshandler.DepsStats{
			Service:      deps.Service,
			Logger:       deps.Logger,
			MapperAmount: protomapper.TransferStatsAmountProtoMapper,
			MapperStatus: protomapper.TransferStatsStatusProtoMapper,
		}),
	}
}

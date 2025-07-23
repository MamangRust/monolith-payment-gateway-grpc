package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/topup"
	topupstatshandler "github.com/MamangRust/monolith-payment-gateway-topup/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
)

// Deps is a struct that holds the dependencies required by the handler
type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler interface {
	TopupQueryHandleGrpc
	TopupCommandHandleGrpc
	topupstatshandler.HandleStats
}

type handler struct {
	TopupQueryHandleGrpc
	TopupCommandHandleGrpc
	topupstatshandler.HandleStats
}

func NewHandler(deps *Deps) Handler {
	mapper := protomapper.NewTopupProtoMapper()

	return &handler{
		TopupQueryHandleGrpc:   NewTopupQueryHandleGrpc(deps.Service, deps.Logger, mapper.TopupQueryProtoMapper),
		TopupCommandHandleGrpc: NewTopupCommandHandleGrpc(deps.Service, mapper.TopupCommandProtoMapper, deps.Logger),
		HandleStats: topupstatshandler.NewTopupStatsHandleGrpc(&topupstatshandler.DepsStats{
			Service:      deps.Service,
			Logger:       deps.Logger,
			MapperAmount: mapper.TopupStatsAmountProtoMapper,
			MapperMethod: mapper.TopupStatsMethodProtoMapper,
			MapperStatus: mapper.TopupStatsStatusProtoMapper,
		}),
	}
}

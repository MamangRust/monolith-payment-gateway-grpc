package topupstatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/topup"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
)

type DepsStats struct {
	Service service.Service

	Logger logger.LoggerInterface

	MapperAmount protomapper.TopupStatsAmountProtoMapper

	MapperMethod protomapper.TopupStatsMethodProtoMapper

	MapperStatus protomapper.TopupStatsStatusProtoMapper
}

type HandleStats interface {
	TopupStatsAmountHandleGrpc
	TopupStatsMethodHandleGrpc
	TopupStatsStatusHandleGrpc
}

type handlerStats struct {
	TopupStatsAmountHandleGrpc
	TopupStatsMethodHandleGrpc
	TopupStatsStatusHandleGrpc
}

func NewTopupStatsHandleGrpc(deps *DepsStats) HandleStats {
	return &handlerStats{
		TopupStatsAmountHandleGrpc: NewTopupStatsAmountHandleGrpc(deps.Service, deps.Logger, deps.MapperAmount),
		TopupStatsMethodHandleGrpc: NewTopupStatsMethodHandleGrpc(deps.Service, deps.Logger, deps.MapperMethod),
		TopupStatsStatusHandleGrpc: NewTopupStatsStatusHandleGrpc(deps.Service, deps.Logger, deps.MapperStatus),
	}
}

package transferstatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
)

type DepsStats struct {
	Service service.Service

	Logger logger.LoggerInterface

	MapperAmount protomapper.TransferStatsAmountProtoMapper

	MapperStatus protomapper.TransferStatsStatusProtoMapper
}

type HandleStats interface {
	TransferStatsAmountHandleGrpc
	TransferStatsStatusHandleGrpc
}

type handleStats struct {
	TransferStatsAmountHandleGrpc
	TransferStatsStatusHandleGrpc
}

func NewTransferStatsHandleGrpc(deps *DepsStats) HandleStats {

	return &handleStats{
		TransferStatsAmountHandleGrpc: NewTransferStatsAmountHandler(deps.Service, deps.Logger, deps.MapperAmount),
		TransferStatsStatusHandleGrpc: NewTransferStatsStatusHandler(deps.Service, deps.Logger, deps.MapperStatus),
	}
}

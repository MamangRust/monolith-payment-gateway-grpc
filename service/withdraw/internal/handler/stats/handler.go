package withdrawstatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
)

type DepsStats struct {
	Service service.Service
	Logger  logger.LoggerInterface

	MapperAmount protomapper.WithdrawaStatsAmountProtoMapper
	MapperStatus protomapper.WithdrawStatsStatusProtoMapper
}

type HandleStats interface {
	WithdrawStatsAmountHandlerGrpc
	WithdrawStatsStatusHandleGrpc
}

type handleStats struct {
	WithdrawStatsAmountHandlerGrpc
	WithdrawStatsStatusHandleGrpc
}

func NewWithdrawStatsHandleGrpc(deps *DepsStats) HandleStats {
	return &handleStats{
		WithdrawStatsAmountHandlerGrpc: NewWithdrawStatsAmountHandleGrpc(deps.Service, deps.Logger, deps.MapperAmount),
		WithdrawStatsStatusHandleGrpc:  NewWithdrawStatsStatusHandleGrpc(deps.Service, deps.Logger, deps.MapperStatus),
	}
}

package saldostatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/saldo"
)

type DepsStats struct {
	Service service.Service
	Logger  logger.LoggerInterface

	MapperBalance      protomapper.SaldoStatsBalanceProtoMapper
	MapperTotalBalance protomapper.SaldoStatsTotalSaldoProtoMapper
}

type HandleStats interface {
	SaldoStatsBalanceHandleGrpc
	SaldoStatsTotalBalanceHandleGrpc
}

type handleStats struct {
	SaldoStatsBalanceHandleGrpc
	SaldoStatsTotalBalanceHandleGrpc
}

func NewSaldoStatsHandle(deps *DepsStats) HandleStats {
	return &handleStats{
		SaldoStatsBalanceHandleGrpc:      NewSaldoStatsBalanceHandleGrpc(deps.Service, deps.Logger, deps.MapperBalance),
		SaldoStatsTotalBalanceHandleGrpc: NewSaldoStatsTotalBalanceHandleGrpc(deps.Service, deps.Logger, deps.MapperTotalBalance),
	}
}

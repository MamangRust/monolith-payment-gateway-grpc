package handlerstats

import (
	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
)

type DepsStats struct {
	Service       service.Service
	Logger        logger.LoggerInterface
	MapperBalance protomapper.CardStatsBalanceProtoMapper
	MapperAmount  protomapper.CardStatsAmountProtoMapper
}

type HandlerStats interface {
	CardStatsBalanceService
	CardStatsTransferService
	CardStatsWithdrawService
	CardStatsTopupService
	CardStatsTransactionService
}

type handlerStats struct {
	CardStatsBalanceService
	CardStatsTransferService
	CardStatsWithdrawService
	CardStatsTopupService
	CardStatsTransactionService
}

func NewHandlerStats(deps *DepsStats) HandlerStats {
	return &handlerStats{
		NewCardStatsBalanceGrpc(deps.Service, deps.Service, deps.Logger, deps.MapperBalance),
		NewCardStatsTransferGrpc(deps.Service, deps.Service, deps.Logger, deps.MapperAmount),
		NewCardStatsWithdrawGrpc(deps.Service, deps.Service, deps.Logger, deps.MapperAmount),
		NewCardStatsTopupGrpc(deps.Service, deps.Service, deps.Logger, deps.MapperAmount),
		NewCardStatsTransactionGrpc(deps.Service, deps.Service, deps.Logger, deps.MapperAmount),
	}
}

package cardstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/stats"
	repositorystats "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
)

type CardStatsService interface {
	CardStatsBalanceService
	CardStatsTopupService
	CardStatsWithdrawService
	CardStatsTransferService
	CardStatsTransactionService
}

type cardStatsService struct {
	CardStatsBalanceService
	CardStatsTopupService
	CardStatsWithdrawService
	CardStatsTransferService
	CardStatsTransactionService
}

type DepsStats struct {
	Mencache      mencache.CardStatsCache
	ErrorHandler  errorhandler.CardStatisticErrorHandler
	Repositories  repositorystats.CardStatsRepository
	Logger        logger.LoggerInterface
	MapperBalance responseservice.CardStatisticBalanceResponseMapper
	MapperAmount  responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsService(deps *DepsStats) CardStatsService {
	return &cardStatsService{
		NewCardStatsBalanceService(&cardStatsBalanceServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperBalance,
		}),
		NewCardStatsTopupService(&cardStatsTopupServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewCardStatsWithdrawService(&cardStatsWithdrawServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewCardStatsTransferService(&cardStatsTransferServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewCardStatsTransactionService(&cardStatsTransactionServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
	}
}

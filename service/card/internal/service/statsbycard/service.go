package cardstatsbycard

import (
	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/statsbycard"
	repositorystats "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
)

type CardStatsByCardService interface {
	CardStatsBalanceByCardService
	CardStatsTopupByCardService
	CardStatsWithdrawByCardService
	CardStatsTransferByCardService
	CardStatsTransactionByCardService
}

type cardStatsByCardService struct {
	CardStatsBalanceByCardService
	CardStatsTopupByCardService
	CardStatsWithdrawByCardService
	CardStatsTransferByCardService
	CardStatsTransactionByCardService
}

type DepsStatsByCard struct {
	Mencache      mencache.CardStatsByCardCache
	ErrorHandler  errorhandler.CardStatisticByNumberErrorHandler
	Repositories  repositorystats.CardStatsByCardRepository
	Logger        logger.LoggerInterface
	MapperBalance responseservice.CardStatisticBalanceResponseMapper
	MapperAmount  responseservice.CardStatisticAmountResponseMapper
}

func NewCardStatsByCardService(deps *DepsStatsByCard) CardStatsByCardService {
	return &cardStatsByCardService{
		NewCardStatsBalanceByCardService(&cardStatsBalanceByCardServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperBalance,
		}),
		NewCardStatsTopupByCardService(&cardStatsTopupByCardServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewCardStatsWithdrawByCardService(&cardStatsWithdrawByCardServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewCardStatsTransferByCardService(&cardStatsTransferByCardServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewCardStatsTransactionByCardService(&cardStatsTransactionByCardServiceDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
	}
}

package saldostatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository/stats"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/saldo"
)

type SaldoStatsService interface {
	SaldoStatsBalanceService
	SaldoStatsTotalBalanceService
}

type saldoStatsService struct {
	SaldoStatsBalanceService
	SaldoStatsTotalBalanceService
}

type DepsStats struct {
	Mencache mencache.SaldoStatsCache

	Errorhandler errorhandler.SaldoStatisticErrorHandler

	Logger logger.LoggerInterface

	Repository repository.SaldoStatsRepository

	MapperBalance responseservice.SaldoStatisticBalanceResponseMapper

	MapperTotalBalance responseservice.SaldoStatisticTotalBalanceResponseMapper
}

func NewSaldoStatsService(deps *DepsStats) SaldoStatsService {
	return &saldoStatsService{
		NewSaldoStatsBalanceService(&saldoStatsBalanceDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.Errorhandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperBalance,
		}),
		NewSaldoStatsTotalBalanceService(&saldoStatsTotalBalanceDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.Errorhandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperTotalBalance,
		}),
	}
}

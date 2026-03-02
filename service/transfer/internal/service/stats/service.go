package transferstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/stats"
)

type TransferStatsService interface {
	TransferStatsAmountService
	TransferStatsStatusService
}

type transferStatsService struct {
	TransferStatsAmountService
	TransferStatsStatusService
}

type DepsStats struct {
	Cache         mencache.TransferStatsCache
	Observability observability.TraceLoggerObservability
	Logger        logger.LoggerInterface
	Repository    repository.TransferStatsRepository
}

func NewTransferStatsService(deps *DepsStats) TransferStatsService {
	return &transferStatsService{
		TransferStatsAmountService: NewTransferStatsAmountService(&transferStatsAmountDeps{
			Cache:         deps.Cache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		TransferStatsStatusService: NewTransferStatsStatusService(&transferStatsStatusDeps{
			Cache:         deps.Cache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}

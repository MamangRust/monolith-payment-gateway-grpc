package merchantstatsservice

import (
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
)

type MerchantStatsService interface {
	MerchantStatsAmountService
	MerchantStatsTotalAmountService
	MerchantStatsMethodService
}

type merchantStatsService struct {
	MerchantStatsAmountService
	MerchantStatsTotalAmountService
	MerchantStatsMethodService
}

type DepsStats struct {
	Mencache      mencache.MerchantStatsCache
	Logger        logger.LoggerInterface
	Repository    repository.MerchantStatsRepository
	Observability observability.TraceLoggerObservability
}

func NewMerchantStatsService(deps *DepsStats) MerchantStatsService {
	return &merchantStatsService{
		NewMerchantStatsAmountService(&merchantStatsAmountDeps{
			Cache:         deps.Mencache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		NewMerchantStatsTotalAmountService(&merchantStatsTotalAmountDeps{
			Cache:         deps.Mencache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		NewMerchantStatsMethodService(&merchantStatsMethodDeps{
			Cache:         deps.Mencache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}

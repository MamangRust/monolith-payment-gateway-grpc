package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/repository"
	topupstatsservice "github.com/MamangRust/monolith-payment-gateway-topup/service/stats"
	topupstatsbycardservice "github.com/MamangRust/monolith-payment-gateway-topup/service/statsbycard"
)

type Service interface {
	TopupQueryService
	TopupCommandService
	topupstatsservice.TopupStatsService
	topupstatsbycardservice.TopupStatsByCardService
}

type service struct {
	TopupQueryService
	TopupCommandService
	topupstatsservice.TopupStatsService
	topupstatsbycardservice.TopupStatsByCardService
}

type Deps struct {
	Kafka        *kafka.Kafka
	Cache        *cache.CacheStore
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) Service {
	cache := mencache.NewMencache(deps.Cache)

	observability, _ := observability.NewObservability("topup-service", deps.Logger)

	return &service{
		TopupQueryService:   newTopupQueryService(deps, observability, cache),
		TopupCommandService: newTopupCommandService(deps, observability, cache),
		TopupStatsService: topupstatsservice.NewTopupStatsService(&topupstatsservice.DepsStats{
			Cache:         cache,
			Logger:        deps.Logger,
			Repository:    deps.Repositories,
			Observability: observability,
		}),
		TopupStatsByCardService: topupstatsbycardservice.NewTopupStatsByCardService(&topupstatsbycardservice.DepsStatsByCard{
			Cache:         cache,
			Logger:        deps.Logger,
			Repository:    deps.Repositories,
			Observability: observability,
		}),
	}
}

func newTopupQueryService(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.TopupQueryCache) TopupQueryService {
	return NewTopupQueryService(&topupQueryDeps{
		Cache:         cache,
		Repository:    deps.Repositories,
		Logger:        deps.Logger,
		Observability: observability,
	})
}

func newTopupCommandService(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.TopupCommandCache) TopupCommandService {
	return NewTopupCommandService(&topupCommandDeps{
		Kafka:                  deps.Kafka,
		Cache:                  cache,
		CardRepository:         deps.Repositories,
		TopupQueryRepository:   deps.Repositories,
		TopupCommandRepository: deps.Repositories,
		SaldoRepository:        deps.Repositories,
		Logger:                 deps.Logger,
		Observability:          observability,
	})
}

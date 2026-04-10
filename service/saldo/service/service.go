package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	saldostatsservice "github.com/MamangRust/monolith-payment-gateway-saldo/service/stats"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
)

type service struct {
	SaldoQueryService
	SaldoCommandService
	saldostatsservice.SaldoStatsService
}

type Service interface {
	SaldoQueryService
	SaldoCommandService
	saldostatsservice.SaldoStatsService
}

type Deps struct {
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
	Cache        *cache.CacheStore
}

func NewService(deps *Deps) Service {
	observability, _ := observability.NewObservability("saldo-service", deps.Logger)
	cache := mencache.NewMencache(deps.Cache)

	return &service{
		SaldoQueryService:   newSaldoQueryService(deps, observability, cache),
		SaldoCommandService: newSaldoCommandService(deps, observability, cache),
		SaldoStatsService: saldostatsservice.NewSaldoStatsService(&saldostatsservice.DepsStats{
			Mencache:   cache,
			Logger:     deps.Logger,
			Repository: deps.Repositories,
		}),
	}
}

func newSaldoQueryService(deps *Deps, observabilty observability.TraceLoggerObservability, cache mencache.Mencache) SaldoQueryService {
	return NewSaldoQueryService(&saldoQueryParams{
		Cache:         cache,
		Repository:    deps.Repositories,
		Logger:        deps.Logger,
		Observability: observabilty,
	})
}

func newSaldoCommandService(deps *Deps, observabilty observability.TraceLoggerObservability, cache mencache.Mencache) SaldoCommandService {
	return NewSaldoCommandService(&saldoCommandParams{
		Cache:                  cache,
		saldoCommandRepository: deps.Repositories,
		CardRepository:         deps.Repositories,
		Logger:                 deps.Logger,
		Observability:          observabilty,
	})
}

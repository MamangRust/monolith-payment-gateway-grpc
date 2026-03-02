package service

import (
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	cardstatsservice "github.com/MamangRust/monolith-payment-gateway-card/internal/service/stats"
	cardstatsbycard "github.com/MamangRust/monolith-payment-gateway-card/internal/service/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
)

type Service interface {
	CardQueryService
	CardDashboardService
	cardstatsservice.CardStatsService
	cardstatsbycard.CardStatsByCardService
	CardCommandService
}

type service struct {
	CardQueryService
	CardDashboardService
	cardstatsservice.CardStatsService
	cardstatsbycard.CardStatsByCardService
	CardCommandService
}

type Deps struct {
	Cache        *cache.CacheStore
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
	Kafka        *kafka.Kafka
}

func NewService(deps *Deps) Service {
	observability, _ := observability.NewObservability("card-server", deps.Logger)

	cache := mencache.NewMencache(deps.Cache)

	return &service{
		CardQueryService:     newCardQuery(deps, observability, cache),
		CardCommandService:   newCardCommand(deps, observability, cache),
		CardDashboardService: newCardDashboard(deps, observability, cache),
		CardStatsService: cardstatsservice.NewCardStatsService(&cardstatsservice.DepsStats{
			Mencache:      cache,
			Repositories:  deps.Repositories.CardStatistic,
			Logger:        deps.Logger,
			Observability: observability,
		}),
		CardStatsByCardService: cardstatsbycard.NewCardStatsByCardService(&cardstatsbycard.DepsStatsByCard{
			Mencache:      cache,
			Repositories:  deps.Repositories.CardStatisticByCard,
			Logger:        deps.Logger,
			Observability: observability,
		}),
	}
}

// newCardQuery initializes a new instance of the CardQueryService.
// It takes a pointer to Deps and a mapper for CardResponse.
// It returns a pointer to CardQueryService.
func newCardQuery(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.Mencache) CardQueryService {
	return NewCardQueryService(&cardQueryServiceDeps{
		Cache:               cache,
		CardQueryRepository: deps.Repositories.CardQuery,
		Logger:              deps.Logger,
		Observability:       observability,
	})
}

// newCardDashboard initializes a new instance of the CardDashboardService.
// It takes a pointer to Deps and a mapper for CardResponse.
// It returns a pointer to CardDashboardService.
func newCardDashboard(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.Mencache) CardDashboardService {
	return NewCardDashboardService(&cardDashboardDeps{
		Cache:                   cache,
		CardDashboardRepository: deps.Repositories.CardDashboard,
		Logger:                  deps.Logger,
		Observability:           observability,
	})
}

// newCardCommand initializes a new instance of the CardCommandService.
// It takes a pointer to Deps and a mapper for CardResponse.
// It returns a pointer to CardCommandService.
func newCardCommand(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.Mencache) CardCommandService {
	return NewCardCommandService(&cardCommandServiceDeps{
		Cache:                 cache,
		Kafka:                 deps.Kafka,
		UserRepository:        deps.Repositories.User,
		CardCommandRepository: deps.Repositories.CardCommand,
		Logger:                deps.Logger,
		Observability:         observability,
	})
}

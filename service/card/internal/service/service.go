package service

import (
	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	cardstatsservice "github.com/MamangRust/monolith-payment-gateway-card/internal/service/stats"
	cardstatsbycard "github.com/MamangRust/monolith-payment-gateway-card/internal/service/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
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
	Mencache     mencache.Mencache
	ErrorHandler *errorhandler.ErrorHandler
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
	Kafka        *kafka.Kafka
}

func NewService(deps *Deps) Service {
	cardMapper := responseservice.NewCardResponseMapper()

	return &service{
		CardQueryService:     newCardQuery(deps, cardMapper.QueryMapper()),
		CardCommandService:   newCardCommand(deps, cardMapper.CommandMapper()),
		CardDashboardService: newCardDashboard(deps),
		CardStatsService: cardstatsservice.NewCardStatsService(&cardstatsservice.DepsStats{
			ErrorHandler:  deps.ErrorHandler.CardStatisticError,
			Mencache:      deps.Mencache,
			Repositories:  deps.Repositories.CardStatistic,
			Logger:        deps.Logger,
			MapperBalance: cardMapper.BalanceStatsMapper(),
			MapperAmount:  cardMapper.AmountStatsMapper(),
		}),
		CardStatsByCardService: cardstatsbycard.NewCardStatsByCardService(&cardstatsbycard.DepsStatsByCard{
			ErrorHandler:  deps.ErrorHandler.CardStatisticByCardError,
			Mencache:      deps.Mencache,
			Repositories:  deps.Repositories.CardStatisticByCard,
			Logger:        deps.Logger,
			MapperBalance: cardMapper.BalanceStatsMapper(),
			MapperAmount:  cardMapper.AmountStatsMapper(),
		}),
	}
}

// newCardQuery initializes a new instance of the CardQueryService.
// It takes a pointer to Deps and a mapper for CardResponse.
// It returns a pointer to CardQueryService.
func newCardQuery(deps *Deps, mapper responseservice.CardQueryResponseMapper) CardQueryService {
	return NewCardQueryService(&cardQueryServiceDeps{
		ErrorHandler:        deps.ErrorHandler.CardQueryError,
		Cache:               deps.Mencache,
		CardQueryRepository: deps.Repositories.CardQuery,
		Logger:              deps.Logger,
		Mapper:              mapper,
	})
}

// newCardDashboard initializes a new instance of the CardDashboardService.
// It takes a pointer to Deps and a mapper for CardResponse.
// It returns a pointer to CardDashboardService.
func newCardDashboard(deps *Deps) CardDashboardService {
	return NewCardDashboardService(&cardDashboardDeps{
		ErrorHandler:            deps.ErrorHandler.CardDashboardError,
		Cache:                   deps.Mencache,
		CardDashboardRepository: deps.Repositories.CardDashboard,
		Logger:                  deps.Logger,
	})
}

// newCardCommand initializes a new instance of the CardCommandService.
// It takes a pointer to Deps and a mapper for CardResponse.
// It returns a pointer to CardCommandService.
func newCardCommand(deps *Deps, mapper responseservice.CardCommandResponseMapper) CardCommandService {
	return NewCardCommandService(&cardCommandServiceDeps{
		ErrorHandler:          deps.ErrorHandler.CardCommandError,
		Cache:                 deps.Mencache,
		Kafka:                 deps.Kafka,
		UserRepository:        deps.Repositories.User,
		CardCommandRepository: deps.Repositories.CardCommand,
		Logger:                deps.Logger,
		Mapper:                mapper,
	})
}

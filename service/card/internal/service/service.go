package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
)

type Service struct {
	CardQuery           CardQueryService
	CardDashboard       CardDashboardService
	CardStatistic       CardStatisticService
	CardStatisticByCard CardStatisticByNumberService
	CardCommand         CardCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
	Kafka        kafka.Kafka
}

func NewService(deps Deps) *Service {
	cardMapper := responseservice.NewCardResponseMapper()

	return &Service{
		CardQuery: NewCardQueryService(
			deps.Ctx,
			deps.Repositories.CardQuery, deps.Logger, cardMapper),
		CardDashboard: NewCardDashboardService(
			deps.Ctx,
			deps.Repositories.CardDashboard,
			deps.Logger,
			cardMapper,
		),
		CardStatistic: NewCardStatisticService(
			deps.Ctx,
			deps.Repositories.CardStatistic, deps.Logger, cardMapper),
		CardStatisticByCard: NewCardStatisticBycardService(
			deps.Ctx,
			deps.Repositories.CardStatisticByCard,
			deps.Logger,
			cardMapper,
		),
		CardCommand: NewCardCommandService(
			deps.Ctx,
			deps.Kafka, deps.Repositories.User, deps.Repositories.CardCommand, deps.Logger, cardMapper),
	}
}

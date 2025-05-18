package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
)

type Service struct {
	TopupQuery           TopupQueryService
	TopupStatistic       TopupStatisticService
	TopupStatisticByCard TopupStatisticByCardService
	TopupCommand         TopupCommandService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
	Mapper       responseservice.ResponseServiceMapper
}

func NewService(deps Deps) *Service {
	topupMapper := responseservice.NewTopupResponseMapper()

	return &Service{
		TopupQuery:           NewTopupQueryService(deps.Ctx, deps.Repositories.TopupQuery, deps.Logger, topupMapper),
		TopupStatistic:       NewTopupStasticService(deps.Ctx, deps.Repositories.TopupStatistic, deps.Logger, topupMapper),
		TopupStatisticByCard: NewTopupStatisticByCardService(deps.Ctx, deps.Repositories.TopupStatistisByCard, deps.Logger, topupMapper),
		TopupCommand: NewTopupCommandService(
			deps.Kafka,
			deps.Ctx,
			deps.Repositories.Card,
			deps.Repositories.TopupQuery,
			deps.Repositories.TopupCommand,
			deps.Repositories.Saldo,
			deps.Logger,
			topupMapper,
		),
	}
}

package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
)

type Service struct {
	TransferQuery           TransferQueryService
	TransferStatistic       TransferStatisticsService
	TransferStatisticByCard TransferStatisticByCardService
	TransferCommand         TransferCommandService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	transferMapper := responseservice.NewTransferResponseMapper()

	return &Service{
		TransferQuery: NewTransferQueryService(
			deps.Ctx,
			deps.Repositories.TransferQuery,
			deps.Logger,
			transferMapper,
		),
		TransferStatistic:       NewTransferStatisticService(deps.Ctx, deps.Repositories.TransferStats, deps.Logger, transferMapper),
		TransferStatisticByCard: NewTransferStatisticByCardService(deps.Ctx, deps.Repositories.Card, deps.Repositories.TransferStatsByCard, deps.Repositories.Saldo, deps.Logger, transferMapper),
		TransferCommand: NewTransferCommandService(
			deps.Kafka,
			deps.Ctx,
			deps.Repositories.Card,
			deps.Repositories.Saldo,
			deps.Repositories.TransferQuery,
			deps.Repositories.TransferCommand,
			deps.Logger,
			transferMapper,
		),
	}
}

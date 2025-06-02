package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
)

type Service struct {
	WithdrawQuery       WithdrawQueryService
	WithdrawCommand     WithdrawCommandService
	WithdrawStats       WithdrawStatisticService
	WithdrawStatsByCard WithdrawStatisticByCardService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
	ErrorHander  errorhandler.ErrorHandler
	Mencache     mencache.Mencache
}

func NewService(deps Deps) *Service {
	withdrawMapper := responseservice.NewWithdrawResponseMapper()

	return &Service{
		WithdrawQuery: NewWithdrawQueryService(deps.Ctx, deps.ErrorHander.WithdrawQueryError, deps.Mencache.WithdrawQueryCache, deps.Repositories.WithdrawQuery, deps.Logger, withdrawMapper),
		WithdrawCommand: NewWithdrawCommandService(
			deps.Ctx,
			deps.ErrorHander.WithdrawCommandError, deps.Mencache.WithdrawCommand,
			deps.Kafka,
			deps.Repositories.Card,
			deps.Repositories.Saldo,
			deps.Repositories.WithdrawCommand,
			deps.Repositories.WithdrawQuery,
			deps.Logger,
			withdrawMapper,
		),
		WithdrawStats: NewWithdrawStatisticService(deps.Ctx, deps.ErrorHander.WithdrawStatisticError, deps.Mencache.WithdrawStatisticCache, deps.Repositories.WithdrawStats, deps.Logger, withdrawMapper),
		WithdrawStatsByCard: NewWithdrawStatisticByCardService(
			deps.Ctx,
			deps.Mencache.WithdrawStatisticByCardCache,
			deps.ErrorHander.WithdrawStatisticByCardError,
			deps.Repositories.WIthdrawStatsByCard, deps.Repositories.Saldo, deps.Logger, withdrawMapper,
		),
	}
}

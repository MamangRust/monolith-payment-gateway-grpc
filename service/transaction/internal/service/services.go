package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionStatistic       TransactionStatisticService
	TransactionStatisticByCard TransactionsStatistcByCardService
	TransactionCommand         TransactionCommandService
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
	transaction := responseservice.NewTransactionResponseMapper()

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Ctx, deps.Mencache.TransactionQueryCache, deps.ErrorHander.TransactionQueryError, deps.Repositories.TransactionQuery, deps.Logger, transaction),
		TransactionStatistic:       NewTransactionStatisticService(deps.Ctx, deps.Mencache.TransactonStatisticCache, deps.ErrorHander.TransactionStatisticError, deps.Repositories.TransactionStat, deps.Logger, transaction),
		TransactionStatisticByCard: NewTransactionStatisticByCardService(deps.Ctx, deps.ErrorHander.TransactionStatisticByCard, deps.Mencache.TransactionStatisticByCardCache, deps.Repositories.TransactionStatByCard, deps.Logger, transaction),
		TransactionCommand: NewTransactionCommandService(
			deps.Kafka,
			deps.Ctx,
			deps.ErrorHander.TransactonCommandError,
			deps.Mencache.TransactionCommandCache,
			deps.Repositories.Merchant,
			deps.Repositories.Card,
			deps.Repositories.Saldo,
			deps.Repositories.TransactionCommand,
			deps.Repositories.TransactionQuery,
			deps.Logger,
			transaction,
		),
	}
}

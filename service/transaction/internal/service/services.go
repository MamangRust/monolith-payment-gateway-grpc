package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	transactionstatsservice "github.com/MamangRust/monolith-payment-gateway-transaction/internal/service/stats"
	transactionstatsbycardservice "github.com/MamangRust/monolith-payment-gateway-transaction/internal/service/statsbycard"
)

// Service is a struct that contains all the services
type service struct {
	TransactionQueryService
	TransactionCommandService
	transactionstatsservice.TransactionStatsService
	transactionstatsbycardservice.TransactionStatsByCardService
}

type Service interface {
	TransactionQueryService
	TransactionCommandService
	transactionstatsservice.TransactionStatsService
	transactionstatsbycardservice.TransactionStatsByCardService
}

// Deps is a struct that contains all the dependencies
type Deps struct {
	Kafka        *kafka.Kafka
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
	ErrorHander  *errorhandler.ErrorHandler
	Mencache     mencache.Mencache
}

// NewService initializes and returns a new instance of Service with all sub-services wired.
func NewService(deps *Deps) Service {
	transaction := responseservice.NewTransactionResponseMapper()

	return &service{
		TransactionQueryService:   newTransactionQueryService(deps, transaction.QueryMapper()),
		TransactionCommandService: newTransactionCommandService(deps, transaction.CommandMapper()),
		TransactionStatsService: transactionstatsservice.NewTransactionStatsService(&transactionstatsservice.DepsStats{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHander.TransactionStatisticError,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			MapperAmount: transaction.AmountStatsMapper(),
			MapperMethod: transaction.MethodStatsMapper(),
			MapperStatus: transaction.StatusStatsMapper(),
		}),
		TransactionStatsByCardService: transactionstatsbycardservice.NewTransactionStatsByCardService(&transactionstatsbycardservice.DepsStats{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHander.TransactionStatisticByCard,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			MapperAmount: transaction.AmountStatsMapper(),
			MapperMethod: transaction.MethodStatsMapper(),
			MapperStatus: transaction.StatusStatsMapper(),
		}),
	}
}

// newTransactionQueryService constructs the TransactionQueryService with its dependencies.
func newTransactionQueryService(deps *Deps, mapper responseservice.TransactionQueryResponseMapper) TransactionQueryService {
	return NewTransactionQueryService(&transactionQueryServiceDeps{
		ErrorHandler:               deps.ErrorHander.TransactionQueryError,
		Cache:                      deps.Mencache,
		TransactionQueryRepository: deps.Repositories,
		Logger:                     deps.Logger,
		Mapper:                     mapper,
	})
}

// newTransactionCommandService constructs the TransactionCommandService with its dependencies.
func newTransactionCommandService(deps *Deps, mapper responseservice.TransactionCommandResponseMapper) TransactionCommandService {
	return NewTransactionCommandService(&transactionCommandServiceDeps{
		Kafka:                        deps.Kafka,
		ErrorHandler:                 deps.ErrorHander.TransactonCommandError,
		Mencache:                     deps.Mencache,
		MerchantRepository:           deps.Repositories,
		CardRepository:               deps.Repositories,
		SaldoRepository:              deps.Repositories,
		TransactionCommandRepository: deps.Repositories,
		TransactionQueryRepository:   deps.Repositories,
		Logger:                       deps.Logger,
		Mapping:                      mapper,
	})
}

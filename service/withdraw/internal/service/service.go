package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	withdrawstatsservice "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service/stats"
	withdrawstatsbycardservice "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service/statsbycard"
)

type Service interface {
	WithdrawQueryService
	WithdrawCommandService
	withdrawstatsservice.WithdrawStatsService
	withdrawstatsbycardservice.WithdrawStatsByCardService
}

type service struct {
	WithdrawQueryService
	WithdrawCommandService
	withdrawstatsservice.WithdrawStatsService
	withdrawstatsbycardservice.WithdrawStatsByCardService
}

// Deps is a struct that contains all the dependencies for the withdraw module
type Deps struct {
	Kafka        *kafka.Kafka
	ErrorHander  *errorhandler.ErrorHandler
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
	Mencache     mencache.Mencache
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a suite of withdraw services including query, command,
// statistics, and statistics by card services. It sets up these services
// using the provided dependencies and response mapper.
func NewService(deps *Deps) Service {
	withdrawMapper := responseservice.NewWithdrawResponseMapper()

	return &service{
		WithdrawQueryService:   newWithdrawQueryService(deps, withdrawMapper.QueryMapper()),
		WithdrawCommandService: newWithdrawCommandService(deps, withdrawMapper.CommandMapper()),
		WithdrawStatsService: withdrawstatsservice.NewWithdrawStatsService(&withdrawstatsservice.DepsStats{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHander.WithdrawStatisticError,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			MapperAmount: withdrawMapper.AmountStatsMapper(),
			MapperStatus: withdrawMapper.StatusStatsMapper(),
		}),
		WithdrawStatsByCardService: withdrawstatsbycardservice.NewWithdrawStatsByCardService(&withdrawstatsbycardservice.DepsStatsByCard{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHander.WithdrawStatisticByCardError,
			Logger:       deps.Logger,
			Repository:   deps.Repositories,
			MapperAmount: withdrawMapper.AmountStatsMapper(),
			MapperStatus: withdrawMapper.StatusStatsMapper(),
		}),
	}
}

// newWithdrawQueryService constructs the WithdrawQueryService with its dependencies.
// It sets up the service to handle withdraw query requests.
//
// Parameters:
// - deps: A pointer to Deps containing the necessary dependencies.
// - mapper: A pointer to responseservice.WithdrawResponseMapper to map WithdrawRecord domain models to WithdrawResponse API-compatible response types.
//
// Returns:
// - A pointer to a newly created WithdrawQueryService.
func newWithdrawQueryService(deps *Deps, mapper responseservice.WithdrawQueryResponseMapper) WithdrawQueryService {
	return NewWithdrawQueryService(
		&withdrawQueryServiceDeps{
			ErrorHandler: deps.ErrorHander.WithdrawQueryError,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       mapper,
		},
	)
}

// newWithdrawCommandService constructs the WithdrawCommandService with its dependencies.
// It sets up the service to handle withdraw command requests.
//
// Parameters:
// - deps: A pointer to Deps containing the necessary dependencies.
// - mapper: A pointer to responseservice.WithdrawResponseMapper to map WithdrawRecord domain models to WithdrawResponse API-compatible response types.
//
// Returns:
// - A pointer to a newly created WithdrawCommandService.
func newWithdrawCommandService(deps *Deps, mapper responseservice.WithdrawCommandResponseMapper) WithdrawCommandService {
	return NewWithdrawCommandService(
		&withdrawCommandServiceDeps{
			ErrorHandler:      deps.ErrorHander.WithdrawCommandError,
			Cache:             deps.Mencache,
			Kafka:             deps.Kafka,
			CardRepository:    deps.Repositories,
			SaldoRepository:   deps.Repositories,
			CommandRepository: deps.Repositories,
			QueryRepository:   deps.Repositories,
			Logger:            deps.Logger,
			Mapper:            mapper,
		},
	)
}

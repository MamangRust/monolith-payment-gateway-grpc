package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	transferstatsservice "github.com/MamangRust/monolith-payment-gateway-transfer/internal/service/stats"
	transferstatsbycardservice "github.com/MamangRust/monolith-payment-gateway-transfer/internal/service/statsbycard"
)

type Service interface {
	TransferQueryService
	TransferCommandService
	transferstatsservice.TransferStatsService
	transferstatsbycardservice.TransferStatsByCardService
}

type service struct {
	TransferQueryService
	TransferCommandService
	transferstatsservice.TransferStatsService
	transferstatsbycardservice.TransferStatsByCardService
}

// Deps is a struct that contains the dependencies for the transfer module
type Deps struct {
	Kafka        *kafka.Kafka
	Mencache     mencache.Mencache
	ErrorHandler *errorhandler.ErrorHandler
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a suite of transfer services including query, statistics,
// and command services. It sets up these services using the provided
// dependencies and response mapper.
func NewService(deps *Deps) Service {
	transferMapper := responseservice.NewTransferResponseMapper()

	return &service{
		TransferQueryService:   newTransferQueryService(deps, transferMapper.QueryMapper()),
		TransferCommandService: newTransferCommandService(deps, transferMapper.CommandMapper()),
		TransferStatsService: transferstatsservice.NewTransferStatsService(&transferstatsservice.DepsStats{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler.TransferStatisticError,
			Logger:       deps.Logger,
			Repository:   deps.Repositories,
			MappeAmount:  transferMapper.AmountStatsMapper(),
			MapperStatus: transferMapper.StatusStatsMapper(),
		}),
		TransferStatsByCardService: transferstatsbycardservice.NewTransferStatsByCardService(&transferstatsbycardservice.DepsStats{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler.TransferStatisticByCardError,
			Logger:       deps.Logger,
			Repository:   deps.Repositories,
			MappeAmount:  transferMapper.AmountStatsMapper(),
			MapperStatus: transferMapper.StatusStatsMapper(),
		}),
	}
}

// newTransferQueryService initializes and returns a new instance of
// TransferQueryService. It sets up the necessary components such as
// context, error handler, cache, repository, logger, and response mapper,
// using the provided dependencies for performing transfer query operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A TransferResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A TransferQueryService, responsible for handling transfer query operations.
func newTransferQueryService(deps *Deps, mapper responseservice.TransferQueryResponseMapper) TransferQueryService {
	return NewTransferQueryService(
		&transferQueryDeps{
			ErrorHandler: deps.ErrorHandler.TransferQueryError,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories,
			Logger:       deps.Logger,
			Mapper:       mapper,
		},
	)
}

// newTransferCommandService initializes and returns a new instance of
// TransferCommandService. It sets up the necessary components, including
// Kafka connection, context, error handler, cache, logger, and response
// mapper, using the provided dependencies to perform transfer command
// operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A TransferResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A TransferCommandService, responsible for handling transfer command operations.
func newTransferCommandService(deps *Deps, mapper responseservice.TransferCommandResponseMapper) TransferCommandService {
	return NewTransferCommandService(
		&transferCommandDeps{
			Kafka:                     deps.Kafka,
			ErrorHandler:              deps.ErrorHandler.TransferCommandError,
			Cache:                     deps.Mencache,
			CardRepository:            deps.Repositories,
			SaldoRepository:           deps.Repositories,
			TransferQueryRepository:   deps.Repositories,
			TransferCommandRepository: deps.Repositories,
			Logger:                    deps.Logger,
			Mapper:                    mapper,
		},
	)
}

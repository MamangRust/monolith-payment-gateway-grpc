package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	topupstatsservice "github.com/MamangRust/monolith-payment-gateway-topup/internal/service/stats"
	topupstatsbycardservice "github.com/MamangRust/monolith-payment-gateway-topup/internal/service/statsbycard"
)

type Service interface {
	TopupQueryService
	TopupCommandService
	topupstatsservice.TopupStatsService
	topupstatsbycardservice.TopupStatsByCardService
}

type service struct {
	TopupQueryService
	TopupCommandService
	topupstatsservice.TopupStatsService
	topupstatsbycardservice.TopupStatsByCardService
}

// Deps holds the dependencies required to construct the topup services.
type Deps struct {
	// Kafka provides the kafka client for publishing messages.
	Kafka *kafka.Kafka
	// ErrorHandler provides the error handler for handling errors.
	ErrorHandler *errorhandler.ErrorHandler
	// Mencache provides the redis client for caching data.
	Mencache mencache.Mencache
	// Ctx provides the context of the service.
	Ctx context.Context
	// Repositories provides the repository for accessing the database.
	Repositories repository.Repositories
	// Logger provides the logger for logging.
	Logger logger.LoggerInterface
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a suite of topup services including query, statistics,
// and command services. It sets up these services using the provided
// dependencies and response mapper.
func NewService(deps *Deps) Service {
	topupMapper := responseservice.NewTopupResponseMapper()

	return &service{
		TopupQueryService:   newTopupQueryService(deps, topupMapper.QueryMapper()),
		TopupCommandService: newTopupCommandService(deps, topupMapper.CommandMapper()),
		TopupStatsService: topupstatsservice.NewTopupStatsService(&topupstatsservice.DepsStats{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler.TopupStatisticError,
			Logger:       deps.Logger,
			Repository:   deps.Repositories,
			MappeAmount:  topupMapper.AmountStatsMapper(),
			MapperMethod: topupMapper.MethodStatsMapper(),
			MapperStatus: topupMapper.StatusStatsMapper(),
		}),
		TopupStatsByCardService: topupstatsbycardservice.NewTopupStatsByCardService(&topupstatsbycardservice.DepsStatsByCard{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler.TopupStatisticByCard,
			Logger:       deps.Logger,
			Repository:   deps.Repositories,
			MappeAmount:  topupMapper.AmountStatsMapper(),
			MapperMethod: topupMapper.MethodStatsMapper(),
			MapperStatus: topupMapper.StatusStatsMapper(),
		}),
	}
}

// newTopupQueryService initializes and returns a new instance of
// TopupQueryService. It sets up the required components, including context,
// error handler, cache, logger, and response mapper, using the provided
// dependencies to perform topup query operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A TopupResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A TopupQueryService, responsible for handling topup query operations.
func newTopupQueryService(deps *Deps, mapper responseservice.TopupQueryResponseMapper) TopupQueryService {
	return NewTopupQueryService(&topupQueryDeps{
		ErrorHandler: deps.ErrorHandler.TopupQueryError,
		Cache:        deps.Mencache,
		Repository:   deps.Repositories,
		Logger:       deps.Logger,
		Mapper:       mapper,
	})
}

// newTopupCommandService initializes and returns a new instance of
// TopupCommandService. It sets up the required components, including Kafka
// connection, context, error handler, cache, logger, and response mapper,
// using the provided dependencies to perform topup command operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A TopupResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A TopupCommandService, responsible for handling topup command operations.
func newTopupCommandService(deps *Deps, mapper responseservice.TopupCommandResponseMapper) TopupCommandService {
	return NewTopupCommandService(&topupCommandDeps{
		Kafka:                  deps.Kafka,
		ErrorHandler:           deps.ErrorHandler.TopupCommandError,
		Cache:                  deps.Mencache,
		CardRepository:         deps.Repositories,
		TopupQueryRepository:   deps.Repositories,
		TopupCommandRepository: deps.Repositories,
		SaldoRepository:        deps.Repositories,
		Logger:                 deps.Logger,
		Mapper:                 mapper,
	})
}

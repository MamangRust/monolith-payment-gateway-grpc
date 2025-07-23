package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	saldostatsservice "github.com/MamangRust/monolith-payment-gateway-saldo/internal/service/stats"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/saldo"
)

// service struct groups all saldo-related domain services.
type service struct {
	SaldoQueryService
	SaldoCommandService
	saldostatsservice.SaldoStatsService
}

// Service interface defines the contract for saldo-related services,
// including query, command, and statistics operations.
type Service interface {
	SaldoQueryService
	SaldoCommandService
	saldostatsservice.SaldoStatsService
}

// Deps holds the external dependencies required to construct the saldo services.
type Deps struct {
	// ErrorHandler provides domain-specific error handlers for saldo operations.
	ErrorHandler *errorhandler.ErrorHandler

	// Mencache provides in-memory caching for query, stats, and command services.
	Mencache mencache.Mencache

	// Repositories provides access to saldo-related data persistence layers.
	Repositories repository.Repositories

	// Logger provides structured and leveled logging support.
	Logger logger.LoggerInterface
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a suite of saldo-related services including query, statistics,
// and command services. It sets up these services using the provided dependencies
// and response mapper.
func NewService(deps *Deps) Service {
	saldoMapper := responseservice.NewSaldoResponseMapper()

	return &service{
		SaldoQueryService:   newSaldoQueryService(deps, saldoMapper.QueryMapper()),
		SaldoCommandService: newSaldoCommandService(deps, saldoMapper.CommandMapper()),
		SaldoStatsService: saldostatsservice.NewSaldoStatsService(&saldostatsservice.DepsStats{
			Mencache:           deps.Mencache,
			Errorhandler:       deps.ErrorHandler.SaldoStatisticError,
			Logger:             deps.Logger,
			Repository:         deps.Repositories,
			MapperBalance:      saldoMapper.StatisticBalanceMapper(),
			MapperTotalBalance: saldoMapper.TotalBalanceMapper(),
		}),
	}
}

// newSaldoQueryService creates a new instance of SaldoQueryService using the provided dependencies and mapper.
// It initializes the service with the context, error handler, cache, repository, logger, and mapper
// from the dependencies. This service handles read-only operations for saldo data.
func newSaldoQueryService(deps *Deps, mapper responseservice.SaldoQueryResponseMapper) SaldoQueryService {
	return NewSaldoQueryService(&saldoQueryParams{
		ErrorHandler: deps.ErrorHandler.SaldoQueryError,
		Cache:        deps.Mencache,
		Repository:   deps.Repositories,
		Logger:       deps.Logger,
		Mapper:       mapper,
	})
}

// newSaldoCommandService creates a new instance of SaldoCommandService using the provided dependencies and mapper.
// It initializes the service with the context, error handler, cache, repository, logger, and mapper
// from the dependencies. This service handles write operations for saldo data, such as top-up and adjustment.
func newSaldoCommandService(deps *Deps, mapper responseservice.SaldoCommandResponseMapper) SaldoCommandService {
	return NewSaldoCommandService(&saldoCommandParams{
		ErrorHandler:    deps.ErrorHandler.SaldoCommandError,
		Cache:           deps.Mencache,
		SaldoRepository: deps.Repositories,
		CardRepository:  deps.Repositories,
		Logger:          deps.Logger,
		Mapper:          mapper,
	})
}

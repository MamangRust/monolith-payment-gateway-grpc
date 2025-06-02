package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	SaldoQuery   SaldoQueryService
	SaldoStats   SaldoStatisticService
	SaldoCommand SaldoCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler errorhandler.ErrorHandler
	Mencache     mencache.Mencache
	Redis        *redis.Client
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	saldoMapper := responseservice.NewSaldoResponseMapper()

	return &Service{
		SaldoQuery:   NewSaldoQueryService(deps.Ctx, deps.ErrorHandler.SaldoQueryError, deps.Mencache.SaldoQueryCache, deps.Repositories.SaldoQuery, deps.Logger, saldoMapper),
		SaldoStats:   NewSaldoStatisticsService(deps.Ctx, deps.ErrorHandler.SaldoStatisticError, deps.Mencache.SaldoStatisticCache, deps.Repositories.SaldoStats, deps.Logger, saldoMapper),
		SaldoCommand: NewSaldoCommandService(deps.Ctx, deps.ErrorHandler.SaldoCommandError, deps.Mencache.SaldoCommandCache, deps.Repositories.SaldoCommand, deps.Repositories.Card, deps.Logger, saldoMapper),
	}
}

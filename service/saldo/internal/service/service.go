package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
)

type Service struct {
	SaldoQuery   SaldoQueryService
	SaldoStats   SaldoStatisticService
	SaldoCommand SaldoCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	saldoMapper := responseservice.NewSaldoResponseMapper()

	return &Service{
		SaldoQuery:   NewSaldoQueryService(deps.Ctx, deps.Repositories.SaldoQuery, deps.Logger, saldoMapper),
		SaldoStats:   NewSaldoStatisticsService(deps.Ctx, deps.Repositories.SaldoStats, deps.Logger, saldoMapper),
		SaldoCommand: NewSaldoCommandService(deps.Ctx, deps.Repositories.SaldoCommand, deps.Repositories.Card, deps.Logger, saldoMapper),
	}
}

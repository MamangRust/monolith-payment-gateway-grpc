package handler

import (
	handlerstats "github.com/MamangRust/monolith-payment-gateway-card/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
)

// Deps represents the dependencies required by the CardHandler.
type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler interface {
	CardQueryService
	CardCommandService
	CardDashboardService
	handlerstats.HandlerStats
}

type handler struct {
	CardQueryService
	CardCommandService
	CardDashboardService
	handlerstats.HandlerStats
}

// NewHandler creates a new CardHandler instance.
func NewHandler(deps *Deps) Handler {
	mapper := protomapper.NewCardProtoMapper()

	return &handler{
		CardQueryService:     NewCardQueryHandleGrpc(deps.Service, deps.Logger, mapper.CardQueryProtoMapper),
		CardCommandService:   NewCardCommandHandleGrpc(deps.Service, deps.Logger, mapper.CardCommandProtoMapper),
		CardDashboardService: NewCardDashboardHandleGrpc(deps.Service, deps.Logger, mapper.CardDashboardProtoMapper),
		HandlerStats: handlerstats.NewHandlerStats(&handlerstats.DepsStats{
			Service:       deps.Service,
			Logger:        deps.Logger,
			MapperBalance: mapper.CardStatsBalanceProtoMapper,
			MapperAmount:  mapper.CardStatsAmountProtoMapper,
		}),
	}
}

package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/withdraw"
	withdrawstatshandler "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
)

// Deps represents the dependencies required by the handler.
type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

// Handler provides methods for handling withdraw requests.
type Handler interface {
	WithdrawQueryHandlerGrpc
	WithdrawCommandHandlerGrpc
	withdrawstatshandler.HandleStats
}

type handler struct {
	WithdrawQueryHandlerGrpc
	WithdrawCommandHandlerGrpc
	withdrawstatshandler.HandleStats
}

// NewHandler creates a new Handler instance.
//
// It takes a pointer to a Deps struct as an argument, which contains the
// dependencies required by the handler.
//
// The returned Handler is ready to be used.
func NewHandler(deps *Deps) Handler {
	mapper := protomapper.NewWithdrawProtoMapper()

	return &handler{
		WithdrawQueryHandlerGrpc:   NewWithdrawQueryHandleGrpc(deps.Service, deps.Logger, mapper.WithdrawQueryProtoMapper),
		WithdrawCommandHandlerGrpc: NewWithdrawCommandHandleGrpc(deps.Service, deps.Logger, mapper.WithdrawCommandProtoMapper),
		HandleStats: withdrawstatshandler.NewWithdrawStatsHandleGrpc(&withdrawstatshandler.DepsStats{
			Service:      deps.Service,
			Logger:       deps.Logger,
			MapperAmount: mapper.WithdrawaStatsAmountProtoMapper,
			MapperStatus: mapper.WithdrawStatsStatusProtoMapper,
		}),
	}
}

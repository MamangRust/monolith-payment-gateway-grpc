package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldostatshandler "github.com/MamangRust/monolith-payment-gateway-saldo/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/saldo"
)

// Deps is a struct that holds the dependencies for the handler
type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

// Handler is a struct that holds the handler
type handler struct {
	SaldoQueryHandleGrpc
	SaldoCommandHandleGrpc
	saldostatshandler.HandleStats
}

type Handler interface {
	SaldoQueryHandleGrpc
	SaldoCommandHandleGrpc
	saldostatshandler.HandleStats
}

// NewHandler creates a new Handler instance.
//
// It takes a pointer to a Deps struct as argument, which contains the dependencies
// required to set up the handler.
//
// The handler contains the gRPC handlers for saldo operations.
//
// The returned handler is ready to be used.
func NewHandler(deps *Deps) Handler {
	protomapepr := protomapper.NewSaldoProtoMapper()

	return &handler{
		SaldoQueryHandleGrpc:   NewSaldoQueryHandleGrpc(deps.Service, deps.Logger, protomapepr.SaldoQueryProtoMapper),
		SaldoCommandHandleGrpc: NewSaldoCommandHandleGrpc(deps.Service, deps.Logger, protomapper.NewSaldoCommandProtoMapper()),
		HandleStats: saldostatshandler.NewSaldoStatsHandle(&saldostatshandler.DepsStats{
			Service:            deps.Service,
			Logger:             deps.Logger,
			MapperBalance:      protomapepr.SaldoStatsBalanceProtoMapper,
			MapperTotalBalance: protomapepr.SaldoStatsTotalSaldoProtoMapper,
		}),
	}
}

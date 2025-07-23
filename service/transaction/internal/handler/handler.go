package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transaction"
	transactionstatshandler "github.com/MamangRust/monolith-payment-gateway-transaction/internal/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
)

// Deps represents the dependencies required to initialize a Handler.
type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler interface {
	TransactionQueryHandleGrpc
	TransactionCommandHandleGrpc
	transactionstatshandler.HandleStats
}

// Handler represents the gRPC handler for transaction operations.
type handler struct {
	TransactionQueryHandleGrpc
	TransactionCommandHandleGrpc
	transactionstatshandler.HandleStats
}

// NewHandler creates a new Handler instance.
func NewHandler(deps *Deps) Handler {
	protomapper := protomapper.NewTransactionProtoMapper()

	return &handler{
		TransactionQueryHandleGrpc:   NewTransactionQueryHandleGrpc(deps.Service, deps.Logger, protomapper.TransactionQueryProtoMapper),
		TransactionCommandHandleGrpc: NewTransactionCommandHandleGrpc(deps.Service, deps.Logger, protomapper.TransactionCommandProtoMapper),
		HandleStats: transactionstatshandler.NewTransactionStatsHandleGrpc(&transactionstatshandler.DepsStats{
			Service:      deps.Service,
			Logger:       deps.Logger,
			MapperAmount: protomapper.TransactionStatsAmountProtoMapper,
			MapperMethod: protomapper.TransactionStatsMethodProtoMapper,
			MapperStatus: protomapper.TransactionStatsStatusProtoMapper,
		}),
	}
}

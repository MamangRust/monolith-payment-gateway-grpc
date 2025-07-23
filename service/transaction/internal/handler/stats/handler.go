package transactionstatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
)

type DepsStats struct {
	Service service.Service

	Logger logger.LoggerInterface

	MapperAmount protomapper.TransactionStatsAmountProtoMapper

	MapperMethod protomapper.TransactionStatsMethodProtoMapper

	MapperStatus protomapper.TransactionStatsStatusProtoMapper
}

type HandleStats interface {
	TransactionStatsAmountHandlerGrpc
	TransactionStatsMethodHandleGrpc
	TransactionStatsStatusHandleGrpc
}

type handlerStats struct {
	TransactionStatsAmountHandlerGrpc
	TransactionStatsMethodHandleGrpc
	TransactionStatsStatusHandleGrpc
}

func NewTransactionStatsHandleGrpc(deps *DepsStats) HandleStats {
	return &handlerStats{
		TransactionStatsAmountHandlerGrpc: NewTransactionStatsAmountHandleGrpc(deps.Service, deps.Logger, deps.MapperAmount),
		TransactionStatsMethodHandleGrpc:  NewTransactionStatsMethodHandleGrpc(deps.Service, deps.Logger, deps.MapperMethod),
		TransactionStatsStatusHandleGrpc:  NewTransactionStatsStatusHandleGrpc(deps.Service, deps.Logger, deps.MapperStatus),
	}
}

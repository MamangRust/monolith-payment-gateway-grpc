package transactionstatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
)

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

func NewTransactionStatsHandleGrpc(service service.Service) HandleStats {
	return &handlerStats{
		TransactionStatsAmountHandlerGrpc: NewTransactionStatsAmountHandleGrpc(service),
		TransactionStatsMethodHandleGrpc:  NewTransactionStatsMethodHandleGrpc(service),
		TransactionStatsStatusHandleGrpc:  NewTransactionStatsStatusHandleGrpc(service),
	}
}

package transactionstatshandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
)

type TransactionStatsAmountHandlerGrpc interface {
	pb.TransactionsStatsAmountServiceServer
}

type TransactionStatsMethodHandleGrpc interface {
	pb.TransactionStatsMethodServiceServer
}

type TransactionStatsStatusHandleGrpc interface {
	pb.TransactionStatsStatusServiceServer
}

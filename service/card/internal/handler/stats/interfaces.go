package handlerstats

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
)

type CardStatsBalanceService interface {
	pb.CardStatsBalanceServiceServer
}

type CardStatsTopupService interface {
	pb.CardStatsTopupServiceServer
}

type CardStatsTransactionService interface {
	pb.CardStatsTransactonServiceServer
}

type CardStatsTransferService interface {
	pb.CardStatsTransferServiceServer
}

type CardStatsWithdrawService interface {
	pb.CardStatsWithdrawServiceServer
}

package withdrawstatshandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
)

type WithdrawStatsAmountHandlerGrpc interface {
	pb.WithdrawStatsAmountServiceServer
}

type WithdrawStatsStatusHandleGrpc interface {
	pb.WithdrawStatsStatusServer
}

package saldostatshandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
)

type SaldoStatsBalanceHandleGrpc interface {
	pb.SaldoStatsBalanceServiceServer
}

type SaldoStatsTotalBalanceHandleGrpc interface {
	pb.SaldoStatsTotalBalanceServer
}

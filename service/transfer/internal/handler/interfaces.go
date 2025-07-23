package handler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
)

type TransferQueryHandleGrpc interface {
	pb.TransferQueryServiceServer
}

type TransferCommandHandleGrpc interface {
	pb.TransferCommandServiceServer
}

type TransferStatsAmountHandleGrpc interface {
	pb.TransferStatsAmountServiceServer
}

type TransferStatsStatusHandleGrpc interface{
	pb.TransferStatsStatusServiceServer
}
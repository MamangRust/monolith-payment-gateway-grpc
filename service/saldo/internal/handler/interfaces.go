package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SaldoHandleGrpc interface {
	pb.SaldoServiceServer

	FindAllSaldo(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldo, error)
	FindByIdSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error)

	FindMonthlyTotalSaldoBalance(ctx context.Context, req *pb.FindMonthlySaldoTotalBalance) (*pb.ApiResponseMonthTotalSaldo, error)
	FindYearTotalSaldoBalance(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearTotalSaldo, error)

	FindMonthlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseMonthSaldoBalances, error)
	FindYearlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearSaldoBalances, error)

	FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseSaldo, error)
	FindByActive(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error)
	CreateSaldo(ctx context.Context, req *pb.CreateSaldoRequest) (*pb.ApiResponseSaldo, error)
	UpdateSaldo(ctx context.Context, req *pb.UpdateSaldoRequest) (*pb.ApiResponseSaldo, error)
	TrashedSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error)
	RestoreSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error)
	DeleteSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDelete, error)

	RestoreAllSaldo(context.Context, *emptypb.Empty) (*pb.ApiResponseSaldoAll, error)
	DeleteAllSaldoPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseSaldoAll, error)
}

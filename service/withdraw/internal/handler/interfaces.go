package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WithdrawHandleGrpc interface {
	pb.WithdrawServiceServer

	FindAllWithdraw(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdraw, error)
	FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error)

	FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error)
	FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error)
	FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawMonthAmount, error)
	FindYearlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawYearAmount, error)

	FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponsesWithdraw, error)
	FindByActive(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error)
	CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error)

	RestoreAllWithdraw(context.Context, *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error)
	DeleteAllWithdrawPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error)
}

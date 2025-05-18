package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransferHandleGrpc interface {
	pb.TransferServiceServer

	FindAllTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransfer, error)
	FindByIdTransfer(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error)

	FindMonthlyTransferStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusSuccess, error)
	FindYearlyTransferStatusSuccess(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusSuccess, error)
	FindMonthlyTransferStatusFailed(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusFailed, error)
	FindYearlyTransferStatusFailed(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusFailed, error)

	FindMonthlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusSuccess, error)
	FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusSuccess, error)
	FindMonthlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusFailed, error)
	FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusFailed, error)

	FindMonthlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferMonthAmount, error)
	FindYearlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearAmount, error)

	FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error)
	FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error)
	FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error)
	FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error)

	FindByTransferByTransferFrom(ctx context.Context, request *pb.FindTransferByTransferFromRequest) (*pb.ApiResponseTransfers, error)
	FindByTransferByTransferTo(ctx context.Context, request *pb.FindTransferByTransferToRequest) (*pb.ApiResponseTransfers, error)
	FindByActiveTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error)
	FindByTrashedTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error)
	CreateTransfer(ctx context.Context, req *pb.CreateTransferRequest) (*pb.ApiResponseTransfer, error)
	UpdateTransfer(ctx context.Context, req *pb.UpdateTransferRequest) (*pb.ApiResponseTransfer, error)
	TrashedTransfer(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error)
	RestoreTransfer(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error)
	DeleteTransferPermanent(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDelete, error)

	RestoreAllTransfer(context.Context, *emptypb.Empty) (*pb.ApiResponseTransferAll, error)
	DeleteAllTransferPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseTransferAll, error)
}

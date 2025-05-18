package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TopupHandleGrpc interface {
	pb.TopupServiceServer

	FindAllTopup(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopup, error)
	FindAllTopupByCardNumber(ctx context.Context, req *pb.FindAllTopupByCardNumberRequest) (*pb.ApiResponsePaginationTopup, error)
	FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error)

	FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error)
	FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error)
	FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error)
	FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error)

	FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error)
	FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error)
	FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error)
	FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error)

	FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error)
	FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error)

	FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error)
	FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error)

	FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error)
	FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error)

	FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error)
	FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error)

	FindByActive(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error)
	CreateTopup(ctx context.Context, req *pb.CreateTopupRequest) (*pb.ApiResponseTopup, error)
	UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error)
	TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error)
	RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error)
	DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error)

	RestoreAllTopup(context.Context, *emptypb.Empty) (*pb.ApiResponseTopupAll, error)
	DeleteAllTopupPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseTopupAll, error)
}

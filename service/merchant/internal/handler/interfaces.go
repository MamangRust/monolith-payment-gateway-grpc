package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MerchantDocumentHandleGrpc interface {
	pb.MerchantDocumentServiceServer
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer

	FindAllMerchant(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error)
	FindByIdMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error)

	FindMonthlyPaymentMethodsMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error)
	FindYearlyPaymentMethodMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyPaymentMethod, error)
	FindMonthlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyAmount, error)
	FindYearlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyAmount, error)

	FindAllTransactionByMerchant(ctx context.Context, req *pb.FindAllMerchantTransaction) (*pb.ApiResponsePaginationMerchantTransaction, error)
	FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error)
	FindYearlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyPaymentMethod, error)
	FindMonthlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyAmount, error)
	FindYearlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyAmount, error)
	FindMonthlyTotalAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyTotalAmount, error)
	FindYearlyTotalAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyTotalAmount, error)

	FindAllTransactionByApikey(ctx context.Context, req *pb.FindAllMerchantApikey) (*pb.ApiResponsePaginationMerchantTransaction, error)
	FindMonthlyPaymentMethodByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error)
	FindYearlyPaymentMethodByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyPaymentMethod, error)
	FindMonthlyAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyAmount, error)
	FindYearlyAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyAmount, error)
	FindMonthlyTotalAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyTotalAmount, error)
	FindYearlyTotalAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyTotalAmount, error)

	FindByApiKey(ctx context.Context, req *pb.FindByApiKeyRequest) (*pb.ApiResponseMerchant, error)

	FindByMerchantUserId(ctx context.Context, req *pb.FindByMerchantUserIdRequest) (*pb.ApiResponsesMerchant, error)
	FindByActive(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error)
	CreateMerchant(ctx context.Context, req *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error)
	UpdateMerchant(ctx context.Context, req *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error)
	TrashedMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error)
	RestoreMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error)
	DeleteMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error)

	RestoreAllMerchant(context.Context, *emptypb.Empty) (*pb.ApiResponseMerchantAll, error)
	DeleteAllMerchantPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseMerchantAll, error)
}

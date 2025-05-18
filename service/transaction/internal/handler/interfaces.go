package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionHandleGrpc interface {
	pb.TransactionServiceServer

	FindAllTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error)
	FindByIdTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error)

	FindMonthlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthMethod, error)
	FindYearlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearMethod, error)
	FindMonthlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthAmount, error)
	FindYearlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearAmount, error)

	FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthMethod, error)
	FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearMethod, error)
	FindMonthlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthAmount, error)
	FindYearlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearAmount, error)

	FindTransactionByMerchantIdRequest(ctx context.Context, request *pb.FindTransactionByMerchantIdRequest) (*pb.ApiResponseTransactions, error)
	FindByActiveTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error)
	FindByTrashedTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error)
	CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error)
	UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error)
	TrashedTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error)
	RestoreTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error)
	DeleteTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error)

	RestoreAllTransaction(context.Context, *emptypb.Empty) (*pb.ApiResponseTransactionAll, error)
	DeleteAllTransactionPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseTransactionAll, error)
}

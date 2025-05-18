package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CardHandleGrpc interface {
	pb.CardServiceServer

	FindAllCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCard, error)
	FindByIdCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error)
	FindByUserIdCard(ctx context.Context, req *pb.FindByUserIdCardRequest) (*pb.ApiResponseCard, error)
	FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseCard, error)

	DashboardCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseDashboardCard, error)
	DashboardCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseDashboardCardNumber, error)

	FindMonthlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseMonthlyBalance, error)
	FindYearlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseYearlyBalance, error)
	FindMonthlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)

	FindMonthlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseMonthlyBalance, error)
	FindYearlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseYearlyBalance, error)
	FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)

	FindMonthlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)

	FindByActiveCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error)
	FindByTrashedCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error)
	CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.ApiResponseCard, error)
	UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.ApiResponseCard, error)
	TrashedCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error)
	RestoreCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error)
	DeleteCardPermanent(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDelete, error)

	RestoreAllCard(context.Context, *emptypb.Empty) (*pb.ApiResponseCardAll, error)
	DeleteAllCardPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseCardAll, error)
}

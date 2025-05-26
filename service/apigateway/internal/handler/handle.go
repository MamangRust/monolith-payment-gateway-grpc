package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type ServiceConnections struct {
	Auth        *grpc.ClientConn
	Role        *grpc.ClientConn
	Card        *grpc.ClientConn
	Merchant    *grpc.ClientConn
	User        *grpc.ClientConn
	Saldo       *grpc.ClientConn
	Topup       *grpc.ClientConn
	Transaction *grpc.ClientConn
	Transfer    *grpc.ClientConn
	Withdraw    *grpc.ClientConn
}

type Deps struct {
	Conn               *grpc.ClientConn
	Kafka              kafka.Kafka
	Token              auth.TokenManager
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	Mapping            apimapper.ResponseApiMapper
	ServiceConnections ServiceConnections
}

func NewHandler(deps Deps) {
	clientAuth := pb.NewAuthServiceClient(deps.ServiceConnections.Auth)
	clientRole := pb.NewRoleServiceClient(deps.ServiceConnections.Role)
	clientCard := pb.NewCardServiceClient(deps.ServiceConnections.Card)
	clientMerchant := pb.NewMerchantServiceClient(deps.ServiceConnections.Merchant)
	clientUser := pb.NewUserServiceClient(deps.ServiceConnections.User)
	clientSaldo := pb.NewSaldoServiceClient(deps.ServiceConnections.Saldo)
	clientTopup := pb.NewTopupServiceClient(deps.ServiceConnections.Topup)
	clientTransaction := pb.NewTransactionServiceClient(deps.ServiceConnections.Transaction)
	clientTransfer := pb.NewTransferServiceClient(deps.ServiceConnections.Transfer)
	clientWithdraw := pb.NewWithdrawServiceClient(deps.ServiceConnections.Withdraw)

	NewHandlerAuth(clientAuth, deps.E, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(clientRole, deps.E, deps.Logger, deps.Mapping.RoleResponseMapper)
	NewHandlerUser(clientUser, deps.E, deps.Logger, deps.Mapping.UserResponseMapper)
	NewHandlerCard(clientCard, deps.E, deps.Logger, deps.Mapping.CardResponseMapper)
	NewHandlerMerchant(clientMerchant, deps.E, deps.Logger, deps.Mapping.MerchantResponseMapper)
	NewHandlerTransaction(clientTransaction, clientMerchant, deps.E, deps.Logger, deps.Mapping.TransactionResponseMapper, deps.Kafka)
	NewHandlerSaldo(clientSaldo, deps.E, deps.Logger, deps.Mapping.SaldoResponseMapper)
	NewHandlerTopup(clientTopup, deps.E, deps.Logger, deps.Mapping.TopupResponseMapper)
	NewHandlerTransfer(clientTransfer, deps.E, deps.Logger, deps.Mapping.TransferResponseMapper)
	NewHandlerWithdraw(clientWithdraw, deps.E, deps.Logger, deps.Mapping.WithdrawResponseMapper)
}

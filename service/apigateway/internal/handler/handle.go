package handler

import (
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
	Kafka              *kafka.Kafka
	Token              auth.TokenManager
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	Mapping            *apimapper.ResponseApiMapper
	ServiceConnections *ServiceConnections
}

func NewHandler(deps *Deps) {
	clientAuth := pb.NewAuthServiceClient(deps.ServiceConnections.Auth)
	clientRole := pb.NewRoleServiceClient(deps.ServiceConnections.Role)
	clientCard := pb.NewCardServiceClient(deps.ServiceConnections.Card)
	clientMerchant := pb.NewMerchantServiceClient(deps.ServiceConnections.Merchant)
	clientMerchantDocument := pb.NewMerchantDocumentServiceClient(deps.ServiceConnections.Merchant)
	clientUser := pb.NewUserServiceClient(deps.ServiceConnections.User)
	clientSaldo := pb.NewSaldoServiceClient(deps.ServiceConnections.Saldo)
	clientTopup := pb.NewTopupServiceClient(deps.ServiceConnections.Topup)
	clientTransaction := pb.NewTransactionServiceClient(deps.ServiceConnections.Transaction)
	clientTransfer := pb.NewTransferServiceClient(deps.ServiceConnections.Transfer)
	clientWithdraw := pb.NewWithdrawServiceClient(deps.ServiceConnections.Withdraw)

	NewHandlerAuth(clientAuth, deps.E, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(clientRole, deps.E, deps.Logger, deps.Mapping.RoleResponseMapper, deps.Kafka)
	NewHandlerUser(clientUser, deps.E, deps.Logger, deps.Mapping.UserResponseMapper)
	NewHandlerCard(clientCard, deps.E, deps.Logger, deps.Mapping.CardResponseMapper)
	NewHandlerMerchant(clientMerchant, deps.E, deps.Logger, deps.Mapping.MerchantResponseMapper)
	NewHandlerMerchantDocument(clientMerchantDocument, deps.E, deps.Logger, deps.Mapping.MerchantDocumentProMapper)
	NewHandlerTransaction(clientTransaction, clientMerchant, deps.E, deps.Logger, deps.Mapping.TransactionResponseMapper, deps.Kafka)
	NewHandlerSaldo(clientSaldo, deps.E, deps.Logger, deps.Mapping.SaldoResponseMapper)
	NewHandlerTopup(clientTopup, deps.E, deps.Logger, deps.Mapping.TopupResponseMapper)
	NewHandlerTransfer(clientTransfer, deps.E, deps.Logger, deps.Mapping.TransferResponseMapper)
	NewHandlerWithdraw(clientWithdraw, deps.E, deps.Logger, deps.Mapping.WithdrawResponseMapper)
}

func parseQueryInt(c echo.Context, key string, defaultValue int) int {
	param := c.QueryParam(key)
	if param == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(param)
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

func parseQueryCard(c echo.Context, logger logger.LoggerInterface) (string, error) {
	cardNumber := strings.TrimSpace(c.QueryParam("card_number"))
	if cardNumber == "" {
		logger.Error("card number is empty")
		return "", card_errors.ErrApiInvalidCardNumber(c)
	}
	return cardNumber, nil
}

func parseQueryMonth(c echo.Context, logger logger.LoggerInterface) (int, error) {
	monthStr := c.QueryParam("month")
	month, err := strconv.Atoi(monthStr)

	if err != nil || month < 1 || month > 12 {
		logger.Error("invalid month", zap.String("month", monthStr))
		return 0, card_errors.ErrApiInvalidMonth(c)
	}

	return month, nil
}

func parseQueryYear(c echo.Context, logger logger.LoggerInterface) (int, error) {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2023 {
		logger.Error("invalid year", zap.String("year", yearStr))
		return 0, card_errors.ErrApiInvalidYear(c)
	}
	return year, nil
}

package handler

import (
	authhandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/auth"
	cardhandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/card"
	merchanthandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/merchant"
	merchantdocumenthandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/merchantdocument"
	rolehandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/role"
	saldohandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/saldo"
	topuphandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/topup"
	transactionhandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/transaction"
	userhandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/user"
	withdrawhandler "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler/withdraw"
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// ServiceConnections holds gRPC connections to external monolith.
type ServiceConnections struct {
	// Auth is the gRPC connection to the authentication service.
	Auth *grpc.ClientConn

	// Role is the gRPC connection to the role management service.
	Role *grpc.ClientConn

	// Card is the gRPC connection to the card service.
	Card *grpc.ClientConn

	// Merchant is the gRPC connection to the merchant service.
	Merchant *grpc.ClientConn

	// User is the gRPC connection to the user management service.
	User *grpc.ClientConn

	// Saldo is the gRPC connection to the saldo/balance service.
	Saldo *grpc.ClientConn

	// Topup is the gRPC connection to the top-up service.
	Topup *grpc.ClientConn

	// Transaction is the gRPC connection to the transaction service.
	Transaction *grpc.ClientConn

	// Transfer is the gRPC connection to the fund transfer service.
	Transfer *grpc.ClientConn

	// Withdraw is the gRPC connection to the withdrawal service.
	Withdraw *grpc.ClientConn
}

// Deps holds dependencies required to initialize HTTP API handlers.
type Deps struct {
	// Kafka is the Kafka instance used for producing and consuming messages.
	Kafka *kafka.Kafka

	// Token is responsible for creating and verifying authentication tokens.
	Token auth.TokenManager

	// E is the Echo instance used for HTTP routing.
	E *echo.Echo

	Mencache mencache.CacheApiGateway

	// Logger is the logging instance for structured and leveled logs.
	Logger logger.LoggerInterface

	// ServiceConnections holds all gRPC connections to other monolith.
	ServiceConnections *ServiceConnections
}

// NewHandler sets up all the handlers for the API Gateway.
// It takes a pointer to a Deps struct, which contains all the dependencies
// required to set up the handlers.
func NewHandler(deps *Deps) {
	authhandler.RegisterAuthHandler(&authhandler.DepsAuth{
		Client: deps.ServiceConnections.Auth,
		E:      deps.E,
		Logger: deps.Logger,
	})

	cardhandler.RegisterCardHandler(&cardhandler.DepsCard{
		Client: deps.ServiceConnections.Card,
		E:      deps.E,
		Logger: deps.Logger,
	})

	merchanthandler.RegisterMerchantHandler(&merchanthandler.DepsMerchant{
		Client: deps.ServiceConnections.Merchant,
		E:      deps.E,
		Logger: deps.Logger,
	})

	merchantdocumenthandler.RegisterMerchantDocumentHandler(&merchantdocumenthandler.DepsMerchantDocument{
		Client: deps.ServiceConnections.Merchant,
		E:      deps.E,
		Logger: deps.Logger,
	})

	rolehandler.RegisterRoleHandler(&rolehandler.DepsRole{
		Kafka:  deps.Kafka,
		Client: deps.ServiceConnections.Role,
		E:      deps.E,
		Logger: deps.Logger,
		Cache:  deps.Mencache,
	})

	saldohandler.RegisterSaldoHandler(&saldohandler.DepsSaldo{
		Client: deps.ServiceConnections.Saldo,
		E:      deps.E,
		Logger: deps.Logger,
	})

	topuphandler.RegisterTopupHandler(&topuphandler.DepsTopup{
		Client: deps.ServiceConnections.Topup,
		E:      deps.E,
		Logger: deps.Logger,
	})

	transactionhandler.RegisterTransactionHandler(&transactionhandler.DepsTransaction{
		Kafka:  deps.Kafka,
		Client: deps.ServiceConnections.Transaction,
		E:      deps.E,
		Logger: deps.Logger,
		Cache:  deps.Mencache,
	})

	userhandler.RegisterUserHandler(&userhandler.DepsUser{
		Client: deps.ServiceConnections.User,
		E:      deps.E,
		Logger: deps.Logger,
	})

	withdrawhandler.RegisterWithdrawHandler(&withdrawhandler.DepsWithdraw{
		Client: deps.ServiceConnections.Withdraw,
		E:      deps.E,
		Logger: deps.Logger,
	})
}

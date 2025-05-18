package app

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	_ "github.com/MamangRust/payment-gateway-monolith-grpc/service/apigateway/docs"
	"github.com/MamangRust/payment-gateway-monolith-grpc/service/apigateway/internal/handler"
	"github.com/MamangRust/payment-gateway-monolith-grpc/service/apigateway/internal/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceAddresses struct {
	Auth        string
	Role        string
	Card        string
	Merchant    string
	User        string
	Saldo       string
	Topup       string
	Transaction string
	Transfer    string
	Withdraw    string
}

// @title PaymentGateway gRPC
// @version 1.0
// @description gRPC based Payment Gateway service

// @host localhost:5000
// @BasePath /api/

// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func RunClient() {
	addresses := ServiceAddresses{
		Auth:        *flag.String("auth-addr", getEnvOrDefault("GRPC_AUTH_ADDR", "localhost:50051"), "Auth service address"),
		Role:        *flag.String("role-addr", getEnvOrDefault("GRPC_ROLE_ADDR", "localhost:50052"), "Role service address"),
		Card:        *flag.String("card-addr", getEnvOrDefault("GRPC_CARD_ADDR", "localhost:50053"), "Card service address"),
		Merchant:    *flag.String("merchant-addr", getEnvOrDefault("GRPC_MERCHANT_ADDR", "localhost:50054"), "Merchant service address"),
		User:        *flag.String("user-addr", getEnvOrDefault("GRPC_USER_ADDR", "localhost:50055"), "User service address"),
		Saldo:       *flag.String("saldo-addr", getEnvOrDefault("GRPC_SALDO_ADDR", "localhost:50056"), "Saldo service address"),
		Topup:       *flag.String("topup-addr", getEnvOrDefault("GRPC_TOPUP_ADDR", "localhost:50057"), "Topup service address"),
		Transaction: *flag.String("transaction-addr", getEnvOrDefault("GRPC_TRANSACTION_ADDR", "localhost:50058"), "Transaction service address"),
		Transfer:    *flag.String("transfer-addr", getEnvOrDefault("GRPC_TRANSFER_ADDR", "localhost:50059"), "Transfer service address"),
		Withdraw:    *flag.String("withdraw-addr", getEnvOrDefault("GRPC_WITHDRAW_ADDR", "localhost:50060"), "Withdraw service address"),
	}

	flag.Parse()

	logger, err := logger.NewLogger()
	limiter := middlewares.NewRateLimiter(20, 50)

	if err != nil {
		logger.Fatal("Failed to create logger", zap.Error(err))
	}

	err = dotenv.Viper()
	if err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	connections, err := createServiceConnections(addresses, logger)
	if err != nil {
		logger.Fatal("Failed to connect to one or more services", zap.Error(err))
	}

	e := echo.New()

	e.Use(limiter.Limit)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1420", "http://localhost:33451"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-API-Key",
		},
		AllowCredentials: true,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	middlewares.WebSecurityConfig(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	token, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		logger.Fatal("Failed to create token manager", zap.Error(err))
	}

	mapping := apimapper.NewResponseApiMapper()

	depsHandler := handler.Deps{
		Conn:               connections.Auth,
		Token:              token,
		E:                  e,
		Logger:             logger,
		Mapping:            *mapping,
		ServiceConnections: handler.ServiceConnections(connections),
	}

	handler.NewHandler(depsHandler)

	go func() {
		if err := e.Start(":5000"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Server.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	for name, conn := range map[string]*grpc.ClientConn{
		"Auth":        connections.Auth,
		"Role":        connections.Role,
		"Card":        connections.Card,
		"Merchant":    connections.Merchant,
		"User":        connections.User,
		"Saldo":       connections.Saldo,
		"Topup":       connections.Topup,
		"Transaction": connections.Transaction,
		"Transfer":    connections.Transfer,
		"Withdraw":    connections.Withdraw,
	} {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.Error(fmt.Sprintf("Failed to close %s connection", name), zap.Error(err))
			}
		}
	}
}

func createServiceConnections(addresses ServiceAddresses, logger logger.LoggerInterface) (handler.ServiceConnections, error) {
	var connections handler.ServiceConnections
	var err error

	connections.Auth, err = createConnection(addresses.Auth, "Auth", logger)
	if err != nil {
		return connections, err
	}

	connections.Role, err = createConnection(addresses.Role, "Role", logger)
	if err != nil {
		return connections, err
	}

	connections.Card, err = createConnection(addresses.Card, "Card", logger)
	if err != nil {
		return connections, err
	}

	connections.Merchant, err = createConnection(addresses.Merchant, "Merchant", logger)
	if err != nil {
		return connections, err
	}

	connections.User, err = createConnection(addresses.User, "User", logger)
	if err != nil {
		return connections, err
	}

	connections.Saldo, err = createConnection(addresses.Saldo, "Saldo", logger)
	if err != nil {
		return connections, err
	}

	connections.Topup, err = createConnection(addresses.Topup, "Topup", logger)
	if err != nil {
		return connections, err
	}

	connections.Transaction, err = createConnection(addresses.Transaction, "Transaction", logger)
	if err != nil {
		return connections, err
	}

	connections.Transfer, err = createConnection(addresses.Transfer, "Transfer", logger)
	if err != nil {
		return connections, err
	}

	connections.Withdraw, err = createConnection(addresses.Withdraw, "Withdraw", logger)
	if err != nil {
		return connections, err
	}

	return connections, nil
}

func createConnection(address, serviceName string, logger logger.LoggerInterface) (*grpc.ClientConn, error) {
	logger.Info(fmt.Sprintf("Connecting to %s service at %s", serviceName, address))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to %s service", serviceName), zap.Error(err))
		return nil, err
	}
	return conn, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

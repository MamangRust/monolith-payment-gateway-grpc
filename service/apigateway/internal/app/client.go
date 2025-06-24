package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/MamangRust/monolith-payment-gateway-apigateway/docs"
	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/middlewares"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
type Client struct {
	App    *echo.Echo
	Logger logger.LoggerInterface
}

func (c *Client) Shutdown(ctx context.Context) error {
	return c.App.Shutdown(ctx)
}

func RunClient() (*Client, func(), error) {
	flag.Parse()

	addresses := loadServiceAddresses()

	log, err := logger.NewLogger("apigateway")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	log.Debug("Loading environment variables")
	if err := dotenv.Viper(); err != nil {
		log.Fatal("Failed to load .env file", zap.Error(err))
	}

	log.Debug("Creating gRPC connections...")
	conns, err := createServiceConnections(addresses, log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect services: %w", err)
	}

	e := setupEcho()

	token, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		log.Fatal("Failed to create token manager", zap.Error(err))
	}

	ctx := context.Background()
	shutdownTracer, err := otel_pkg.InitTracerProvider("apigateway", ctx)
	if err != nil {
		log.Fatal("Failed to initialize tracer provider", zap.Error(err))
	}

	myKafka := kafka.NewKafka(log, []string{os.Getenv("KAFKA_BROKERS")})
	mapping := apimapper.NewResponseApiMapper()

	deps := &handler.Deps{
		Kafka:              myKafka,
		Token:              token,
		E:                  e,
		Logger:             log,
		Mapping:            mapping,
		ServiceConnections: handler.ServiceConnections(conns),
	}

	handler.NewHandler(deps)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Info("Starting API Gateway server on :5000")
		if err := e.Start(":5000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Echo server error", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		log.Info("Starting Prometheus metrics server on :8091")
		if err := http.ListenAndServe(":8091", promhttp.Handler()); err != nil {
			log.Fatal("Metrics server error", zap.Error(err))
		}
	}()

	shutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Info("Shutting down API Gateway...")
		if err := e.Shutdown(ctx); err != nil {
			log.Error("Echo shutdown failed", zap.Error(err))
		}

		closeConnections(conns, log)

		if shutdownTracer != nil {
			if err := shutdownTracer(context.Background()); err != nil {
				log.Error("Tracer shutdown failed", zap.Error(err))
			}
		}
	}

	return &Client{App: e, Logger: log}, shutdown, nil
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

func loadServiceAddresses() ServiceAddresses {
	return ServiceAddresses{
		Auth:        getEnvOrDefault("GRPC_AUTH_ADDR", "localhost:50051"),
		Role:        getEnvOrDefault("GRPC_ROLE_ADDR", "localhost:50052"),
		Card:        getEnvOrDefault("GRPC_CARD_ADDR", "localhost:50053"),
		Merchant:    getEnvOrDefault("GRPC_MERCHANT_ADDR", "localhost:50054"),
		User:        getEnvOrDefault("GRPC_USER_ADDR", "localhost:50055"),
		Saldo:       getEnvOrDefault("GRPC_SALDO_ADDR", "localhost:50056"),
		Topup:       getEnvOrDefault("GRPC_TOPUP_ADDR", "localhost:50057"),
		Transaction: getEnvOrDefault("GRPC_TRANSACTION_ADDR", "localhost:50058"),
		Transfer:    getEnvOrDefault("GRPC_TRANSFER_ADDR", "localhost:50059"),
		Withdraw:    getEnvOrDefault("GRPC_WITHDRAW_ADDR", "localhost:50060"),
	}
}

func setupEcho() *echo.Echo {
	e := echo.New()

	limiter := middlewares.NewRateLimiter(20, 50)
	e.Use(limiter.Limit, middleware.Recover(), middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1420", "http://localhost:33451"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-API-Key"},
		AllowCredentials: true,
	}))

	middlewares.WebSecurityConfig(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

func closeConnections(conns handler.ServiceConnections, log logger.LoggerInterface) {
	for name, conn := range map[string]*grpc.ClientConn{
		"Auth":        conns.Auth,
		"Role":        conns.Role,
		"Card":        conns.Card,
		"Merchant":    conns.Merchant,
		"User":        conns.User,
		"Saldo":       conns.Saldo,
		"Topup":       conns.Topup,
		"Transaction": conns.Transaction,
		"Transfer":    conns.Transfer,
		"Withdraw":    conns.Withdraw,
	} {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Error(fmt.Sprintf("Failed to close %s connection", name), zap.Error(err))
			}
		}
	}
}

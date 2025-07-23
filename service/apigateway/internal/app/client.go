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
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	redisclient "github.com/MamangRust/monolith-payment-gateway-pkg/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceAddresses holds the addresses of all the monolith.
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

// loadServiceAddresses loads the addresses of all the monolith.
//
// It gets the addresses from environment variables with default values.
//
// The environment variables are:
//   - GRPC_AUTH_ADDR
//   - GRPC_ROLE_ADDR
//   - GRPC_CARD_ADDR
//   - GRPC_MERCHANT_ADDR
//   - GRPC_USER_ADDR
//   - GRPC_SALDO_ADDR
//   - GRPC_TOPUP_ADDR
//   - GRPC_TRANSACTION_ADDR
//   - GRPC_TRANSFER_ADDR
//   - GRPC_WITHDRAW_ADDR
func loadServiceAddresses() *ServiceAddresses {
	return &ServiceAddresses{
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

// createServiceConnections creates a set of gRPC connections to the monolith.
//
// It takes a pointer to a ServiceAddresses struct and a logger interface as arguments.
// The ServiceAddresses struct contains the addresses of all the monolith.
//
// The function returns a pointer to a ServiceConnections struct which contains the gRPC connections
// to the monolith or an error if any of the connections fail.
//
// The connections are created using the createConnection function which returns a gRPC connection
// to the specified microservice or an error if the connection fails.
//
// The function will return an error if any of the connections fail.
func createServiceConnections(addresses *ServiceAddresses, logger logger.LoggerInterface) (*handler.ServiceConnections, error) {
	var connections handler.ServiceConnections

	conns := map[string]*string{
		"Auth":        &addresses.Auth,
		"Role":        &addresses.Role,
		"Card":        &addresses.Card,
		"Merchant":    &addresses.Merchant,
		"User":        &addresses.User,
		"Saldo":       &addresses.Saldo,
		"Topup":       &addresses.Topup,
		"Transaction": &addresses.Transaction,
		"Transfer":    &addresses.Transfer,
		"Withdraw":    &addresses.Withdraw,
	}

	for name, addr := range conns {
		conn, err := createConnection(*addr, name, logger)
		if err != nil {
			return nil, err
		}

		switch name {
		case "Auth":
			connections.Auth = conn
		case "Role":
			connections.Role = conn
		case "Card":
			connections.Card = conn
		case "Merchant":
			connections.Merchant = conn
		case "User":
			connections.User = conn
		case "Saldo":
			connections.Saldo = conn
		case "Topup":
			connections.Topup = conn
		case "Transaction":
			connections.Transaction = conn
		case "Transfer":
			connections.Transfer = conn
		case "Withdraw":
			connections.Withdraw = conn
		}
	}

	return &connections, nil
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

// Shutdown the server.
//
// This method is idempotent and will not return an error if the server is already shutdown.
//
// The context is used to receive any error encountered that caused the shutdown.
//
// The client will not be usable after this method is called.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.App.Shutdown(ctx)
}

// RunClient initializes and runs the API Gateway client.
//
// It sets up environment variables, creates a logger, establishes gRPC connections to monolith,
// initializes an authentication token manager, and sets up tracing using OpenTelemetry.
// The function also starts the Echo server and a Prometheus metrics server concurrently.
//
// Returns a Client instance, a shutdown function to gracefully stop services, and an error
// if initialization fails.
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

	myredis := redisclient.NewRedisClient(&redisclient.Config{
		Host:         viper.GetString("REDIS_HOST"),
		Port:         viper.GetString("REDIS_PORT"),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB_APIGATEWAY"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := myredis.Client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to ping redis", zap.Error(err))
	}

	mencache := mencache.NewCacheApiGateway(&mencache.Deps{
		Redis:  myredis.Client,
		Logger: log,
	})

	deps := &handler.Deps{
		Kafka:              myKafka,
		Token:              token,
		E:                  e,
		Logger:             log,
		ServiceConnections: conns,
		Mencache:           mencache,
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

// createConnection creates a gRPC client connection to a service given its address and
// service name. It logs the connection attempt and any errors that occur.
//
// The service name is used for logging purposes only.
//
// The function uses the insecure gRPC transport credentials, as the services are assumed
// to be running in the same network and no encryption is needed.
func createConnection(address, serviceName string, logger logger.LoggerInterface) (*grpc.ClientConn, error) {
	logger.Info(fmt.Sprintf("Connecting to %s service at %s", serviceName, address))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to %s service", serviceName), zap.Error(err))
		return nil, err
	}
	return conn, nil
}

// closeConnections gracefully closes all gRPC connections contained within the
// ServiceConnections struct. It logs any errors encountered during the closure
// of each connection using the provided logger.
//
// Parameters:
//   - conns: A pointer to a ServiceConnections struct holding the gRPC connections
//     to be closed.
//   - log: An instance of LoggerInterface used to log errors if any connection fails
//     to close.
//
// The function iterates over each connection in the ServiceConnections struct.
// If a connection is not nil, it attempts to close it and logs any errors that
// occur during the closure process.
func closeConnections(conns *handler.ServiceConnections, log logger.LoggerInterface) {
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

// setupEcho sets up the Echo web framework with the necessary middleware and routes.
//
// It returns a pointer to the Echo instance.
//
// The middleware used are:
//   - middleware.Logger(): logs each request
//   - middleware.Recover(): recovers from panics and logs the error
//   - limiter.Limit: limits the number of requests per second to 20 and the number of requests per minute to 50
//   - middleware.CORSWithConfig: enables CORS with the specified configuration
//   - middlewares.WebSecurityConfig: sets up the web security configuration for the Echo instance
//   - echoSwagger.WrapHandler: sets up the Swagger UI handler
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

// getEnvOrDefault retrieves the value of an environment variable and returns a default
// value if it is not set.
func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

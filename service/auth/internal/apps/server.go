package apps

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/errorhandler"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/middleware"
	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/service"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/database"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	redisclient "github.com/MamangRust/monolith-payment-gateway-pkg/redis"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	port int
)

// init initializes the server port for the gRPC authentication service.
// It retrieves the port number from the environment configuration using Viper.
// If the port is not specified, it defaults to 50051.
// The port can also be overridden via a command-line flag.
func init() {
	port = viper.GetInt("GRPC_AUTH_PORT")
	if port == 0 {
		port = 50051
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

// Server represents the gRPC server for the authentication service.
type Server struct {
	Logger       logger.LoggerInterface
	DB           *db.Queries
	TokenManager *auth.Manager
	Services     *service.Service
	Handlers     *handler.Handler
	Ctx          context.Context
}

// NewServer returns a new instance of the Server struct, along with a function
// to be used to shut down the OpenTelemetry tracer provider and an error.
// The function will be used to shut down the OpenTelemetry tracer provider
// when the server is stopped.
//
// The function initializes the logger, configuration, database connection,
// token manager, repositories, services, and handlers, and returns a new
// instance of the Server struct.
//
// It also returns an error if any of the initialization steps fail.
func NewServer(ctx context.Context) (*Server, func(context.Context) error, error) {
	flag.Parse()

	logger, err := logger.NewLogger("auth-service")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize jlogger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))

		return nil, nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		logger.Fatal("Failed to create token manager", zap.Error(err))

		return nil, nil, err
	}

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))

		return nil, nil, err
	}

	DB := db.New(conn)

	hash := hash.NewHashingPassword()

	depsRepo := &repository.Deps{
		DB: DB,
	}
	repositories := repository.NewRepositories(depsRepo)

	kafka := kafka.NewKafka(logger, []string{viper.GetString("KAFKA_BROKERS")})

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("auth-service", ctx)

	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
	}

	myredis := redisclient.NewRedisClient(&redisclient.Config{
		Host:         viper.GetString("REDIS_HOST"),
		Port:         viper.GetString("REDIS_PORT"),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB_AUTH"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := myredis.Client.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to ping redis", zap.Error(err))
	}

	mencache := mencache.NewMencache(&mencache.Deps{
		Redis:  myredis.Client,
		Logger: logger,
	})

	errorhandler := errorhandler.NewErrorHandler(logger)

	services := service.NewService(&service.Deps{
		ErrorHandler: errorhandler,
		Mencache:     mencache,
		Repositories: repositories,
		Token:        tokenManager,
		Hash:         hash,
		Logger:       logger,
		Kafka:        kafka,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Logger:  logger,
		Service: services,
	})

	return &Server{
		Logger:       logger,
		DB:           DB,
		TokenManager: tokenManager,
		Services:     services,
		Handlers:     handlers,
		Ctx:          ctx,
	}, shutdownTracerProvider, nil
}

// Run starts the gRPC and metrics servers for the auth service.
// It sets up network listeners for the gRPC server and metrics server using ports
// specified in the environment configuration. The function initializes and starts
// a gRPC server with OpenTelemetry instrumentation and registers the auth service handlers.
// Additionally, it creates an HTTP server to serve Prometheus metrics for monitoring.
// The function runs both servers concurrently and waits for them to finish using a wait group.
// If any server encounters an error during execution, the error is logged, and the application
// terminates with a fatal error.
func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_AUTH_ADDR"))
	metricsLis, err := net.Listen("tcp", metricsAddr)
	if err != nil {
		s.Logger.Fatal("failed to listen on", zap.Error(err))
	}

	if err != nil {
		log.Fatalf("Failed to listen for metrics: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			),
		),
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryMiddleware(s.Logger),
			middleware.ContextMiddleware(60*time.Second, s.Logger),
		),
	)

	pb.RegisterAuthServiceServer(grpcServer, s.Handlers.Auth)

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	s.Logger.Info(fmt.Sprintf("Server running on port %d", port))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.Logger.Info("Metrics server listening on :8081")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			s.Logger.Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		s.Logger.Debug("gRPC server listening on port", zap.Int("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	wg.Wait()
}

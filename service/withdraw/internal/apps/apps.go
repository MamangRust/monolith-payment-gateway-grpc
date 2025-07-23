package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/database"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	redisclient "github.com/MamangRust/monolith-payment-gateway-pkg/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/middleware"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
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

// init initializes the gRPC server port for the withdraw service.
// It retrieves the port number from the environment configuration using Viper.
// If the port is not specified, it defaults to 50060.
// The port can also be overridden via a command-line flag.
func init() {
	port = viper.GetInt("GRPC_WITHDRAW_PORT")
	if port == 0 {
		port = 50060
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

// Server represents the gRPC server for the withdraw service.
type Server struct {
	Logger   logger.LoggerInterface
	DB       *db.Queries
	Services service.Service
	Handlers handler.Handler
	Ctx      context.Context
}

// NewServer initializes and returns a new instance of the Server struct for the withdraw service,
// along with a shutdown function for the OpenTelemetry tracer provider and an error.
// It sets up the logger, loads environment configurations using Viper, establishes a database connection,
// initializes the Redis client, tracer provider, Kafka client, cache, error handler, services, and handlers.
// If any of the initialization steps fail, it returns an error.
func NewServer(ctx context.Context) (*Server, func(context.Context) error, error) {
	logger, err := logger.NewLogger("withdraw-service")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
		return nil, nil, err
	}
	flag.Parse()

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
		return nil, nil, err
	}
	DB := db.New(conn)

	repositories := repository.NewRepositories(DB)

	kafka := kafka.NewKafka(logger, []string{viper.GetString("KAFKA_BROKERS")})

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("withdraw-service", ctx)
	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
		return nil, nil, err
	}

	myredis := redisclient.NewRedisClient(&redisclient.Config{
		Host:         viper.GetString("REDIS_HOST"),
		Port:         viper.GetString("REDIS_PORT"),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB_WITHDRAW"),
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
		Mencache:     mencache,
		ErrorHander:  errorhandler,
		Kafka:        kafka,
		Repositories: repositories,
		Logger:       logger,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  logger,
	})

	return &Server{
		Logger:   logger,
		DB:       DB,
		Services: services,
		Handlers: handlers,
		Ctx:      ctx,
	}, shutdownTracerProvider, nil
}

// Run starts the gRPC and metrics servers for the withdraw service.
// It sets up network listeners for the gRPC server and metrics server using ports
// specified in the environment configuration. The function initializes and starts
// a gRPC server with OpenTelemetry instrumentation and registers the withdraw service handlers.
// Additionally, it creates an HTTP server to serve Prometheus metrics for monitoring.
// The function runs both servers concurrently and waits for them to finish using a wait group.
// If any server encounters an error during execution, the error is logged, and the application
// terminates with a fatal error.
func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}
	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_WITHDRAW_ADDR"))
	metricsLis, err := net.Listen("tcp", metricsAddr)
	if err != nil {
		s.Logger.Fatal("failed to listen on", zap.Error(err))
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

	s.RegisterHandleGrpc(grpcServer, s.Handlers)

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	s.Logger.Info(fmt.Sprintf("Server running on port %d", port))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.Logger.Info("Metrics server listening on :8090")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			s.Logger.Fatal("Metrics server error", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		s.Logger.Info("gRPC server listening on :50060")
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	wg.Wait()
}

func (s *Server) RegisterHandleGrpc(grpcServer *grpc.Server, handler handler.Handler) {
	pb.RegisterWithdrawQueryServiceServer(grpcServer, handler)
	pb.RegisterWithdrawCommandServiceServer(grpcServer, handler)
	pb.RegisterWithdrawStatsAmountServiceServer(grpcServer, handler)
	pb.RegisterWithdrawStatsStatusServer(grpcServer, handler)
}

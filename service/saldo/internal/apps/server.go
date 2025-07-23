package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/database"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	redisclient "github.com/MamangRust/monolith-payment-gateway-pkg/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/handler"
	myhandlerkafka "github.com/MamangRust/monolith-payment-gateway-saldo/internal/kafka"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/middleware"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
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

// init initializes the gRPC server port for the saldo service.
// It retrieves the port number from the environment configuration using Viper.
// If the port is not specified, it defaults to 50056.
// The port can also be overridden via a command-line flag.
func init() {
	port = viper.GetInt("GRPC_SALDO_PORT")
	if port == 0 {
		port = 50056
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

// Server represents the gRPC server for the saldo service.
type Server struct {
	Logger   logger.LoggerInterface
	DB       *db.Queries
	Services service.Service
	Handlers handler.Handler
	Ctx      context.Context
}

// NewServer creates a new instance of Server, which is the gRPC server for the saldo service.
// It initializes the logger, database connection, OpenTelemetry tracer provider, Redis connection,
// and Kafka consumer. It also initializes the service and handler for the saldo service.
// The function returns the Server instance, a shutdown function for the OpenTelemetry tracer provider,
// and an error if any of the initialization steps fail.
func NewServer(ctx context.Context) (*Server, func(context.Context) error, error) {
	flag.Parse()

	logger, err := logger.NewLogger("saldo-service")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
		return nil, nil, err
	}

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
		return nil, nil, err
	}
	DB := db.New(conn)

	repositories := repository.NewRepositories(DB)

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("saldo-service", ctx)
	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
		return nil, nil, err
	}

	myredis := redisclient.NewRedisClient(&redisclient.Config{
		Host:         viper.GetString("REDIS_HOST"),
		Port:         viper.GetString("REDIS_PORT"),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB_SALDO"),
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
		Ctx:    ctx,
		Redis:  myredis.Client,
		Logger: logger,
	})

	errorhandler := errorhandler.NewErrorHandler(logger)

	services := service.NewService(&service.Deps{
		Mencache:     mencache,
		ErrorHandler: errorhandler,
		Repositories: repositories,
		Logger:       logger,
	})

	myKafka := kafka.NewKafka(logger, []string{viper.GetString("KAFKA_BROKERS")})

	handler_kafka_saldo := myhandlerkafka.NewSaldoKafkaHandler(services, logger)

	err = myKafka.StartConsumers([]string{
		"saldo-service-topic-create-saldo",
	}, "saldo-service-group", handler_kafka_saldo)

	if err != nil {
		logger.Fatal("Failed to start consumers", zap.Error(err))
		return nil, nil, err
	}

	handlers := handler.NewHandler(&handler.Deps{
		Logger:  logger,
		Service: services,
	})

	return &Server{
		Logger:   logger,
		DB:       DB,
		Services: services,
		Handlers: handlers,
		Ctx:      ctx,
	}, shutdownTracerProvider, nil
}

// Run starts the gRPC server and a metrics server. It serves the gRPC server
// on port 50056 and the metrics server on port 8086. It blocks until the
// server is stopped.
func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}
	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_SALDO_ADDR"))
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
		s.Logger.Info("Metrics server listening on :8086")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			s.Logger.Fatal("Metrics server error", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		s.Logger.Info("gRPC server listening on :50056")
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("Failed to serve gRPC server", zap.Error(err))
		}
	}()

	wg.Wait()
}

func (s *Server) RegisterHandleGrpc(grpcServer *grpc.Server, handler handler.Handler) {
	pb.RegisterSaldoQueryServiceServer(grpcServer, handler)
	pb.RegisterSaldoCommandServiceServer(grpcServer, handler)
	pb.RegisterSaldoStatsBalanceServiceServer(grpcServer, handler)
	pb.RegisterSaldoStatsTotalBalanceServer(grpcServer, handler)
}

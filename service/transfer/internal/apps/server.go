package apps

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
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

func init() {
	port = viper.GetInt("GRPC_TRANSFER_ADDR")
	if port == 0 {
		port = 50059
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

type Server struct {
	Logger   logger.LoggerInterface
	DB       *db.Queries
	Services *service.Service
	Handlers *handler.Handler
	Ctx      context.Context
}

func NewServer() (*Server, error) {
	logger, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}
	flag.Parse()

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	DB := db.New(conn)

	ctx := context.Background()

	mapperRecord := recordmapper.NewRecordMapper()

	depsRepo := repository.Deps{
		DB:           DB,
		Ctx:          ctx,
		MapperRecord: mapperRecord,
	}

	repositories := repository.NewRepositories(depsRepo)
	kafka := kafka.NewKafka(logger, []string{viper.GetString("KAFKA_BROKERS")})

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("Transfer-service", ctx)
	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			logger.Fatal("Failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	services := service.NewService(service.Deps{
		Kafka:        *kafka,
		Ctx:          ctx,
		Repositories: repositories,
		Logger:       logger,
	})

	handlers := handler.NewHandler(handler.Deps{
		Service: *services,
	})

	return &Server{
		Logger:   logger,
		DB:       DB,
		Services: services,
		Handlers: handlers,
		Ctx:      ctx,
	}, nil
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_TRANSFER_ADDR"))
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
	)

	pb.RegisterTransferServiceServer(grpcServer, s.Handlers.Transfer)

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	s.Logger.Info(fmt.Sprintf("Server running on port %d", port))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("Metrics server listening on :8089")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			s.Logger.Fatal("Metrics server error", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		log.Println("gRPC server listening on :50059")
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	wg.Wait()
}

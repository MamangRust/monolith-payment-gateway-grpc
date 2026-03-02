package apps

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/middleware"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/grafana/pyroscope-go"
	"github.com/redis/go-redis/v9"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/service"

	pb "github.com/MamangRust/monolith-payment-gateway-pb"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/database"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	"github.com/MamangRust/monolith-payment-gateway-pkg/resilience"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort              = 50051
	defaultMaxConcurrentConn = 1024
	defaultWindowSize        = 16 * 1024 * 1024
	defaultKeepaliveTime     = 20 * time.Second
	defaultKeepaliveTimeout  = 5 * time.Second
	defaultMinKeepaliveTime  = 5 * time.Second

	monitoringInterval     = 30 * time.Second
	cleanupInterval        = 120 * time.Second
	cacheRefCountThreshold = 500

	shutdownTimeout = 30 * time.Second

	redisDialTimeout  = 5 * time.Second
	redisReadTimeout  = 3 * time.Second
	redisWriteTimeout = 3 * time.Second
	redisPoolSize     = 10
	redisMinIdleConns = 3
)

var (
	port = flag.Int("port", defaultPort, "gRPC server port")
)

type Server struct {
	Logger       logger.LoggerInterface
	DB           *db.Queries
	TokenManager *auth.Manager
	Services     *service.Service
	Handlers     *handler.Handler
	Ctx          context.Context
	Cancel       context.CancelFunc
	CacheStore   *cache.CacheStore
	Redis        *redis.Client
	Telemetry    *otel_pkg.Telemetry
}

type Config struct {
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	Environment    string `mapstructure:"environment"`
	OtelEndpoint   string `mapstructure:"otel_endpoint"`
}

func NewServer(cfg *Config) (*Server, error) {
	flag.Parse()

	if err := dotenv.Viper(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	if err := initPyroscope(); err != nil {
		log.Fatal("Failed to initialize pyroscope:", err)
	}

	cfg, err := InitConfig()

	if err != nil {
		log.Fatal(err)
	}

	telemetry, err := initTelemetry(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	cacheMetrics, err := observability.NewCacheMetrics("cache")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache metrics: %w", err)
	}

	logger, err := logger.NewLogger(cfg.ServiceName, telemetry.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	dbConn, err := database.NewClient(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	queries := db.New(dbConn)

	ctx, cancel := context.WithCancel(context.Background())

	redisClient, err := initRedisServer(ctx, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	cacheStore := cache.NewCacheStore(redisClient, logger, cacheMetrics)

	hasher := hash.NewHashingPassword()

	repositories := repository.NewRepositories(queries)

	kafka := kafka.NewKafka(logger, []string{viper.GetString("KAFKA_BROKERS")})

	services := service.NewService(&service.Deps{
		Cache:        cacheStore,
		Repositories: repositories,
		Token:        tokenManager,
		Hash:         hasher,
		Logger:       logger,
		Kafka:        kafka,
	})

	handlers := handler.NewHandler(&handler.Deps{Service: services, Logger: logger})

	server := &Server{
		Logger:       logger,
		DB:           queries,
		TokenManager: tokenManager,
		Services:     services,
		Handlers:     handlers,
		Ctx:          ctx,
		Cancel:       cancel,
		CacheStore:   cacheStore,
		Redis:        redisClient,
		Telemetry:    telemetry,
	}

	logger.Info("Server initialized successfully",
		zap.String("service", cfg.ServiceName),
		zap.String("version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
	)

	return server, nil
}

func (s *Server) Run() error {
	defer s.Cleanup()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", defaultPort, err)
	}

	resilienceManager := s.initResilience()

	grpcServer := s.createGRPCServer(resilienceManager)

	s.registerServices(grpcServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	if getEnv("ENABLE_REFLECTION", "false") == "true" {
		reflection.Register(grpcServer)
		s.Logger.Info("gRPC reflection enabled")
	}

	monitoringDone := spawnMonitoringTask(s.Ctx, s.CacheStore)
	cleanupDone := spawnCleanupTask(s.Ctx, s.CacheStore)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	errChan := make(chan error, 1)
	go func() {
		s.Logger.Info("gRPC server starting",
			zap.Int("port", *port),
			zap.String("address", lis.Addr().String()),
		)
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case sig := <-sigChan:
		s.Logger.Info("Received shutdown signal",
			zap.String("signal", sig.String()),
		)
	case err := <-errChan:
		s.Logger.Error("Server error", zap.Error(err))
		return err
	}

	return s.gracefulShutdown(grpcServer, healthServer, monitoringDone, cleanupDone)
}

func (s *Server) initResilience() *middleware.ResilienceInterceptor {
	loadMonitor := resilience.NewLoadMonitor()
	circuitBreaker := resilience.NewCircuitBreaker(100, 10, s.Logger)
	requestLimiter := resilience.NewRequestLimiter(800, s.Logger)

	return middleware.NewResilienceInterceptor(loadMonitor, circuitBreaker, requestLimiter)
}

func (s *Server) createGRPCServer(resilienceManager *middleware.ResilienceInterceptor) *grpc.Server {
	return grpc.NewServer(
		grpc.MaxConcurrentStreams(defaultMaxConcurrentConn),
		grpc.InitialConnWindowSize(defaultWindowSize),
		grpc.InitialWindowSize(defaultWindowSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    defaultKeepaliveTime,
			Timeout: defaultKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             defaultMinKeepaliveTime,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			middleware.PyroscopeUnaryInterceptor(),
			resilienceManager.UnaryInterceptor(),
		),
	)
}

func (s *Server) registerServices(grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, s.Handlers.Auth)

	s.Logger.Info("All gRPC services registered successfully")
}

func (s *Server) gracefulShutdown(
	grpcServer *grpc.Server,
	healthServer *health.Server,
	monitoringDone, cleanupDone <-chan struct{},
) error {
	s.Logger.Info("Starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	s.Cancel()

	tasksDone := make(chan struct{})
	go func() {
		<-monitoringDone
		<-cleanupDone
		close(tasksDone)
	}()

	select {
	case <-tasksDone:
		s.Logger.Info("Background tasks stopped successfully")
	case <-shutdownCtx.Done():
		s.Logger.Warn("Background tasks shutdown timeout, forcing stop")
	}

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.Logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		s.Logger.Warn("Graceful shutdown timeout, forcing stop")
		grpcServer.Stop()
	}

	s.Logger.Info("Graceful shutdown completed")
	return nil
}

func (s *Server) Cleanup() {
	s.Logger.Info("Cleaning up resources...")

	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil {
			s.Logger.Error("Failed to close Redis connection", zap.Error(err))
		} else {
			s.Logger.Info("Redis connection closed")
		}
	}

	if s.Telemetry != nil {
		if err := s.Telemetry.Shutdown(context.Background()); err != nil {
			s.Logger.Error("Failed to shutdown telemetry", zap.Error(err))
		} else {
			s.Logger.Info("Telemetry shutdown successfully")
		}
	}

	s.Logger.Info("Cleanup completed")
}

func initPyroscope() error {
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "auth-service",
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER"),
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
		Tags: map[string]string{
			"service": "auth-service",
			"env":     os.Getenv("ENV"),
			"version": os.Getenv("VERSION"),
		},
	})
	return err
}

func InitConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("service_name", "auth-service")
	v.SetDefault("service_version", "1.0.0")
	v.SetDefault("environment", "production")
	v.SetDefault("otel_endpoint", "otel-collector:4317")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initTelemetry(cfg *Config) (*otel_pkg.Telemetry, error) {
	telemetry := otel_pkg.NewTelemetry(otel_pkg.Config{
		ServiceName:            cfg.ServiceName,
		ServiceVersion:         cfg.ServiceVersion,
		Environment:            cfg.Environment,
		Endpoint:               cfg.OtelEndpoint,
		Insecure:               true,
		EnableRuntimeMetrics:   true,
		RuntimeMetricsInterval: 15 * time.Second,
	})

	if err := telemetry.Init(context.Background()); err != nil {
		return nil, err
	}

	return telemetry, nil
}

func initRedisServer(ctx context.Context, logger logger.LoggerInterface) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_AUTH"), viper.GetString("REDIS_PORT_AUTH")),
		Password:     viper.GetString("REDIS_PASSWORD_AUTH"),
		DB:           viper.GetInt("REDIS_DB_AUTH"),
		DialTimeout:  redisDialTimeout,
		ReadTimeout:  redisReadTimeout,
		WriteTimeout: redisWriteTimeout,
		PoolSize:     redisPoolSize,
		MinIdleConns: redisMinIdleConns,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Info("Redis connection established",
		zap.String("addr", fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_AUTH"), viper.GetString("REDIS_PORT_AUTH"))),
		zap.Int("db", viper.GetInt("REDIS_DB_AUTH")),
	)

	return client, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func spawnMonitoringTask(ctx context.Context, cache *cache.CacheStore) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ticker := time.NewTicker(monitoringInterval)
		defer ticker.Stop()

		cache.Logger.Info("Cache monitoring task started",
			zap.Duration("interval", monitoringInterval),
		)

		for {
			select {
			case <-ctx.Done():
				cache.Logger.Info("Cache monitoring task stopped")
				return
			case <-ticker.C:
				monitorCache(ctx, cache)
			}
		}
	}()

	return done
}

func monitorCache(ctx context.Context, cache *cache.CacheStore) {
	refCount := cache.GetRefCount()

	stats, err := cache.GetStats(ctx)
	if err != nil {
		cache.Logger.Error("Failed to get cache stats", zap.Error(err))
		return
	}

	logLevel := zap.InfoLevel
	if refCount > cacheRefCountThreshold {
		logLevel = zap.WarnLevel
	}

	if ce := cache.Logger.Check(logLevel, "Cache statistics"); ce != nil {
		ce.Write(
			zap.Int64("ref_count", refCount),
			zap.Int64("total_keys", stats.TotalKeys),
			zap.Float64("hit_rate", stats.HitRate),
			zap.String("memory_used", stats.MemoryUsedHuman),
			zap.Bool("high_ref_count", refCount > cacheRefCountThreshold),
		)
	}
}

func spawnCleanupTask(ctx context.Context, cache *cache.CacheStore) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		cache.Logger.Info("Cache cleanup task started",
			zap.Duration("interval", cleanupInterval),
		)

		for {
			select {
			case <-ctx.Done():
				cache.Logger.Info("Cache cleanup task stopped")
				return
			case <-ticker.C:
				cleanupCache(ctx, cache)
			}
		}
	}()

	return done
}

func cleanupCache(ctx context.Context, cache *cache.CacheStore) {
	cache.Logger.Info("Starting periodic cache cleanup")

	statsBefore, err := cache.GetStats(ctx)
	if err != nil {
		cache.Logger.Error("Failed to get cache stats before cleanup", zap.Error(err))
		statsBefore = nil
	}

	scanned, err := cache.ClearExpired(ctx)
	if err != nil {
		cache.Logger.Error("Cache cleanup failed", zap.Error(err))
		return
	}

	statsAfter, err := cache.GetStats(ctx)
	if err != nil {
		cache.Logger.Error("Failed to get cache stats after cleanup", zap.Error(err))
		statsAfter = nil
	}

	logFields := []zap.Field{
		zap.Int64("scanned_keys", scanned),
		zap.Int64("ref_count", cache.GetRefCount()),
	}

	if statsBefore != nil && statsAfter != nil {
		keysRemoved := statsBefore.TotalKeys - statsAfter.TotalKeys
		logFields = append(logFields,
			zap.Int64("keys_before", statsBefore.TotalKeys),
			zap.Int64("keys_after", statsAfter.TotalKeys),
			zap.Int64("keys_removed", keysRemoved),
			zap.String("memory_before", statsBefore.MemoryUsedHuman),
			zap.String("memory_after", statsAfter.MemoryUsedHuman),
		)
	}

	cache.Logger.Info("Cache cleanup completed", logFields...)
}

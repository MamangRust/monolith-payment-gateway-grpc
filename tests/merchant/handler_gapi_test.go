package merchant_test

import (
	"context"
	"net"
	"testing"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-merchant/handler"
	"github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	"github.com/MamangRust/monolith-payment-gateway-merchant/service"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantGapiTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	dbPool          *pgxpool.Pool
	redisClient     *redis.Client
	grpcServer      *grpc.Server
	commandClient   pb.MerchantCommandServiceClient
	queryClient     pb.MerchantQueryServiceClient
	conn            *grpc.ClientConn
	userRepo        user_repo.UserCommandRepository
	userID          int
	merchantID      int
}

func (s *MerchantGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)
	s.userRepo = user_repo.NewUserCommandRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	merchantService := service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	user, err := s.userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Gapi",
		LastName:  "Merchant",
		Email:     "gapi.merchant@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Start gRPC Server
	merchantHandler := handler.NewHandler(merchantService)
	server := grpc.NewServer()
	pb.RegisterMerchantCommandServiceServer(server, merchantHandler)
	pb.RegisterMerchantQueryServiceServer(server, merchantHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.commandClient = pb.NewMerchantCommandServiceClient(conn)
	s.queryClient = pb.NewMerchantQueryServiceClient(conn)
}

func (s *MerchantGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantGapiTestSuite) Test1_CreateMerchant() {
	ctx := context.Background()

	createReq := &pb.CreateMerchantRequest{
		Name:   "Gapi Merchant",
		UserId: int32(s.userID),
	}
	res, err := s.commandClient.CreateMerchant(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Name, res.Data.Name)
	s.merchantID = int(res.Data.Id)
}

func (s *MerchantGapiTestSuite) Test2_FindMerchantById() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	findReq := &pb.FindByIdMerchantRequest{
		MerchantId: int32(s.merchantID),
	}
	found, err := s.queryClient.FindByIdMerchant(ctx, findReq)
	s.NoError(err)
	s.Equal(int32(s.merchantID), found.Data.Id)
}

func TestMerchantGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantGapiTestSuite))
}

package card_test

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-card/handler"
	"github.com/MamangRust/monolith-payment-gateway-card/repository"
	"github.com/MamangRust/monolith-payment-gateway-card/service"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardGapiTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	dbPool       *pgxpool.Pool
	redisClient *redis.Client
	grpcServer   *grpc.Server
	conn         *grpc.ClientConn
	queryClient  pb.CardQueryServiceClient
	cmdClient    pb.CardCommandServiceClient
	userID       int
	cardID       int
}

func (s *CardGapiTestSuite) SetupSuite() {
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
	userRepo := user_repo.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	cardService := service.NewService(&service.Deps{
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
		Kafka:        nil,
	})

	cardHandler := handler.NewHandler(cardService)
	server := grpc.NewServer()
	pb.RegisterCardQueryServiceServer(server, cardHandler)
	pb.RegisterCardCommandServiceServer(server, cardHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.queryClient = pb.NewCardQueryServiceClient(conn)
	s.cmdClient = pb.NewCardCommandServiceClient(conn)

	// Create user
	user, err := userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Gapi",
		LastName:  "Card",
		Email:     "gapi.card@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *CardGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *CardGapiTestSuite) Test1_CreateCard() {
	req := &pb.CreateCardRequest{
		UserId:       int32(s.userID),
		CardType:     "debit",
		ExpireDate:   timestamppb.New(time.Now().AddDate(5, 0, 0)),
		Cvv:          "123",
		CardProvider: "Visa",
	}

	res, err := s.cmdClient.CreateCard(context.Background(), req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	s.cardID = int(res.Data.Id)
}

func (s *CardGapiTestSuite) Test2_FindById() {
	s.Require().NotZero(s.cardID)
	req := &pb.FindByIdCardRequest{
		CardId: int32(s.cardID),
	}

	res, err := s.queryClient.FindByIdCard(context.Background(), req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(int32(s.cardID), res.Data.Id)
}

func (s *CardGapiTestSuite) Test3_UpdateCard() {
	s.Require().NotZero(s.cardID)
	req := &pb.UpdateCardRequest{
		CardId:       int32(s.cardID),
		UserId:       int32(s.userID),
		CardType:     "credit",
		ExpireDate:   timestamppb.New(time.Now().AddDate(6, 0, 0)),
		Cvv:          "456",
		CardProvider: "MasterCard",
	}

	res, err := s.cmdClient.UpdateCard(context.Background(), req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	s.Equal("credit", res.Data.CardType)
}

func (s *CardGapiTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.cardID)
	ctx := context.Background()

	trashRes, err := s.cmdClient.TrashedCard(ctx, &pb.FindByIdCardRequest{CardId: int32(s.cardID)})
	s.NoError(err)
	s.Equal("success", trashRes.Status)

	restoreRes, err := s.cmdClient.RestoreCard(ctx, &pb.FindByIdCardRequest{CardId: int32(s.cardID)})
	s.NoError(err)
	s.Equal("success", restoreRes.Status)
}

func (s *CardGapiTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.cardID)
	ctx := context.Background()

	_, _ = s.cmdClient.TrashedCard(ctx, &pb.FindByIdCardRequest{CardId: int32(s.cardID)})

	delRes, err := s.cmdClient.DeleteCardPermanent(ctx, &pb.FindByIdCardRequest{CardId: int32(s.cardID)})
	s.NoError(err)
	s.Equal("success", delRes.Status)
}

func TestCardGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardGapiTestSuite))
}

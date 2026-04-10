package saldo_test

import (
	"context"
	"net"
	"testing"
	"time"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	card_pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	gapi "github.com/MamangRust/monolith-payment-gateway-saldo/handler"
)

type SaldoGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	commandClient pb.SaldoCommandServiceClient
	queryClient   pb.SaldoQueryServiceClient
	conn        *grpc.ClientConn
	
	userRepo     user_repo.UserCommandRepository
	cardRepo     card_repo.CardCommandRepository
	saldoRepo    saldo_repo.Repositories

	cardNumber string
	saldoID    int32
}

func (s *SaldoGapiTestSuite) SetupSuite() {
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
	saldoRepos := saldo_repo.NewRepositories(queries)
	userRepos := user_repo.NewRepositories(queries)
	cardRepos := card_repo.NewRepositories(queries)

	s.userRepo = userRepos.UserCommand()
	s.cardRepo = cardRepos.CardCommand
	s.saldoRepo = saldoRepos

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	saldoService := service.NewService(&service.Deps{
		Repositories: s.saldoRepo,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User and Card
	user, err := s.userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Saldo", LastName: "Gapi", Email: "saldo.gapi@test.com", Password: "password123",
	})
	s.Require().NoError(err)
	
	card, err := s.cardRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "333", CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber

	// Start gRPC Server
	saldoHandler := gapi.NewHandler(saldoService)
	server := grpc.NewServer()
	pb.RegisterSaldoCommandServiceServer(server, saldoHandler)
	pb.RegisterSaldoQueryServiceServer(server, saldoHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.commandClient = pb.NewSaldoCommandServiceClient(conn)
	s.queryClient = pb.NewSaldoQueryServiceClient(conn)
}

func (s *SaldoGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *SaldoGapiTestSuite) Test1_Create() {
	ctx := context.Background()

	createReq := &pb.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 1000000,
	}
	res, err := s.commandClient.CreateSaldo(ctx, createReq)
	s.NoError(err)
	s.Equal(int32(1000000), res.Data.TotalBalance)

	s.saldoID = res.Data.SaldoId
}

func (s *SaldoGapiTestSuite) Test2_FindByCardNumber() {
	s.Require().NotEmpty(s.cardNumber)
	ctx := context.Background()

	found, err := s.queryClient.FindByCardNumber(ctx, &card_pb.FindByCardNumberRequest{CardNumber: s.cardNumber})
	s.NoError(err)
	s.Equal(int32(1000000), found.Data.TotalBalance)
}

func (s *SaldoGapiTestSuite) Test3_Update() {
	s.Require().NotZero(s.saldoID)
	ctx := context.Background()

	_, err := s.commandClient.UpdateSaldo(ctx, &pb.UpdateSaldoRequest{
		SaldoId:      s.saldoID,
		CardNumber:   s.cardNumber,
		TotalBalance: 1200000,
	})
	s.NoError(err)

	// Verify update
	updated, _ := s.queryClient.FindByCardNumber(ctx, &card_pb.FindByCardNumberRequest{CardNumber: s.cardNumber})
	s.Equal(int32(1200000), updated.Data.TotalBalance)
}

func (s *SaldoGapiTestSuite) Test4_Trashed() {
	s.Require().NotZero(s.saldoID)
	ctx := context.Background()

	_, err := s.commandClient.TrashedSaldo(ctx, &pb.FindByIdSaldoRequest{SaldoId: s.saldoID})
	s.NoError(err)
}

func (s *SaldoGapiTestSuite) Test5_Restore() {
	s.Require().NotZero(s.saldoID)
	ctx := context.Background()

	_, err := s.commandClient.RestoreSaldo(ctx, &pb.FindByIdSaldoRequest{SaldoId: s.saldoID})
	s.NoError(err)
}

func (s *SaldoGapiTestSuite) Test6_DeletePermanent() {
	s.Require().NotZero(s.saldoID)
	ctx := context.Background()

	_, err := s.commandClient.DeleteSaldoPermanent(ctx, &pb.FindByIdSaldoRequest{SaldoId: s.saldoID})
	s.NoError(err)
}

func TestSaldoGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoGapiTestSuite))
}

package topup_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	topup_repo "github.com/MamangRust/monolith-payment-gateway-topup/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/MamangRust/monolith-payment-gateway-topup/service"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gapi "github.com/MamangRust/monolith-payment-gateway-topup/handler"
)

type TopupGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	commandClient pb.TopupCommandServiceClient
	queryClient   pb.TopupQueryServiceClient
	conn        *grpc.ClientConn
	
	userRepo     user_repo.UserCommandRepository
	cardRepo     card_repo.CardCommandRepository
	saldoRepo    saldo_repo.Repositories
	topupRepo    topup_repo.Repositories

	cardNumber string
	topupID    int32
}

func (s *TopupGapiTestSuite) SetupSuite() {
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
	
	userRepos := user_repo.NewRepositories(queries)
	cardRepos := card_repo.NewRepositories(queries)
	saldoRepos := saldo_repo.NewRepositories(queries)
	
	cardAdapter := &topupCardRepoAdapter{
		CardQueryRepository:   cardRepos.CardQuery,
		CardCommandRepository: cardRepos.CardCommand,
	}
	s.topupRepo = topup_repo.NewRepositories(queries, cardAdapter, saldoRepos)
	s.userRepo = userRepos.UserCommand()
	s.cardRepo = cardRepos.CardCommand
	s.saldoRepo = saldoRepos

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	topupService := service.NewService(&service.Deps{
		Kafka:        nil,
		Cache:        cacheStore,
		Repositories: s.topupRepo,
		Logger:       log,
	})

	// Seed User, Card, Saldo
	user, err := s.userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Topup", LastName: "Gapi", Email: "topup.gapi@test.com", Password: "password123",
	})
	s.Require().NoError(err)
	
	card, err := s.cardRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "444", CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber
	
	_, err = s.saldoRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber: s.cardNumber, TotalBalance: 0,
	})
	s.Require().NoError(err)

	// Start gRPC Server
	topupHandler := gapi.NewHandler(topupService)
	server := grpc.NewServer()
	pb.RegisterTopupCommandServiceServer(server, topupHandler)
	pb.RegisterTopupQueryServiceServer(server, topupHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.commandClient = pb.NewTopupCommandServiceClient(conn)
	s.queryClient = pb.NewTopupQueryServiceClient(conn)
}

func (s *TopupGapiTestSuite) TearDownSuite() {
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

func (s *TopupGapiTestSuite) Test1_Create() {
	ctx := context.Background()

	createReq := &pb.CreateTopupRequest{
		CardNumber:  s.cardNumber,
		TopupAmount: 100000,
		TopupMethod: "bri",
	}
	res, err := s.commandClient.CreateTopup(ctx, createReq)
	s.NoError(err)
	s.Equal(int32(100000), res.Data.TopupAmount)

	s.topupID = res.Data.Id

	// Verify balance
	saldo, _ := s.saldoRepo.FindByCardNumber(ctx, s.cardNumber)
	s.Equal(int32(100000), saldo.TotalBalance)
}

func (s *TopupGapiTestSuite) Test2_FindById() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	found, err := s.queryClient.FindByIdTopup(ctx, &pb.FindByIdTopupRequest{TopupId: s.topupID})
	s.NoError(err)
	s.Equal(s.topupID, found.Data.Id)
}

func (s *TopupGapiTestSuite) Test3_Update() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	updateReq := &pb.UpdateTopupRequest{
		TopupId:     s.topupID,
		CardNumber:  s.cardNumber,
		TopupAmount: 150000,
		TopupMethod: "bri",
	}
	updated, err := s.commandClient.UpdateTopup(ctx, updateReq)
	s.NoError(err)
	s.Equal(int32(150000), updated.Data.TopupAmount)

	// Verify adjusted balance
	saldo, _ := s.saldoRepo.FindByCardNumber(ctx, s.cardNumber)
	s.Equal(int32(150000), saldo.TotalBalance)
}

func (s *TopupGapiTestSuite) Test4_Trashed() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	_, err := s.commandClient.TrashedTopup(ctx, &pb.FindByIdTopupRequest{TopupId: s.topupID})
	s.NoError(err)
}

func (s *TopupGapiTestSuite) Test5_Restore() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	_, err := s.commandClient.RestoreTopup(ctx, &pb.FindByIdTopupRequest{TopupId: s.topupID})
	s.NoError(err)
}

func (s *TopupGapiTestSuite) Test6_DeletePermanent() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	_, err := s.commandClient.DeleteTopupPermanent(ctx, &pb.FindByIdTopupRequest{TopupId: s.topupID})
	s.NoError(err)
}
func TestTopupGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupGapiTestSuite))
}

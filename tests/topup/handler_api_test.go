package topup_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/topup"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	app_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	gapi "github.com/MamangRust/monolith-payment-gateway-topup/handler"
	topup_repo "github.com/MamangRust/monolith-payment-gateway-topup/repository"
	"github.com/MamangRust/monolith-payment-gateway-topup/service"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TopupHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	commandClient pb.TopupCommandServiceClient
	queryClient   pb.TopupQueryServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
	
	userRepo  user_repo.UserCommandRepository
	cardRepo  card_repo.CardCommandRepository
	saldoRepo saldo_repo.Repositories
	topupRepo topup_repo.Repositories

	cardNumber string
	topupID    int
}

func (s *TopupHandlerTestSuite) SetupSuite() {
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
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	topupService := service.NewService(&service.Deps{
		Kafka:        nil,
		Cache:        cacheStore,
		Repositories: s.topupRepo,
		Logger:       log,
	})

	// Seed User and Card
	user, err := s.userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Topup",
		LastName:  "Owner",
		Email:     "topup.handler@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	card, err := s.cardRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber

	// Seed Saldo
	_, err = s.saldoRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)

	// Start gRPC Server
	topupGapiHandler := gapi.NewHandler(topupService)
	server := grpc.NewServer()
	pb.RegisterTopupCommandServiceServer(server, topupGapiHandler)
	pb.RegisterTopupQueryServiceServer(server, topupGapiHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	// Create gRPC Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.commandClient = pb.NewTopupCommandServiceClient(conn)
	s.queryClient = pb.NewTopupQueryServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)

	api.RegisterTopupHandler(&api.DepsTopup{
		Client:     s.conn,
		E:          s.router,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiErrorHandler,
	})
}

func (s *TopupHandlerTestSuite) TearDownSuite() {
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

func (s *TopupHandlerTestSuite) Test1_CreateTopup() {
	req := requests.CreateTopupRequest{
		CardNumber:  s.cardNumber,
		TopupAmount: 100000,
		TopupMethod: "visa",
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/topup-command/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	topupData := createRes["data"].(map[string]interface{})
	s.topupID = int(topupData["id"].(float64))
}

func (s *TopupHandlerTestSuite) Test2_FindById() {
	s.Require().NotZero(s.topupID)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/topup-query/%d", s.topupID), nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *TopupHandlerTestSuite) Test3_Update() {
	s.Require().NotZero(s.topupID)
	updateReq := requests.UpdateTopupRequest{
		TopupID:     &s.topupID,
		CardNumber:  s.cardNumber,
		TopupAmount: 150000,
		TopupMethod: "mastercard",
	}
	updateBody, _ := json.Marshal(updateReq)
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/topup-command/update/%d", s.topupID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *TopupHandlerTestSuite) Test4_DeletePermanent() {
	s.Require().NotZero(s.topupID)
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/topup-command/permanent/%d", s.topupID), nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestTopupHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupHandlerTestSuite))
}

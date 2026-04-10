package saldo_test

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

	api "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/saldo"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	app_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	gapi "github.com/MamangRust/monolith-payment-gateway-saldo/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SaldoHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	commandClient pb.SaldoCommandServiceClient
	queryClient   pb.SaldoQueryServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
	
	userRepo  user_repo.UserCommandRepository
	cardRepo  card_repo.CardCommandRepository
	saldoRepo saldo_repo.Repositories

	cardNumber string
	saldoID    int
}

func (s *SaldoHandlerTestSuite) SetupSuite() {
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
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	saldoService := service.NewService(&service.Deps{
		Repositories: s.saldoRepo,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User and Card
	user, err := s.userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Saldo",
		LastName:  "Owner",
		Email:     "saldo.handler@example.com",
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

	// Start gRPC Server
	saldoGapiHandler := gapi.NewHandler(saldoService)
	server := grpc.NewServer()
	pb.RegisterSaldoCommandServiceServer(server, saldoGapiHandler)
	pb.RegisterSaldoQueryServiceServer(server, saldoGapiHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	// Create gRPC Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.commandClient = pb.NewSaldoCommandServiceClient(conn)
	s.queryClient = pb.NewSaldoQueryServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)

	api.RegisterSaldoHandler(&api.DepsSaldo{
		Client:     s.conn,
		E:          s.router,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiErrorHandler,
	})
}

func (s *SaldoHandlerTestSuite) TearDownSuite() {
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

func (s *SaldoHandlerTestSuite) Test1_CreateSaldo() {
	req := requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 100000,
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/saldo-command/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	saldoData := createRes["data"].(map[string]interface{})
	s.saldoID = int(saldoData["id"].(float64))
}

func (s *SaldoHandlerTestSuite) Test2_FindById() {
	s.Require().NotZero(s.saldoID)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldo-query/%d", s.saldoID), nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *SaldoHandlerTestSuite) Test3_FindByCardNumber() {
	s.Require().NotEmpty(s.cardNumber)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldo-query/card_number/%s", s.cardNumber), nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *SaldoHandlerTestSuite) Test4_Update() {
	s.Require().NotZero(s.saldoID)
	updateReq := requests.UpdateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 150000,
	}
	updateBody, _ := json.Marshal(updateReq)
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/saldo-command/update/%d", s.saldoID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *SaldoHandlerTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.saldoID)
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/saldo-command/permanent/%d", s.saldoID), nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestSaldoHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoHandlerTestSuite))
}

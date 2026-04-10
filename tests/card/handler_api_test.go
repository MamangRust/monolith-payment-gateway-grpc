package card_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	cardhandler "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/card"
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
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CardApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echoApp     *echo.Echo
	userID      int
	cardID      int
}

func (s *CardApiTestSuite) SetupSuite() {
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

	cardGapiHandler := handler.NewHandler(cardService)
	server := grpc.NewServer()
	pb.RegisterCardQueryServiceServer(server, cardGapiHandler)
	pb.RegisterCardCommandServiceServer(server, cardGapiHandler)
	pb.RegisterCardDashboardServiceServer(server, cardGapiHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.echoApp = echo.New()
	obs, _ := observability.NewObservability("test", log)
	apiHandler := errors.NewApiHandler(obs, log)

	cardhandler.RegisterCardHandler(&cardhandler.DepsCard{
		Client:     conn,
		E:          s.echoApp,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiHandler,
	})

	// Create user
	user, err := userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Api",
		LastName:  "Card",
		Email:     "api.card@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Auth Bypass
	s.echoApp.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("userID", strconv.Itoa(s.userID))
			return next(c)
		}
	})
}

func (s *CardApiTestSuite) TearDownSuite() {
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *CardApiTestSuite) Test1_CreateCard() {
	req := requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	}
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/card-command/create", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	s.echoApp.ServeHTTP(rec, httpReq)

	s.Equal(http.StatusOK, rec.Code)
	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	s.cardID = int(data["id"].(float64))
}

func (s *CardApiTestSuite) Test2_FindById() {
	s.Require().NotZero(s.cardID)
	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card-query/%d", s.cardID), nil)

	s.echoApp.ServeHTTP(rec, httpReq)

	s.Equal(http.StatusOK, rec.Code)
	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	s.Equal(float64(s.cardID), data["id"].(float64))
}

func (s *CardApiTestSuite) Test3_UpdateCard() {
	s.Require().NotZero(s.cardID)
	req := requests.UpdateCardRequest{
		CardID:       s.cardID,
		UserID:       s.userID,
		CardType:     "credit",
		ExpireDate:   time.Now().AddDate(6, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card-command/update/%d", s.cardID), bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	s.echoApp.ServeHTTP(rec, httpReq)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CardApiTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.cardID)

	// Trash
	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card-command/trashed/%d", s.cardID), nil)
	s.echoApp.ServeHTTP(rec, httpReq)
	s.Equal(http.StatusOK, rec.Code)

	// Restore
	rec = httptest.NewRecorder()
	httpReq = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card-command/restore/%d", s.cardID), nil)
	s.echoApp.ServeHTTP(rec, httpReq)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *CardApiTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.cardID)

	// Trash first
	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card-command/trashed/%d", s.cardID), nil)
	s.echoApp.ServeHTTP(rec, httpReq)

	// Delete permanent
	rec = httptest.NewRecorder()
	httpReq = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/card-command/permanent/%d", s.cardID), nil)
	s.echoApp.ServeHTTP(rec, httpReq)
	s.Equal(http.StatusOK, rec.Code)
}

func TestCardApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardApiTestSuite))
}

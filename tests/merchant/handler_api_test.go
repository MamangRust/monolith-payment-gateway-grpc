package merchant_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/merchant"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	app_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-merchant/handler"
	"github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	"github.com/MamangRust/monolith-payment-gateway-merchant/service"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	conn        *grpc.ClientConn
	router      *echo.Echo
	userRepo    user_repo.UserCommandRepository
	userID      int
	merchantID  int
}

func (s *MerchantHandlerTestSuite) SetupSuite() {
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
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	merchantService := service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	user, err := s.userRepo.CreateUser(s.ts.Ctx, &requests.CreateUserRequest{
		FirstName: "Handler",
		LastName:  "Merchant",
		Email:     "handler.merchant@example.com",
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

	// Create gRPC Client for Echo
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn

	// Setup Echo
	s.router = echo.New()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)

	api.RegisterMerchantHandler(&api.DepsMerchant{
		Client:     conn,
		E:          s.router,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiErrorHandler,
	})
}

func (s *MerchantHandlerTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantHandlerTestSuite) Test1_CreateMerchant() {
	req := requests.CreateMerchantRequest{
		Name:   "Handler Merchant",
		UserID: s.userID,
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/merchant-command/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	merchantData := createRes["data"].(map[string]interface{})
	s.merchantID = int(merchantData["id"].(float64))
}

func (s *MerchantHandlerTestSuite) Test2_FindMerchantById() {
	s.Require().NotZero(s.merchantID)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-query/%d", s.merchantID), nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantHandlerTestSuite))
}

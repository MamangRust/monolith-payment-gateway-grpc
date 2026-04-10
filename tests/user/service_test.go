package user_test

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/MamangRust/monolith-payment-gateway-user/service"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type UserServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	userService service.Service
	userID      int
}

func (s *UserServiceTestSuite) SetupSuite() {
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

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	s.userService = service.NewService(&service.Deps{
		Cache:        cacheStore,
		Repositories: repos,
		Hash:         hasher,
		Logger:       log,
	})
}

func (s *UserServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *UserServiceTestSuite) Test1_CreateUser() {
	ctx := context.Background()

	req := &requests.CreateUserRequest{
		FirstName: "User",
		LastName:  "Service",
		Email:     "user.service@example.com",
		Password:  "password123",
	}
	user, err := s.userService.CreateUser(ctx, req)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(req.Email, user.Email)
	s.userID = int(user.UserID)
}

func (s *UserServiceTestSuite) Test2_FindUserById() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.userService.FindByID(ctx, s.userID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.userID, int(found.UserID))
}

func (s *UserServiceTestSuite) Test3_UpdateUser() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updateReq := &requests.UpdateUserRequest{
		UserID:    &s.userID,
		FirstName: "Updated",
	}
	updated, err := s.userService.UpdateUser(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated", updated.Firstname)
}

func (s *UserServiceTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	_, err := s.userService.TrashedUser(ctx, s.userID)
	s.NoError(err)

	_, err = s.userService.RestoreUser(ctx, s.userID)
	s.NoError(err)
}

func (s *UserServiceTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	success, err := s.userService.DeleteUserPermanent(ctx, s.userID)
	s.NoError(err)
	s.True(success)
}

func TestUserServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserServiceTestSuite))
}

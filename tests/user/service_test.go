package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/MamangRust/monolith-payment-gateway-user/service"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"

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
	userRepo    repository.Repositories
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
	s.userRepo = repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	s.userService = service.NewService(&service.Deps{
		Cache:        cacheStore,
		Repositories: s.userRepo,
		Hash:         hasher,
		Logger:       log,
	})

	// Flush Redis to ensure isolation
	s.redisClient.FlushAll(s.ts.Ctx)
}

func (s *UserServiceTestSuite) TearDownSuite() {
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

func (s *UserServiceTestSuite) Test1_CreateUser() {
	ctx := context.Background()

	req := &requests.CreateUserRequest{
		FirstName: fmt.Sprintf("User-%d", time.Now().UnixNano()),
		LastName:  "Service",
		Email:     fmt.Sprintf("user.service.%d@example.com", time.Now().UnixNano()),
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

func (s *UserServiceTestSuite) Test3_FindAll() {
	ctx := context.Background()
	
	req := &requests.FindAllUsers{
		Search:   "User",
		Page:     1,
		PageSize: 10,
	}
	users, total, err := s.userService.FindAll(ctx, req)
	s.NoError(err)
	s.NotNil(users)
	s.NotZero(*total)
}

func (s *UserServiceTestSuite) Test4_FindByActive() {
	ctx := context.Background()
	
	req := &requests.FindAllUsers{
		Search:   "User",
		Page:     1,
		PageSize: 10,
	}
	users, total, err := s.userService.FindByActive(ctx, req)
	s.NoError(err)
	s.NotNil(users)
	s.NotZero(*total)
}

func (s *UserServiceTestSuite) Test5_UpdateUser() {
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

func (s *UserServiceTestSuite) Test6_TrashAndRestore() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	// Trash
	_, err := s.userService.TrashedUser(ctx, s.userID)
	s.NoError(err)

	// Find By Trashed
	req := &requests.FindAllUsers{
		Page:     1,
		PageSize: 10,
	}
	trashed, total, err := s.userService.FindByTrashed(ctx, req)
	s.NoError(err)
	s.NotNil(trashed)
	s.NotZero(*total)

	// Restore
	_, err = s.userService.RestoreUser(ctx, s.userID)
	s.NoError(err)
}

func (s *UserServiceTestSuite) Test7_BulkOperations() {
	ctx := context.Background()

	// Trash for bulk test
	_, err := s.userService.TrashedUser(ctx, s.userID)
	s.NoError(err)

	// Restore All
	success, err := s.userService.RestoreAllUser(ctx)
	s.NoError(err)
	s.True(success)

	// Trash again for delete all
	_, err = s.userService.TrashedUser(ctx, s.userID)
	s.NoError(err)

	// Delete All Permanent
	success, err = s.userService.DeleteAllUserPermanent(ctx)
	s.NoError(err)
	s.True(success)

	// Reset userID since it's deleted
	s.userID = 0
}

func (s *UserServiceTestSuite) Test8_DeletePermanent() {
	ctx := context.Background()
	
	// Create another user for permanent delete test
	req := &requests.CreateUserRequest{
		FirstName: "DeleteMe",
		LastName:  "Permanent",
		Email:     fmt.Sprintf("delete.me.%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	}
	user, err := s.userService.CreateUser(ctx, req)
	s.NoError(err)
	
	uid := int(user.UserID)
	
	// Must be trashed first
	_, err = s.userService.TrashedUser(ctx, uid)
	s.NoError(err)

	success, err := s.userService.DeleteUserPermanent(ctx, uid)
	s.NoError(err)
	s.True(success)
}

func TestUserServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserServiceTestSuite))
}

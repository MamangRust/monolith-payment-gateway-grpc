package topup_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-topup/repository"
	"github.com/MamangRust/monolith-payment-gateway-topup/service"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo_impl "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo_impl "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type testCardRepo struct {
	db *db.Queries
}

func (r *testCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*db.GetUserEmailByCardNumberRow, error) {
	return r.db.GetUserEmailByCardNumber(ctx, card_number)
}

func (r *testCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	return r.db.GetCardByCardNumber(ctx, card_number)
}

func (r *testCardRepo) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	expireDate := pgtype.Date{
		Time:  request.ExpireDate,
		Valid: true,
	}
	return r.db.UpdateCard(ctx, db.UpdateCardParams{
		CardID:       int32(request.CardID),
		CardType:     request.CardType,
		ExpireDate:   expireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	})
}

type testSaldoRepo struct {
	db *db.Queries
}

func (r *testSaldoRepo) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	return r.db.GetSaldoByCardNumber(ctx, card_number)
}

func (r *testSaldoRepo) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
	return r.db.UpdateSaldoBalance(ctx, db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	})
}

type TopupServiceTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	topupService service.Service
	dbPool       *pgxpool.Pool
	topupID      int
	cardNumber   string
	userID       int
}

func (s *TopupServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	queries := db.New(pool)
	
	cardRepo := &testCardRepo{db: queries}
	saldoRepo := &testSaldoRepo{db: queries}
	repos := repository.NewRepositories(queries, cardRepo, saldoRepo)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	s.topupService = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	userRepo := user_repo.NewUserCommandRepository(queries)
	user, err := userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Topup",
		LastName:  "Tester",
		Email:     fmt.Sprintf("topup.tester.%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Seed Card
	cardCmdRepo := card_repo_impl.NewCardCommandRepository(queries)
	card, err := cardCmdRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber

	// Seed Saldo
	saldoCmdRepo := saldo_repo_impl.NewSaldoCommandRepository(queries)
	_, err = saldoCmdRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)
}

func (s *TopupServiceTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TopupServiceTestSuite) Test1_TopupLifecycle() {
	ctx := context.Background()

	// Create Topup
	createReq := &requests.CreateTopupRequest{
		CardNumber:  s.cardNumber,
		TopupAmount: 100000,
		TopupMethod: "gopay",
	}
	res, err := s.topupService.CreateTopup(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.topupID = int(res.TopupID)
	s.Equal("success", res.Status)

	// FindById
	found, err := s.topupService.FindById(ctx, s.topupID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.topupID), found.TopupID)

	// Update Topup
	updateReq := &requests.UpdateTopupRequest{
		TopupID:     &s.topupID,
		CardNumber:  s.cardNumber,
		TopupAmount: 200000,
		TopupMethod: "dana",
	}
	updated, err := s.topupService.UpdateTopup(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("success", updated.Status)
}

func (s *TopupServiceTestSuite) Test2_QueryOperations() {
	ctx := context.Background()

	// FindAll
	all, total, err := s.topupService.FindAll(ctx, &requests.FindAllTopups{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// FindByActive
	active, totalActive, err := s.topupService.FindByActive(ctx, &requests.FindAllTopups{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(active)
	s.GreaterOrEqual(*totalActive, 1)

	// FindByCardNumber
	byCard, _, err := s.topupService.FindAllByCardNumber(ctx, &requests.FindAllTopupsByCardNumber{
		CardNumber: s.cardNumber,
		Page:       1,
		PageSize:   10,
	})
	s.NoError(err)
	s.NotNil(byCard)
	s.GreaterOrEqual(len(byCard), 1)
}

func (s *TopupServiceTestSuite) Test3_TrashAndRestore() {
	ctx := context.Background()
	s.Require().NotZero(s.topupID)

	// Trash
	trashed, err := s.topupService.TrashedTopup(ctx, s.topupID)
	s.NoError(err)
	s.NotNil(trashed)

	// FindByTrashed
	trashedList, totalTrashed, err := s.topupService.FindByTrashed(ctx, &requests.FindAllTopups{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// Restore
	restored, err := s.topupService.RestoreTopup(ctx, s.topupID)
	s.NoError(err)
	s.NotNil(restored)
}

func (s *TopupServiceTestSuite) Test4_BulkOperations() {
	ctx := context.Background()

	// Restore All
	ok, err := s.topupService.RestoreAllTopup(ctx)
	s.NoError(err)
	s.True(ok)

	// Delete All Permanent
	ok, err = s.topupService.DeleteAllTopupPermanent(ctx)
	s.NoError(err)
	s.True(ok)
}

func TestTopupServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupServiceTestSuite))
}

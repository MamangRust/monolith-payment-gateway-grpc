package transaction_test

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
	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	"github.com/MamangRust/monolith-payment-gateway-transaction/service"
	user_repo_impl "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo_impl "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo_impl "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	merchant_repo_impl "github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type testMerchantRepo struct {
	db *db.Queries
}

func (r *testMerchantRepo) FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error) {
	return r.db.GetMerchantByApiKey(ctx, api_key)
}

type testCardRepo struct {
	db *db.Queries
}

func (r *testCardRepo) FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	return r.db.GetCardByUserID(ctx, int32(user_id))
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

type TransactionServiceTestSuite struct {
	suite.Suite
	ts                 *tests.TestSuite
	transactionService service.Service
	dbPool             *pgxpool.Pool
	transactionID      int
	cardNumber         string
	merchantID         int
	merchantApiKey     string
	userID             int
}

func (s *TransactionServiceTestSuite) SetupSuite() {
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
	merchantRepo := &testMerchantRepo{db: queries}
	repos := repository.NewRepositories(queries, saldoRepo, cardRepo, merchantRepo)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	s.transactionService = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	userRepo := user_repo_impl.NewUserCommandRepository(queries)
	user, err := userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Transaction",
		LastName:  "Tester",
		Email:     fmt.Sprintf("transaction.tester.%d@example.com", time.Now().UnixNano()),
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

	// Seed Merchant
	merchantRepoImpl := merchant_repo_impl.NewMerchantCommandRepository(queries)
	merchant, err := merchantRepoImpl.CreateMerchant(context.Background(), &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   "Test Merchant",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)
	s.merchantApiKey = merchant.ApiKey

	// Seed Merchant Card & Saldo (needed for transaction credit)
	merchantCard, err := cardCmdRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID, // reusing user for simplicity
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "456",
		CardProvider: "mastercard",
	})
	s.Require().NoError(err)
	_, err = saldoCmdRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)
}

func (s *TransactionServiceTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransactionServiceTestSuite) Test1_TransactionLifecycle() {
	ctx := context.Background()

	// Create Transaction
	merchantID := int(s.merchantID)
	createReq := &requests.CreateTransactionRequest{
		CardNumber:      s.cardNumber,
		Amount:          100000,
		PaymentMethod:   "gopay",
		MerchantID:      &merchantID,
		TransactionTime: time.Now(),
	}
	res, err := s.transactionService.Create(ctx, s.merchantApiKey, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.transactionID = int(res.TransactionID)
	s.Equal("success", res.Status)

	// FindById
	found, err := s.transactionService.FindById(ctx, s.transactionID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.transactionID), found.TransactionID)

	// Update Transaction
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID:   &s.transactionID,
		CardNumber:      s.cardNumber,
		Amount:          200000,
		PaymentMethod:   "dana",
		MerchantID:      &merchantID,
		TransactionTime: time.Now(),
	}
	updated, err := s.transactionService.Update(ctx, s.merchantApiKey, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(200000), updated.Amount)
}

func (s *TransactionServiceTestSuite) Test2_QueryOperations() {
	ctx := context.Background()

	// FindAll
	all, total, err := s.transactionService.FindAll(ctx, &requests.FindAllTransactions{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// FindByActive
	active, totalActive, err := s.transactionService.FindByActive(ctx, &requests.FindAllTransactions{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(active)
	s.GreaterOrEqual(*totalActive, 1)

	// FindByCardNumber
	byCard, totalCard, err := s.transactionService.FindAllByCardNumber(ctx, &requests.FindAllTransactionCardNumber{
		CardNumber: s.cardNumber,
		Page:       1,
		PageSize:   10,
	})
	s.NoError(err)
	s.NotNil(byCard)
	s.GreaterOrEqual(*totalCard, 1)

	// FindByMerchantId
	byMerchant, err := s.transactionService.FindTransactionByMerchantId(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(byMerchant)
	s.GreaterOrEqual(len(byMerchant), 1)
}

func (s *TransactionServiceTestSuite) Test3_TrashAndRestore() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)

	// Trash
	trashed, err := s.transactionService.TrashedTransaction(ctx, s.transactionID)
	s.NoError(err)
	s.NotNil(trashed)

	// FindByTrashed
	trashedList, totalTrashed, err := s.transactionService.FindByTrashed(ctx, &requests.FindAllTransactions{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// Restore
	restored, err := s.transactionService.RestoreTransaction(ctx, s.transactionID)
	s.NoError(err)
	s.NotNil(restored)
}

func (s *TransactionServiceTestSuite) Test4_BulkOperations() {
	ctx := context.Background()

	// Restore All
	ok, err := s.transactionService.RestoreAllTransaction(ctx)
	s.NoError(err)
	s.True(ok)

	// Delete All Permanent
	ok, err = s.transactionService.DeleteAllTransactionPermanent(ctx)
	s.NoError(err)
	s.True(ok)
}

func TestTransactionServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionServiceTestSuite))
}

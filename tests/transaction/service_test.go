package transaction_test

import (
	"context"
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
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	ts                 *tests.TestSuite
	dbPool             *pgxpool.Pool
	redisClient        *redis.Client
	transactionService service.Service
	
	// Repositories for seeding
	userRepo     user_repo.UserCommandRepository
	cardRepo     card_repo.Repositories
	saldoRepo    saldo_repo.Repositories
	merchantRepo merchant_repo.Repositories

	customerCardNumber string
	merchantID         int
	merchantApiKey     string
	merchantCardNumber string
	transactionID      int
}

func (s *TransactionServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.userRepo = user_repo.NewUserCommandRepository(queries)
	s.cardRepo = *card_repo.NewRepositories(queries)
	s.saldoRepo = saldo_repo.NewRepositories(queries)
	s.merchantRepo = merchant_repo.NewRepositories(queries)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	_ = log
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	cardRepoWrapper := &transactionCardRepo{
		query:   s.cardRepo.CardQuery,
		command: s.cardRepo.CardCommand,
	}

	transactionRepos := repository.NewRepositories(queries, s.saldoRepo, cardRepoWrapper, s.merchantRepo)
	s.transactionService = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: transactionRepos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User, Card, Merchant, Saldo
	ctx := context.Background()
	user, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Tx", LastName: "Owner", Email: "tx.service@example.com", Password: "password123",
	})
	card, _ := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(2, 0, 0), CVV: "123", CardProvider: "visa",
	})
	s.customerCardNumber = card.CardNumber

	merchant, _ := s.merchantRepo.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		Name: "Service Merchant", UserID: int(user.UserID),
	})
	s.merchantID = int(merchant.MerchantID)
	s.merchantApiKey = merchant.ApiKey

	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber: s.customerCardNumber, TotalBalance: 1000000,
	})

	merchantCard, _ := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(2, 0, 0), CVV: "456", CardProvider: "mastercard",
	})
	s.merchantCardNumber = merchantCard.CardNumber
	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber: s.merchantCardNumber, TotalBalance: 0,
	})
}

func (s *TransactionServiceTestSuite) TearDownSuite() {
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

func (s *TransactionServiceTestSuite) Test1_CreateTransaction() {
	ctx := context.Background()
	merchantID := s.merchantID
	req := &requests.CreateTransactionRequest{
		CardNumber:      s.customerCardNumber,
		Amount:          100000,
		PaymentMethod:   "visa",
		MerchantID:      &merchantID,
		TransactionTime: time.Now(),
	}
	tx, err := s.transactionService.Create(ctx, s.merchantApiKey, req)
	s.NoError(err)
	s.NotNil(tx)
	s.transactionID = int(tx.TransactionID)
}

func (s *TransactionServiceTestSuite) Test2_FindTransactionById() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	res, err := s.transactionService.FindById(ctx, s.transactionID)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(int32(s.transactionID), res.TransactionID)
}

func (s *TransactionServiceTestSuite) Test3_UpdateTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	merchantID := s.merchantID
	req := &requests.UpdateTransactionRequest{
		TransactionID:   &s.transactionID,
		CardNumber:      s.customerCardNumber,
		Amount:          150000,
		MerchantID:      &merchantID,
		PaymentMethod:   "visa",
		TransactionTime: time.Now(),
	}
	res, err := s.transactionService.Update(ctx, s.merchantApiKey, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(int32(150000), res.Amount)
}

func (s *TransactionServiceTestSuite) Test4_TrashedTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	res, err := s.transactionService.TrashedTransaction(ctx, s.transactionID)
	s.NoError(err)
	s.NotNil(res.DeletedAt)
}

func (s *TransactionServiceTestSuite) Test5_RestoreTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	res, err := s.transactionService.RestoreTransaction(ctx, s.transactionID)
	s.NoError(err)
	s.True(res.DeletedAt.Time.IsZero())
}

func (s *TransactionServiceTestSuite) Test6_PermanentDeleteTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	success, err := s.transactionService.DeleteTransactionPermanent(ctx, s.transactionID)
	s.NoError(err)
	s.True(success)
}

func TestTransactionServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionServiceTestSuite))
}

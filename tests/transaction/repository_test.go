package transaction_test

import (
	"context"
	"testing"
	"time"


	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	ts            *tests.TestSuite
	dbPool        *pgxpool.Pool
	commandRepo   repository.TransactionCommandRepository
	queryRepo     repository.TransactionQueryRepository
	
	// Repositories for seeding
	userRepo     user_repo.UserCommandRepository
	cardRepo     card_repo.Repositories
	merchantRepo merchant_repo.Repositories

	customerCardNumber string
	merchantID         int
	transactionID      int
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.userRepo = user_repo.NewUserCommandRepository(queries)
	s.cardRepo = *card_repo.NewRepositories(queries)
	s.merchantRepo = merchant_repo.NewRepositories(queries)

	transactionRepos := repository.NewRepositories(queries, nil, nil, nil) // We don't need saldo/card/merchant for the base repo tests if we seed directly
	s.commandRepo = transactionRepos
	s.queryRepo = transactionRepos
}

func (s *TransactionRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *TransactionRepositoryTestSuite) Test1_CreateTransaction() {
	ctx := context.Background()
	
	// Seed 
	user, err := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Repo", LastName: "Owner", Email: "repo@test.com", Password: "password123",
	})
	s.Require().NoError(err)

	card, err := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(2, 0, 0), CVV: "123", CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.customerCardNumber = card.CardNumber

	merchant, err := s.merchantRepo.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		Name: "Repo Merchant", UserID: int(user.UserID),
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	req := &requests.CreateTransactionRequest{
		CardNumber:      s.customerCardNumber,
		Amount:          100000,
		MerchantID:      &s.merchantID,
		PaymentMethod:   "visa",
		TransactionTime: time.Now(),
	}

	res, err := s.commandRepo.CreateTransaction(ctx, req)
	s.NoError(err)
	s.Require().NotNil(res)
	s.Equal(int32(req.Amount), res.Amount)
	s.transactionID = int(res.TransactionID)
}

func (s *TransactionRepositoryTestSuite) Test2_FindById() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	res, err := s.queryRepo.FindById(ctx, s.transactionID)
	s.NoError(err)
	s.Require().NotNil(res)
	s.Equal(int32(s.transactionID), res.TransactionID)
}

func (s *TransactionRepositoryTestSuite) Test3_UpdateTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	req := &requests.UpdateTransactionRequest{
		TransactionID:   &s.transactionID,
		CardNumber:      s.customerCardNumber,
		Amount:          200000,
		MerchantID:      &s.merchantID,
		PaymentMethod:   "visa",
		TransactionTime: time.Now(),
	}
	res, err := s.commandRepo.UpdateTransaction(ctx, req)
	s.NoError(err)
	s.Require().NotNil(res)
	s.Equal(int32(200000), res.Amount)
}

func (s *TransactionRepositoryTestSuite) Test4_TrashedTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	res, err := s.commandRepo.TrashedTransaction(ctx, s.transactionID)
	s.NoError(err)
	s.NotNil(res.DeletedAt)
}

func (s *TransactionRepositoryTestSuite) Test5_RestoreTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	res, err := s.commandRepo.RestoreTransaction(ctx, s.transactionID)
	s.NoError(err)
	s.True(res.DeletedAt.Time.IsZero())
}

func (s *TransactionRepositoryTestSuite) Test6_PermanentDeleteTransaction() {
	ctx := context.Background()
	s.Require().NotZero(s.transactionID)
	success, err := s.commandRepo.DeleteTransactionPermanent(ctx, s.transactionID)
	s.NoError(err)
	s.True(success)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

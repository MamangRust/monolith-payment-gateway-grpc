package saldo_test

import (
	"context"
	"testing"
	"time"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type SaldoRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   saldo_repo.Repositories
	
	userRepo  user_repo.UserCommandRepository
	cardRepo  card_repo.CardCommandRepository

	cardNumber string
	saldoID    int32
}

func (s *SaldoRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = saldo_repo.NewRepositories(queries)
	
	userRepos := user_repo.NewRepositories(queries)
	cardRepos := card_repo.NewRepositories(queries)
	s.userRepo = userRepos.UserCommand()
	s.cardRepo = cardRepos.CardCommand

	// Seed User and Card
	ctx := context.Background()
	user, err := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Saldo",
		LastName:  "Owner",
		Email:     "saldo.repo@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	
	card, err := s.cardRepo.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber
}

func (s *SaldoRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *SaldoRepositoryTestSuite) Test1_CreateSaldo() {
	ctx := context.Background()
	req := &requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 1000000,
	}

	saldo, err := s.repo.CreateSaldo(ctx, req)
	s.NoError(err)
	s.Require().NotNil(saldo)
	s.Equal(int32(req.TotalBalance), saldo.TotalBalance)
	s.saldoID = saldo.SaldoID
}

func (s *SaldoRepositoryTestSuite) Test2_FindByCardNumber() {
	s.Require().NotEmpty(s.cardNumber)
	ctx := context.Background()

	found, err := s.repo.FindByCardNumber(ctx, s.cardNumber)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(1000000), found.TotalBalance)
}

func (s *SaldoRepositoryTestSuite) Test3_UpdateBalance() {
	s.Require().NotEmpty(s.cardNumber)
	ctx := context.Background()

	req := &requests.UpdateSaldoBalance{
		CardNumber:   s.cardNumber,
		TotalBalance: 1200000,
	}
	updated, err := s.repo.UpdateSaldoBalance(ctx, req)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(1200000), updated.TotalBalance)
}

func (s *SaldoRepositoryTestSuite) Test4_TrashedAndRestore() {
	s.Require().NotZero(s.saldoID)
	ctx := context.Background()

	_, err := s.repo.TrashedSaldo(ctx, int(s.saldoID))
	s.NoError(err)

	_, err = s.repo.RestoreSaldo(ctx, int(s.saldoID))
	s.NoError(err)
}

func (s *SaldoRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.saldoID)
	ctx := context.Background()

	success, err := s.repo.DeleteSaldoPermanent(ctx, int(s.saldoID))
	s.NoError(err)
	s.True(success)
}

func TestSaldoRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoRepositoryTestSuite))
}

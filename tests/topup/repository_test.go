package topup_test

import (
	"context"
	"testing"
	"time"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	topup_repo "github.com/MamangRust/monolith-payment-gateway-topup/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TopupRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   topup_repo.Repositories
	
	userRepo  user_repo.UserCommandRepository
	cardRepo  card_repo.CardCommandRepository
	saldoRepo saldo_repo.Repositories

	cardNumber string
	topupID    int32
}

func (s *TopupRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	
	// Initialize repos from their modules
	userRepos := user_repo.NewRepositories(queries)
	cardRepos := card_repo.NewRepositories(queries)
	saldoRepos := saldo_repo.NewRepositories(queries)
	
	// Match topup repository interfaces
	cardAdapter := &topupCardRepoAdapter{
		CardQueryRepository:   cardRepos.CardQuery,
		CardCommandRepository: cardRepos.CardCommand,
	}
	s.repo = topup_repo.NewRepositories(queries, cardAdapter, saldoRepos)
	s.userRepo = userRepos.UserCommand()
	s.cardRepo = cardRepos.CardCommand
	s.saldoRepo = saldoRepos

	// Seed User and Card
	ctx := context.Background()
	user, err := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Topup",
		LastName:  "Owner",
		Email:     "topup.repo@example.com",
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

	_, err = s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 0,
	})
	s.Require().NoError(err)
}

func (s *TopupRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *TopupRepositoryTestSuite) Test1_CreateTopup() {
	ctx := context.Background()
	req := &requests.CreateTopupRequest{
		CardNumber:  s.cardNumber,
		TopupAmount: 50000,
		TopupMethod: "visa",
	}

	topup, err := s.repo.CreateTopup(ctx, req)
	s.NoError(err)
	s.NotNil(topup)
	s.Equal(int32(req.TopupAmount), topup.TopupAmount)
	s.topupID = topup.TopupID
}

func (s *TopupRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	found, err := s.repo.FindById(ctx, int(s.topupID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.topupID, found.TopupID)
}

func (s *TopupRepositoryTestSuite) Test3_UpdateStatus() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	req := &requests.UpdateTopupStatus{
		TopupID: int(s.topupID),
		Status:  "success",
	}
	updated, err := s.repo.UpdateTopupStatus(ctx, req)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("success", updated.Status)
}

func (s *TopupRepositoryTestSuite) Test4_TrashedAndRestore() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	_, err := s.repo.TrashedTopup(ctx, int(s.topupID))
	s.NoError(err)

	_, err = s.repo.RestoreTopup(ctx, int(s.topupID))
	s.NoError(err)
}

func (s *TopupRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.topupID)
	ctx := context.Background()

	success, err := s.repo.DeleteTopupPermanent(ctx, int(s.topupID))
	s.NoError(err)
	s.True(success)
}

func TestTopupRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupRepositoryTestSuite))
}

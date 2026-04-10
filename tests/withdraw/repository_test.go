package withdraw_test

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-test"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type WithdrawRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.WithdrawCommandRepository
}

func (s *WithdrawRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewWithdrawCommandRepository(queries)
}

func (s *WithdrawRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *WithdrawRepositoryTestSuite) TestCreateWithdraw() {
	// Seed card first
	userReq := &requests.CreateUserRequest{
		FirstName: "Withdraw",
		LastName:  "Owner",
		Email:     "withdrawowner@example.com",
		Password:  "password123",
	}
	queries := db.New(s.dbPool)
	userRepo := user_repo.NewUserCommandRepository(queries)
	user, err := userRepo.CreateUser(context.Background(), userReq)
	s.Require().NoError(err)

	cardRepos := card_repo.NewRepositories(queries)
	card, err := cardRepos.CardCommand.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)

	req := &requests.CreateWithdrawRequest{
		CardNumber:     card.CardNumber,
		WithdrawAmount: 50000,
		WithdrawTime:   time.Now(),
	}

	withdraw, err := s.repo.CreateWithdraw(context.Background(), req)
	s.NoError(err)
	s.NotNil(withdraw)
	s.Equal(int32(req.WithdrawAmount), withdraw.WithdrawAmount)
}

func TestWithdrawRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawRepositoryTestSuite))
}

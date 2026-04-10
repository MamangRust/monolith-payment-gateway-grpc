package card_test

import (
	"context"
	"testing"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-card/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type CardRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	repo     *repository.Repositories
	userRepo user_repo.Repositories
	userID   int
	cardID   int
}

func (s *CardRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)

	queries := db.New(pool)
	s.repo = repository.NewRepositories(queries)
	s.userRepo = user_repo.NewRepositories(queries)

	// Create a user for card ownership
	user, err := s.userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Card",
		LastName:  "Owner",
		Email:     "card.owner@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *CardRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *CardRepositoryTestSuite) Test1_CreateCard() {
	ctx := context.Background()
	req := &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	}

	res, err := s.repo.CardCommand.CreateCard(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.NotEmpty(res.CardNumber)
	s.cardID = int(res.CardID)
}

func (s *CardRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.cardID)
	ctx := context.Background()

	found, err := s.repo.CardQuery.FindById(ctx, s.cardID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.cardID), found.CardID)
}

func (s *CardRepositoryTestSuite) Test3_UpdateCard() {
	s.Require().NotZero(s.cardID)
	ctx := context.Background()

	req := &requests.UpdateCardRequest{
		CardID:       s.cardID,
		UserID:       s.userID,
		CardType:     "credit",
		ExpireDate:   time.Now().AddDate(6, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	res, err := s.repo.CardCommand.UpdateCard(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("credit", res.CardType)
}

func (s *CardRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.cardID)
	ctx := context.Background()

	trashed, err := s.repo.CardCommand.TrashedCard(ctx, s.cardID)
	s.NoError(err)
	s.NotNil(trashed)

	restored, err := s.repo.CardCommand.RestoreCard(ctx, s.cardID)
	s.NoError(err)
	s.NotNil(restored)
}

func (s *CardRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.cardID)
	ctx := context.Background()

	trashed, _ := s.repo.CardCommand.TrashedCard(ctx, s.cardID)
	s.NotNil(trashed)

	success, err := s.repo.CardCommand.DeleteCardPermanent(ctx, s.cardID)
	s.NoError(err)
	s.True(success)
}

func TestCardRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardRepositoryTestSuite))
}

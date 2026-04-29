package saldo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type SaldoRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	repo     repository.Repositories
	cardRepo *card_repo.Repositories
	userRepo user_repo.Repositories
	userID   int
}

func (s *SaldoRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)

	queries := db.New(pool)
	s.repo = repository.NewRepositories(queries)
	s.cardRepo = card_repo.NewRepositories(queries)
	s.userRepo = user_repo.NewRepositories(queries)

	// Create user
	user, err := s.userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Saldo",
		LastName:  "Owner",
		Email:     fmt.Sprintf("saldo.owner-%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *SaldoRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *SaldoRepositoryTestSuite) createSeedSaldo() (*db.CreateSaldoRow, error) {
    card, err := s.cardRepo.CardCommand.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	})
    if err != nil {
        return nil, err
    }

	return s.repo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 100000,
	})
}

func (s *SaldoRepositoryTestSuite) TestCreateSaldo() {
	ctx := context.Background()
    
    card, _ := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	})

	req := &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 100000,
	}

	res, err := s.repo.CreateSaldo(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *SaldoRepositoryTestSuite) TestFindAllSaldos() {
	_, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindAllSaldos(ctx, &requests.FindAllSaldos{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *SaldoRepositoryTestSuite) TestFindById() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	found, err := s.repo.FindById(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(saldo.SaldoID, found.SaldoID)
}

func (s *SaldoRepositoryTestSuite) TestFindByActive() {
	_, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindByActive(ctx, &requests.FindAllSaldos{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *SaldoRepositoryTestSuite) TestFindByTrashed() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	res, err := s.repo.FindByTrashed(ctx, &requests.FindAllSaldos{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *SaldoRepositoryTestSuite) TestUpdateSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	id := int(saldo.SaldoID)
	req := &requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   saldo.CardNumber,
		TotalBalance: 200000,
	}

	res, err := s.repo.UpdateSaldo(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *SaldoRepositoryTestSuite) TestTrashSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	trashed, err := s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.NotNil(trashed)
}

func (s *SaldoRepositoryTestSuite) TestRestoreSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	restored, err := s.repo.RestoreSaldo(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.NotNil(restored)
}

func (s *SaldoRepositoryTestSuite) TestDeleteSaldoPermanent() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	success, err := s.repo.DeleteSaldoPermanent(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.True(success)
}

func (s *SaldoRepositoryTestSuite) TestRestoreAllSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	success, err := s.repo.RestoreAllSaldo(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *SaldoRepositoryTestSuite) TestDeleteAllSaldoPermanent() {
	_, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	success, err := s.repo.DeleteAllSaldoPermanent(ctx)
	s.NoError(err)
	s.True(success)
}

func TestSaldoRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoRepositoryTestSuite))
}

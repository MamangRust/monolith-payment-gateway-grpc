package merchant_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type MerchantRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	repo     repository.Repositories
	userRepo user_repo.Repositories
	userID   int
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)

	queries := db.New(pool)
	s.repo = repository.NewRepositories(queries)
	s.userRepo = user_repo.NewRepositories(queries)

	// Seed User
	user, err := s.userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     fmt.Sprintf("merchant.owner-%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *MerchantRepositoryTestSuite) createSeedMerchant() (*db.CreateMerchantRow, error) {
	return s.repo.CreateMerchant(context.Background(), &requests.CreateMerchantRequest{
		Name:   fmt.Sprintf("Test Merchant-%d", time.Now().UnixNano()),
		UserID: s.userID,
	})
}

func (s *MerchantRepositoryTestSuite) TestCreateMerchant() {
	ctx := context.Background()
	req := &requests.CreateMerchantRequest{
		Name:   "Test Merchant",
		UserID: s.userID,
	}

	res, err := s.repo.CreateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *MerchantRepositoryTestSuite) TestFindAllMerchants() {
	_, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindAllMerchants(ctx, &requests.FindAllMerchants{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *MerchantRepositoryTestSuite) TestFindById() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	found, err := s.repo.FindByMerchantId(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(merchant.MerchantID, found.MerchantID)
}

func (s *MerchantRepositoryTestSuite) TestFindByActive() {
	_, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindByActive(ctx, &requests.FindAllMerchants{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *MerchantRepositoryTestSuite) TestFindByTrashed() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	res, err := s.repo.FindByTrashed(ctx, &requests.FindAllMerchants{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *MerchantRepositoryTestSuite) TestUpdateMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	id := int(merchant.MerchantID)
	req := &requests.UpdateMerchantRequest{
		MerchantID: &id,
		Name:       "Updated Merchant",
		UserID:     s.userID,
		Status:     "active",
	}

	res, err := s.repo.UpdateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *MerchantRepositoryTestSuite) TestTrashMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	trashed, err := s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(trashed)
}

func (s *MerchantRepositoryTestSuite) TestRestoreMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	restored, err := s.repo.RestoreMerchant(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(restored)
}

func (s *MerchantRepositoryTestSuite) TestDeleteMerchantPermanent() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	success, err := s.repo.DeleteMerchantPermanent(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.True(success)
}

func (s *MerchantRepositoryTestSuite) TestRestoreAllMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	success, err := s.repo.RestoreAllMerchant(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *MerchantRepositoryTestSuite) TestDeleteAllMerchantPermanent() {
	_, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	success, err := s.repo.DeleteAllMerchantPermanent(ctx)
	s.NoError(err)
	s.True(success)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}

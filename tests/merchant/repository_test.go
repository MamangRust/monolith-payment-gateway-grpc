package merchant_test

import (
	"context"
	"testing"

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
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repo       repository.MerchantCommandRepository
	queryRepo  repository.MerchantQueryRepository
	userRepo   user_repo.UserCommandRepository
	userID     int
	merchantID int
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantCommandRepository(queries)
	s.queryRepo = repository.NewMerchantQueryRepository(queries)
	s.userRepo = user_repo.NewUserCommandRepository(queries)

	// Seed User
	user, err := s.userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     "merchant.owner@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantRepositoryTestSuite) Test1_CreateMerchant() {
	req := &requests.CreateMerchantRequest{
		Name:   "Test Merchant",
		UserID: s.userID,
	}

	merchant, err := s.repo.CreateMerchant(context.Background(), req)
	s.NoError(err)
	s.NotNil(merchant)
	s.Equal(req.Name, merchant.Name)
	s.Equal(int32(s.userID), merchant.UserID)
	s.merchantID = int(merchant.MerchantID)
}

func (s *MerchantRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	found, err := s.queryRepo.FindByMerchantId(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.merchantID, int(found.MerchantID))
}

func (s *MerchantRepositoryTestSuite) Test3_UpdateMerchant() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	updateReq := &requests.UpdateMerchantRequest{
		MerchantID: &s.merchantID,
		Name:       "Updated Merchant",
		UserID:     s.userID,
		Status:     "active",
	}

	updated, err := s.repo.UpdateMerchant(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)
	s.Equal(updateReq.Status, updated.Status)
}

func (s *MerchantRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	_, err := s.repo.TrashedMerchant(ctx, s.merchantID)
	s.NoError(err)

	_, err = s.repo.RestoreMerchant(ctx, s.merchantID)
	s.NoError(err)
}

func (s *MerchantRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	success, err := s.repo.DeleteMerchantPermanent(ctx, s.merchantID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}

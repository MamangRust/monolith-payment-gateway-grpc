package merchant_test

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
	"github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	"github.com/MamangRust/monolith-payment-gateway-merchant/service"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type MerchantServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	merchantService service.Service
	dbPool          *pgxpool.Pool
	merchantID      int
	documentID      int
	userID          int
}

func (s *MerchantServiceTestSuite) SetupSuite() {
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
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	s.merchantService = service.NewService(&service.Deps{
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
		Kafka:        nil,
	})

	// Seed User
	userRepo := user_repo.NewUserCommandRepository(queries)
	user, err := userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Tester",
		Email:     fmt.Sprintf("merchant.tester.%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *MerchantServiceTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantServiceTestSuite) Test1_MerchantOperations() {
	ctx := context.Background()
	
	// Create Merchant
	createReq := &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   fmt.Sprintf("Test Merchant %d", time.Now().UnixNano()),
	}
	res, err := s.merchantService.MerchantCommandService().CreateMerchant(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.merchantID = int(res.MerchantID)

	// FindById
	found, err := s.merchantService.MerchantQueryService().FindById(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.merchantID), found.MerchantID)

	// Update Merchant
	updateReq := &requests.UpdateMerchantRequest{
		MerchantID: &s.merchantID,
		Name:       "Updated Merchant Name",
		UserID:     s.userID,
		Status:     "inactive",
	}
	updated, err := s.merchantService.MerchantCommandService().UpdateMerchant(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Merchant Name", updated.Name)

	// Update Status
	statusReq := &requests.UpdateMerchantStatusRequest{
		MerchantID: &s.merchantID,
		Status:     "active",
	}
	statusUpdated, err := s.merchantService.MerchantCommandService().UpdateMerchantStatus(ctx, statusReq)
	s.NoError(err)
	s.NotNil(statusUpdated)
	s.Equal("active", statusUpdated.Status)
}

func (s *MerchantServiceTestSuite) Test2_DocumentOperations() {
	ctx := context.Background()
	s.Require().NotZero(s.merchantID)

	// Create Document
	createDocReq := &requests.CreateMerchantDocumentRequest{
		MerchantID:   s.merchantID,
		DocumentType: "Identity Proof",
		DocumentUrl:  "http://example.com/doc.pdf",
	}
	res, err := s.merchantService.MerchantDocumentCommandService().CreateMerchantDocument(ctx, createDocReq)
	s.NoError(err)
	s.NotNil(res)
	s.documentID = int(res.DocumentID)

	// FindByIdDocument
	found, err := s.merchantService.MerchantDocumentQueryService().FindById(ctx, s.documentID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.documentID), found.DocumentID)

	// Update Document
	updateDocReq := &requests.UpdateMerchantDocumentRequest{
		DocumentID:   &s.documentID,
		MerchantID:   s.merchantID,
		DocumentType: "Updated Type",
		DocumentUrl:  "http://example.com/updated.pdf",
		Status:       "pending",
		Note:         "Please re-upload",
	}
	updated, err := s.merchantService.MerchantDocumentCommandService().UpdateMerchantDocument(ctx, updateDocReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Type", updated.DocumentType)

	// Update Document Status
	statusDocReq := &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &s.documentID,
		MerchantID: s.merchantID,
		Status:     "approved",
		Note:       "All good",
	}
	statusUpdated, err := s.merchantService.MerchantDocumentCommandService().UpdateMerchantDocumentStatus(ctx, statusDocReq)
	s.NoError(err)
	s.NotNil(statusUpdated)
	s.Equal("approved", statusUpdated.Status)
}

func (s *MerchantServiceTestSuite) Test3_TrashAndRestore() {
	ctx := context.Background()
	s.Require().NotZero(s.merchantID)
	s.Require().NotZero(s.documentID)

	// Trash Document
	trashedDoc, err := s.merchantService.MerchantDocumentCommandService().TrashedMerchantDocument(ctx, s.documentID)
	s.NoError(err)
	s.NotNil(trashedDoc)

	// Restore Document
	restoredDoc, err := s.merchantService.MerchantDocumentCommandService().RestoreMerchantDocument(ctx, s.documentID)
	s.NoError(err)
	s.NotNil(restoredDoc)

	// Trash Merchant
	trashedMerchant, err := s.merchantService.MerchantCommandService().TrashedMerchant(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(trashedMerchant)

	// Restore Merchant
	restoredMerchant, err := s.merchantService.MerchantCommandService().RestoreMerchant(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(restoredMerchant)
}

func (s *MerchantServiceTestSuite) Test4_BulkOperations() {
	ctx := context.Background()

	// Restore All
	ok, err := s.merchantService.MerchantCommandService().RestoreAllMerchant(ctx)
	s.NoError(err)
	s.True(ok)

	ok, err = s.merchantService.MerchantDocumentCommandService().RestoreAllMerchantDocument(ctx)
	s.NoError(err)
	s.True(ok)

	// Delete All Permanent
	ok, err = s.merchantService.MerchantCommandService().DeleteAllMerchantPermanent(ctx)
	s.NoError(err)
	s.True(ok)

	ok, err = s.merchantService.MerchantDocumentCommandService().DeleteAllMerchantDocumentPermanent(ctx)
	s.NoError(err)
	s.True(ok)
}

func TestMerchantServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantServiceTestSuite))
}

package transfer_test

import (
	"context"
	"testing"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-transfer/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TransferRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	commandRepo   repository.TransferCommandRepository
	queryRepo     repository.TransferQueryRepository
	userRepo    user_repo.UserCommandRepository
	cardRepo    card_repo.Repositories
	saldoRepo   saldo_repo.Repositories

	senderCardNumber   string
	receiverCardNumber string
	transferID         int
}

func (s *TransferRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	
	// Repositories for seeding
	s.userRepo = user_repo.NewUserCommandRepository(queries)
	s.cardRepo = *card_repo.NewRepositories(queries)
	s.saldoRepo = saldo_repo.NewRepositories(queries)

	s.commandRepo = repository.NewTransferCommandRepository(queries)
	s.queryRepo = repository.NewTransferQueryRepository(queries)
}

func (s *TransferRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *TransferRepositoryTestSuite) Test1_CreateTransfer() {
	ctx := context.Background()

	// Seed Sender
	sender, err := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Sender",
		LastName:  "Repo",
		Email:     "sender.repo@test.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	sCard, err := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(sender.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "111",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.senderCardNumber = sCard.CardNumber

	// Seed Receiver
	receiver, err := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Receiver",
		LastName:  "Repo",
		Email:     "receiver.repo@test.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	rCard, err := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(receiver.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "222",
		CardProvider: "mastercard",
	})
	s.Require().NoError(err)
	s.receiverCardNumber = rCard.CardNumber

	req := &requests.CreateTransferRequest{
		TransferFrom:   s.senderCardNumber,
		TransferTo:     s.receiverCardNumber,
		TransferAmount: 25000,
	}

	transfer, err := s.commandRepo.CreateTransfer(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(transfer)
	s.transferID = int(transfer.TransferID)
	s.Equal(int32(req.TransferAmount), transfer.TransferAmount)
}

func (s *TransferRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.transferID)
	ctx := context.Background()

	res, err := s.queryRepo.FindById(ctx, s.transferID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal(int32(s.transferID), res.TransferID)
}

func TestTransferRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferRepositoryTestSuite))
}

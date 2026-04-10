package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	transactionstatsrepository "github.com/MamangRust/monolith-payment-gateway-transaction/repository/stats"
	transactionbycardrepository "github.com/MamangRust/monolith-payment-gateway-transaction/repository/statsbycard"
)

type Repositories interface {
	SaldoRepository
	MerchantRepository
	CardRepository
	TransactionQueryRepository
	TransactionCommandRepository
	transactionstatsrepository.TransactionStatsRepository
	transactionbycardrepository.TransactionStatsByCardRepository
}

// Repositories is a struct that contains all the repositories used in the transaction service.
type repositories struct {
	SaldoRepository
	MerchantRepository
	CardRepository
	TransactionQueryRepository
	TransactionCommandRepository
	transactionstatsrepository.TransactionStatsRepository
	transactionbycardrepository.TransactionStatsByCardRepository
}

// NewRepositories creates a new instance of Repositories with the provided database
// queries, context, and record mappers. This repository is responsible for
// executing command and query operations related to topup records in the database.
//
// Parameters:
//   - deps: A pointer to Deps containing the required dependencies.
//
// Returns:
//   - A pointer to the newly created Repositories instance.
func NewRepositories(
	db *db.Queries,
	saldo SaldoRepository,
	card CardRepository,
	merchant MerchantRepository,
) Repositories {

	return &repositories{
		SaldoRepository:                  saldo,
		MerchantRepository:               merchant,
		CardRepository:                   card,
		TransactionQueryRepository:       NewTransactionQueryRepository(db),
		TransactionCommandRepository:     NewTransactionCommandRepository(db),
		TransactionStatsRepository:       transactionstatsrepository.NewTransactionStatsRepository(db),
		TransactionStatsByCardRepository: transactionbycardrepository.NewTransactionStatsRepository(db),
	}
}

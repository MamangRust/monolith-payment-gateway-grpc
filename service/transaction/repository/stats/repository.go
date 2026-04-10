package transactionstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type TransactionStatsRepository interface {
	TransactionStatsAmountRepository
	TransactionStatsMethodRepository
	TransactionStatsStatusRepository
}

type repository struct {
	TransactionStatsAmountRepository
	TransactionStatsMethodRepository
	TransactionStatsStatusRepository
}

func NewTransactionStatsRepository(db *db.Queries) TransactionStatsRepository {

	return &repository{
		TransactionStatsAmountRepository: NewTransactionStatsAmountRepository(db),
		TransactionStatsMethodRepository: NewTransactionStatsMethodRepository(db),
		TransactionStatsStatusRepository: NewTransactionStatsStatusRepository(db),
	}
}

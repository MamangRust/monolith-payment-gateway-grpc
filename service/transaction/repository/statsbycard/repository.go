package transactionbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type TransactionStatsByCardRepository interface {
	TransactionStatsByCardAmountRepository
	TransactionStatsByCardMethodRepository
	TransactionStatsByCardStatusRepository
}

type repository struct {
	TransactionStatsByCardAmountRepository
	TransactionStatsByCardMethodRepository
	TransactionStatsByCardStatusRepository
}

func NewTransactionStatsRepository(db *db.Queries) TransactionStatsByCardRepository {

	return &repository{
		TransactionStatsByCardAmountRepository:  NewTransactionStatsByCardAmountRepository(db),
		TransactionStatsByCardMethodRepository: NewTransactionStatsByCardMethodRepository(db),
		TransactionStatsByCardStatusRepository:  NewTransactionStatsByCardStatusRepository(db),
	}
}


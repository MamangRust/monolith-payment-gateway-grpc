package transactionbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type TransactionStatsByCardRepository interface {
	TransactonStatsByCardAmountRepository
	TransactionStatsByCardMethodRepository
	TransactonStatsByCardStatusRepository
}

type repository struct {
	TransactonStatsByCardAmountRepository
	TransactionStatsByCardMethodRepository
	TransactonStatsByCardStatusRepository
}

func NewTransactionStatsRepository(db *db.Queries) TransactionStatsByCardRepository {

	return &repository{
		TransactonStatsByCardAmountRepository:  NewTransactionStatsByCardAmountRepository(db),
		TransactionStatsByCardMethodRepository: NewTransactionStatsByCardMethodRepository(db),
		TransactonStatsByCardStatusRepository:  NewTransactionStatsByCardStatusRepository(db),
	}
}

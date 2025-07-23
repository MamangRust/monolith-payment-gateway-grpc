package transactionbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/statsbycard"
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

func NewTransactionStatsRepository(db *db.Queries, mapper recordmapper.TransactonStatisticByCardMapper) TransactionStatsByCardRepository {

	return &repository{
		TransactonStatsByCardAmountRepository:  NewTransactionStatsByCardAmountRepository(db, mapper),
		TransactionStatsByCardMethodRepository: NewTransactionStatsByCardMethodRepository(db, mapper),
		TransactonStatsByCardStatusRepository:  NewTransactionStatsByCardStatusRepository(db, mapper),
	}
}

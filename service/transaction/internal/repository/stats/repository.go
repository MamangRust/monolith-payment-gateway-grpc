package transactionstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/stats"
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

func NewTransactionStatsRepository(db *db.Queries, mapper recordmapper.TransactonStatisticsRecordMapper) TransactionStatsRepository {

	return &repository{
		TransactionStatsAmountRepository: NewTransactionStatsAmountRepository(db, mapper),
		TransactionStatsMethodRepository: NewTransactionStatsMethodRepository(db, mapper),
		TransactionStatsStatusRepository: NewTransactionStatsStatusRepository(db, mapper),
	}
}

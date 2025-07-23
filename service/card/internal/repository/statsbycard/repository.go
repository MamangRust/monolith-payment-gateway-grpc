package repositorystatsbycard

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/statsbycard"
)

type CardStatsByCardRepository interface {
	CardStatsBalanceByCardRepository
	CardStatsTopupByCardRepository
	CardStatsTransactionByCardRepository
	CardStatsTransferByCardRepository
	CardStatsWithdrawByCardRepository
}

type repository struct {
	CardStatsBalanceByCardRepository
	CardStatsTopupByCardRepository
	CardStatsTransactionByCardRepository
	CardStatsTransferByCardRepository
	CardStatsWithdrawByCardRepository
}

func NewCardStatsByCardRepository(db *db.Queries, mapper recordmapper.CardStatsByCardRecordMapper) CardStatsByCardRepository {

	return &repository{
		CardStatsBalanceByCardRepository:     NewCardStatsBalanceByCardRepository(db, mapper),
		CardStatsTopupByCardRepository:       NewCardStatsTopupByCardRepository(db, mapper),
		CardStatsTransactionByCardRepository: NewCardStatsTransactionByCardRepository(db, mapper),
		CardStatsTransferByCardRepository:    NewCardStatsTransferByCardRepository(db, mapper),
		CardStatsWithdrawByCardRepository:    NewCardStatsWithdrawByCardRepository(db, mapper),
	}
}

package repositorystats

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/stats"
)

type CardStatsRepository interface {
	CardStatsBalanceRepository
	CardStatsTopupRepository
	CardStatsTransactionRepository
	CardStatsTransferRepository
	CardStatsWithdrawRepository
}

type repository struct {
	CardStatsBalanceRepository
	CardStatsTopupRepository
	CardStatsTransactionRepository
	CardStatsTransferRepository
	CardStatsWithdrawRepository
}

func NewCardStatsRepository(db *db.Queries, mapper recordmapper.CardStatsRecordMapper) CardStatsRepository {

	return &repository{
		CardStatsBalanceRepository:     NewCardStatsBalanceRepository(db, mapper),
		CardStatsTopupRepository:       NewCardStatsTopupRepository(db, mapper),
		CardStatsTransactionRepository: NewCardStatsTransactionRepository(db, mapper),
		CardStatsTransferRepository:    NewCardStatsTransferRepository(db, mapper),
		CardStatsWithdrawRepository:    NewCardStatsWithdrawRepository(db, mapper),
	}
}

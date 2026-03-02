package repositorystats

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewCardStatsRepository(db *db.Queries) CardStatsRepository {

	return &repository{
		CardStatsBalanceRepository:     NewCardStatsBalanceRepository(db),
		CardStatsTopupRepository:       NewCardStatsTopupRepository(db),
		CardStatsTransactionRepository: NewCardStatsTransactionRepository(db),
		CardStatsTransferRepository:    NewCardStatsTransferRepository(db),
		CardStatsWithdrawRepository:    NewCardStatsWithdrawRepository(db),
	}
}

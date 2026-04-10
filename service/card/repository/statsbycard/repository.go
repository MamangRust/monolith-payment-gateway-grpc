package repositorystatsbycard

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewCardStatsByCardRepository(db *db.Queries) CardStatsByCardRepository {
	return &repository{
		CardStatsBalanceByCardRepository:     NewCardStatsBalanceByCardRepository(db),
		CardStatsTopupByCardRepository:       NewCardStatsTopupByCardRepository(db),
		CardStatsTransactionByCardRepository: NewCardStatsTransactionByCardRepository(db),
		CardStatsTransferByCardRepository:    NewCardStatsTransferByCardRepository(db),
		CardStatsWithdrawByCardRepository:    NewCardStatsWithdrawByCardRepository(db),
	}
}

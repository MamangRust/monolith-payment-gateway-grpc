package saldostatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type SaldoStatsRepository interface {
	SaldoStatsBalanceRepository
	SaldoStatsTotalSaldoRepository
}

type repository struct {
	SaldoStatsBalanceRepository
	SaldoStatsTotalSaldoRepository
}

func NewSaldoStatsRepository(db *db.Queries) SaldoStatsRepository {
	return &repository{
		SaldoStatsBalanceRepository:    NewSaldoStatsBalanceRepository(db),
		SaldoStatsTotalSaldoRepository: NewSaldoStatsTotalBalanceRepository(db),
	}
}

package saldostatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

type SaldoStatsRepository interface {
	SaldoStatsBalanceRepository
	SaldoStatsTotalSaldoRepository
}

type repository struct {
	SaldoStatsBalanceRepository
	SaldoStatsTotalSaldoRepository
}

func NewSaldoStatsRepository(db *db.Queries, mapper recordmapper.SaldoStatisticRecordMapping) SaldoStatsRepository {

	return &repository{
		SaldoStatsBalanceRepository:    NewSaldoStatsBalanceRepository(db, mapper),
		SaldoStatsTotalSaldoRepository: NewSaldoStatsTotalBalanceRepository(db, mapper),
	}
}

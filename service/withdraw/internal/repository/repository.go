package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	mappercard "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
	mappersaldo "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
	mapperwithdraw "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw"
	withdrawstatsrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository/stats"
	withdrawstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository/statsbycard"
)

type Repositories interface {
	CardRepository
	SaldoRepository
	WithdrawQueryRepository
	WithdrawCommandRepository
	withdrawstatsrepository.WithdrawStatsRepository
	withdrawstatsbycardrepository.WithdrawStatsByCardRepository
}

type repositories struct {
	CardRepository
	SaldoRepository
	WithdrawQueryRepository
	WithdrawCommandRepository
	withdrawstatsrepository.WithdrawStatsRepository
	withdrawstatsbycardrepository.WithdrawStatsByCardRepository
}

func NewRepositories(db *db.Queries) Repositories {
	mapperwithdraw := mapperwithdraw.NewWithdrawRecordMapper()
	mappersaldo := mappersaldo.NewSaldoQueryRecordMapper()
	mappercard := mappercard.NewCardQueryRecordMapper()

	return &repositories{
		CardRepository:                NewCardRepository(db, mappercard),
		SaldoRepository:               NewSaldoRepository(db, mappersaldo),
		WithdrawQueryRepository:       NewWithdrawQueryRepository(db, mapperwithdraw.QueryMapper()),
		WithdrawCommandRepository:     NewWithdrawCommandRepository(db, mapperwithdraw.CommandMapper()),
		WithdrawStatsRepository:       withdrawstatsrepository.NewWithdrawStatsRepository(db, mapperwithdraw.StatsMapper()),
		WithdrawStatsByCardRepository: withdrawstatsbycardrepository.NewWithdrawStatsByCardRepository(db, mapperwithdraw.StatsByCardMapper()),
	}
}

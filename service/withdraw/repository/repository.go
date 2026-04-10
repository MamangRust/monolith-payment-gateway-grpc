package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	withdrawstatsrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/repository/stats"
	withdrawstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/repository/statsbycard"
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

func NewRepositories(
	db *db.Queries,
	card CardRepository,
	saldo SaldoRepository,
) Repositories {
	return &repositories{
		CardRepository:                card,
		SaldoRepository:               saldo,
		WithdrawQueryRepository:       NewWithdrawQueryRepository(db),
		WithdrawCommandRepository:     NewWithdrawCommandRepository(db),
		WithdrawStatsRepository:       withdrawstatsrepository.NewWithdrawStatsRepository(db),
		WithdrawStatsByCardRepository: withdrawstatsbycardrepository.NewWithdrawStatsByCardRepository(db),
	}
}

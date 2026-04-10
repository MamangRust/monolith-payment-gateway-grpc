package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	saldostatsrepository "github.com/MamangRust/monolith-payment-gateway-saldo/repository/stats"
)

type Repositories interface {
	SaldoQueryRepository
	SaldoCommandRepository
	saldostatsrepository.SaldoStatsRepository
	CardRepository
}

type repositories struct {
	SaldoQueryRepository
	SaldoCommandRepository
	saldostatsrepository.SaldoStatsRepository
	CardRepository
}

func NewRepositories(db *db.Queries) Repositories {
	return &repositories{
		SaldoQueryRepository:   NewSaldoQueryRepository(db),
		SaldoCommandRepository: NewSaldoCommandRepository(db),
		SaldoStatsRepository:   saldostatsrepository.NewSaldoStatsRepository(db),
		CardRepository:         NewCardRepository(db),
	}
}

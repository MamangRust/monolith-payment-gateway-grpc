package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	topupstatsrepository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/stats"
	topupstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/statsbycard"
)

type Repositories interface {
	TopupQueryRepository
	TopupCommandRepository
	CardRepository
	SaldoRepository
	topupstatsrepository.TopupStatsRepository
	topupstatsbycardrepository.TopupStatsByCardRepository
}

type repositories struct {
	TopupQueryRepository
	TopupCommandRepository
	CardRepository
	SaldoRepository
	topupstatsrepository.TopupStatsRepository
	topupstatsbycardrepository.TopupStatsByCardRepository
}

func NewRepositories(db *db.Queries) Repositories {
	return &repositories{
		TopupQueryRepository:       NewTopupQueryRepository(db),
		TopupCommandRepository:     NewTopupCommandRepository(db),
		TopupStatsRepository:       topupstatsrepository.NewTopupStatsRepository(db),
		TopupStatsByCardRepository: topupstatsbycardrepository.NewTopupStatsByCardRepository(db),
		CardRepository:             NewCardRepository(db),
		SaldoRepository:            NewSaldoRepository(db),
	}
}

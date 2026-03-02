package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	transferstatsrepository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/stats"
	transferstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/statsbycard"
)

type Repositories interface {
	SaldoRepository
	TransferQueryRepository
	TransferCommandRepository
	CardRepository
	transferstatsrepository.TransferStatsRepository
	transferstatsbycardrepository.TransferStatsByCardRepository
}

type repositories struct {
	SaldoRepository
	TransferQueryRepository
	TransferCommandRepository
	CardRepository
	transferstatsrepository.TransferStatsRepository
	transferstatsbycardrepository.TransferStatsByCardRepository
}

func NewRepositories(db *db.Queries) Repositories {
	return &repositories{
		SaldoRepository:               NewSaldoRepository(db),
		TransferQueryRepository:       NewTransferQueryRepository(db),
		TransferCommandRepository:     NewTransferCommandRepository(db),
		TransferStatsRepository:       transferstatsrepository.NewTransferStatsRepository(db),
		TransferStatsByCardRepository: transferstatsbycardrepository.NewTransferStatsByCardRepository(db),
		CardRepository:                NewCardRepository(db),
	}
}

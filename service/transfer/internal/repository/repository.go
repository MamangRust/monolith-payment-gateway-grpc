package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	mappercard "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
	mappersaldo "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
	mapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer"
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

// NewRepositories creates a new instance of the Repositories interface, which contains all the repositories used in the transfer service.
//
// Parameters:
//   - db: A pointer to the database queries.
//
// Returns:
//   - A pointer to the newly created Repositories instance.
func NewRepositories(db *db.Queries) Repositories {
	saldoMapper := mappersaldo.NewSaldoQueryRecordMapper()
	transferMapper := mapper.NewTransferRecordMapper()
	cardmapper := mappercard.NewCardQueryRecordMapper()

	return &repositories{
		SaldoRepository:               NewSaldoRepository(db, saldoMapper),
		TransferQueryRepository:       NewTransferQueryRepository(db, transferMapper.QueryMapper()),
		TransferCommandRepository:     NewTransferCommandRepository(db, transferMapper.CommandMapper()),
		TransferStatsRepository:       transferstatsrepository.NewTransferStatsRepository(db, transferMapper.StatsMapper()),
		TransferStatsByCardRepository: transferstatsbycardrepository.NewTransferStatsByCardRepository(db, transferMapper.StatsByCardMapper()),
		CardRepository:                NewCardRepository(db, cardmapper),
	}
}

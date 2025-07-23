package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	mappercard "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
	mappersaldo "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup"
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

// Repositories is a struct that contains all the repositories for the topup service
type repositories struct {
	TopupQueryRepository
	TopupCommandRepository
	CardRepository
	SaldoRepository
	topupstatsrepository.TopupStatsRepository
	topupstatsbycardrepository.TopupStatsByCardRepository
}

// NewRepositories creates a new instance of Repositories with the provided database
// queries, context, and record mappers. This repository is responsible for
// executing command and query operations related to topup records in the database.
//
// Parameters:
//   - deps: A pointer to Deps containing the required dependencies.
//
// Returns:
//   - A pointer to the newly created Repositories instance.
func NewRepositories(db *db.Queries) Repositories {
	mapper := recordmapper.NewTopupRecordMapper()
	mappersaldo := mappersaldo.NewSaldoQueryRecordMapper()
	mappercard := mappercard.NewCardQueryRecordMapper()

	return &repositories{
		TopupQueryRepository:       NewTopupQueryRepository(db, mapper.QueryMapper()),
		TopupCommandRepository:     NewTopupCommandRepository(db, mapper.CommandMapper()),
		TopupStatsRepository:       topupstatsrepository.NewTopupStatsRepository(db, mapper.StatsMapper()),
		TopupStatsByCardRepository: topupstatsbycardrepository.NewTopupStatsByCardRepository(db, mapper.StatsByCardMapper()),
		CardRepository:             NewCardRepository(db, mappercard),
		SaldoRepository:            NewSaldoRepository(db, mappersaldo),
	}
}

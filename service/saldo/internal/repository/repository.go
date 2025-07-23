package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	saldostatsrepository "github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository/stats"
	mappercard "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
	mappersaldo "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

type Repositories interface {
	SaldoQueryRepository
	SaldoCommandRepository
	saldostatsrepository.SaldoStatsRepository
	CardRepository
}

// Repositories is a struct that contains the repositories for the saldo service
type repositories struct {
	SaldoQueryRepository
	SaldoCommandRepository
	saldostatsrepository.SaldoStatsRepository
	CardRepository
}

// NewRepositories creates a new instance of Repositories with the provided database
// queries, context, and record mappers. This repository is responsible for
// executing command and query operations related to saldo records in the database.
//
// Parameters:
//   - deps: A pointer to Deps containing the required dependencies.
//
// Returns:
//   - A pointer to the newly created Repositories instance.
func NewRepositories(db *db.Queries) Repositories {
	cardRecordMapper := mappercard.NewCardQueryRecordMapper()
	saldoRecordMapper := mappersaldo.NewSaldoRecordMapper()

	return &repositories{
		SaldoQueryRepository:   NewSaldoQueryRepository(db, saldoRecordMapper.QueryMapper()),
		SaldoCommandRepository: NewSaldoCommandRepository(db, saldoRecordMapper.CommandMapper()),
		SaldoStatsRepository:   saldostatsrepository.NewSaldoStatsRepository(db, saldoRecordMapper.StatisticMapper()),
		CardRepository:         NewCardRepository(db, cardRecordMapper),
	}
}

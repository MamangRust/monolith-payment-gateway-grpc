package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	SaldoQuery   SaldoQueryRepository
	SaldoCommand SaldoCommandRepository
	SaldoStats   SaldoStatisticsRepository

	Card CardRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps *Deps) *Repositories {
	cardRecordMapper := recordmapper.NewCardRecordMapper()
	saldoRecordMapper := recordmapper.NewSaldoRecordMapper()

	return &Repositories{
		SaldoQuery:   NewSaldoQueryRepository(deps.DB, deps.Ctx, saldoRecordMapper),
		SaldoCommand: NewSaldoCommandRepository(deps.DB, deps.Ctx, saldoRecordMapper),
		SaldoStats:   NewSaldoStatisticsRepository(deps.DB, deps.Ctx, saldoRecordMapper),
		Card:         NewCardRepository(deps.DB, deps.Ctx, cardRecordMapper),
	}
}

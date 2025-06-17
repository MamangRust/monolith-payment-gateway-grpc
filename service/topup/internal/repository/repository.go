package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	TopupQuery           TopupQueryRepository
	TopupStatistic       TopupStatisticRepository
	TopupStatistisByCard TopupStatisticByCardRepository
	TopupCommand         TopupCommandRepository

	Card  CardRepository
	Saldo SaldoRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		TopupQuery:           NewTopupQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.TopupRecordMapper),
		TopupStatistic:       NewTopupStatisticRepository(deps.DB, deps.Ctx, deps.MapperRecord.TopupRecordMapper),
		TopupStatistisByCard: NewTopupStatisticByCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.TopupRecordMapper),
		TopupCommand:         NewTopupCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.TopupRecordMapper),
		Saldo:                NewSaldoRepository(deps.DB, deps.Ctx, deps.MapperRecord.SaldoRecordMapper),
		Card:                 NewCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
	}
}

package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	Card                CardRepository
	Saldo               SaldoRepository
	WithdrawQuery       WithdrawQueryRepository
	WithdrawStats       WithdrawStatisticRepository
	WithdrawCommand     WithdrawCommandRepository
	WIthdrawStatsByCard WithdrawStatisticByCardRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		Card:                NewCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		Saldo:               NewSaldoRepository(deps.DB, deps.Ctx, deps.MapperRecord.SaldoRecordMapper),
		WithdrawQuery:       NewWithdrawQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.WithdrawRecordMapper),
		WithdrawStats:       NewWithdrawStatisticRepository(deps.DB, deps.Ctx, deps.MapperRecord.WithdrawRecordMapper),
		WithdrawCommand:     NewWithdrawCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.WithdrawRecordMapper),
		WIthdrawStatsByCard: NewWithdrawStatisticByCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.WithdrawRecordMapper),
	}
}

package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	Saldo               SaldoRepository
	TransferQuery       TransferQueryRepository
	TransferStats       TransferStatisticRepository
	TransferStatsByCard TransferStatisticByCardRepository
	TransferCommand     TransferCommandRepository
	Card                CardRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		Saldo:               NewSaldoRepository(deps.DB, deps.Ctx, deps.MapperRecord.SaldoRecordMapper),
		TransferQuery:       NewTransferQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransferRecordMapper),
		TransferStats:       NewTransferStatisticRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransferRecordMapper),
		TransferStatsByCard: NewTransferStatisticByCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransferRecordMapper),
		TransferCommand:     NewTransferCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransferRecordMapper),
		Card:                NewCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
	}
}

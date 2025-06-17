package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	Saldo                 SaldoRepository
	Merchant              MerchantRepository
	Card                  CardRepository
	TransactionQuery      TransactionQueryRepository
	TransactionStat       TransactionStatisticsRepository
	TransactionStatByCard TransactionStatisticByCardRepository
	TransactionCommand    TransactionCommandRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		Saldo:                 NewSaldoRepository(deps.DB, deps.Ctx, deps.MapperRecord.SaldoRecordMapper),
		Merchant:              NewMerchantRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		Card:                  NewCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		TransactionQuery:      NewTransactionQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransactionRecordMapper),
		TransactionStat:       NewTransactionStatisticRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransactionRecordMapper),
		TransactionStatByCard: NewTransactionStatisticByCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransactionRecordMapper),
		TransactionCommand:    NewTransactionCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransactionRecordMapper),
	}
}

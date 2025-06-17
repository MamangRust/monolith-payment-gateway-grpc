package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantTrans           MerchantTransactionRepository
	MerchantStat            MerchantStatisticRepository
	MerchantStatByApiKey    MerchantStatisticByApiKeyRepository
	MerchantStatByMerchant  MerchantStatisticByMerchantRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	User                    UserRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps *Deps) *Repositories {

	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantTrans:           NewMerchantTransactionRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantStat:            NewMerchantStatisticRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantStatByApiKey:    NewMerchantStatisticByApiKeyRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantStatByMerchant:  NewMerchantStatisticByMerchantRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantCommand:         NewMerchantCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantDocumentRecordMapper),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantDocumentRecordMapper),
		User:                    NewUserRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
	}
}

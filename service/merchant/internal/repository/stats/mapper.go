package merchantstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/stats"
)

type MerchantStatsRepository interface {
	MerchantStatsAmountRepository
	MerchantStatsMethodRepository
	MerchantStatsTotalAmountRepository
}

type repository struct {
	MerchantStatsAmountRepository
	MerchantStatsMethodRepository
	MerchantStatsTotalAmountRepository
}

func NewMerchantStatsRepository(db *db.Queries, mapper recordmapper.MerchantStatisticRecordMapper) MerchantStatsRepository {

	return &repository{
		MerchantStatsAmountRepository:      NewMerchantStatsAmountRepository(db, mapper),
		MerchantStatsMethodRepository:      NewMerchantStatsMethodRepository(db, mapper),
		MerchantStatsTotalAmountRepository: NewMerchantStatsTotalAmountRepository(db, mapper),
	}
}

package merchantstatsmerchantrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/statsByMerchant"
)

type MerchantStatsByMerchantRepository interface {
	MerchantStatsAmountByMerchantRepository
	MerchantStatsMethodByMerchantRepository
	MerchantStatsTotalAmountByMerchantRepository
}

type repository struct {
	MerchantStatsAmountByMerchantRepository
	MerchantStatsMethodByMerchantRepository
	MerchantStatsTotalAmountByMerchantRepository
}

func NewMerchantStatsByMerchantRepository(db *db.Queries, mapper recordmapper.MerchantStatisticByMerchantMapper) MerchantStatsByMerchantRepository {

	return &repository{
		MerchantStatsAmountByMerchantRepository:      NewMerchantStatsAmountByMerchantRepository(db, mapper),
		MerchantStatsMethodByMerchantRepository:      NewMerchantStatsMethodByMerchantRepository(db, mapper),
		MerchantStatsTotalAmountByMerchantRepository: NewMerchantStatsTotalAmountByMerchantRepository(db, mapper),
	}
}

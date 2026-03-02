package merchantstatsmerchantrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewMerchantStatsByMerchantRepository(db *db.Queries) MerchantStatsByMerchantRepository {
	return &repository{
		MerchantStatsAmountByMerchantRepository:      NewMerchantStatsAmountByMerchantRepository(db),
		MerchantStatsMethodByMerchantRepository:      NewMerchantStatsMethodByMerchantRepository(db),
		MerchantStatsTotalAmountByMerchantRepository: NewMerchantStatsTotalAmountByMerchantRepository(db),
	}
}

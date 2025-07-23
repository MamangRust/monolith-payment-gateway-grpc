package merchantstatsapikeyrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/statsByApiKey"
)

type MerchantStatsByApiKeyRepository interface {
	MerchantStatsAmountByApiKeyRepository
	MerchantStatsMethodByApiKeyRepository
	MerchantStatsTotalAmountByApiKeyRepository
}

type repository struct {
	MerchantStatsAmountByApiKeyRepository
	MerchantStatsMethodByApiKeyRepository
	MerchantStatsTotalAmountByApiKeyRepository
}

func NewMerchantStatsByApiKeyRepository(db *db.Queries, mapper recordmapper.MerchantStatisticByApiKeyMapper) MerchantStatsByApiKeyRepository {

	return &repository{
		MerchantStatsAmountByApiKeyRepository:      NewMerchantStatsAmountByApiKeyRepository(db, mapper),
		MerchantStatsMethodByApiKeyRepository:      NewMerchantStatsMethodByApiKeyRepository(db, mapper),
		MerchantStatsTotalAmountByApiKeyRepository: NewMerchantStatsTotalAmountByApiKeyRepository(db, mapper),
	}
}

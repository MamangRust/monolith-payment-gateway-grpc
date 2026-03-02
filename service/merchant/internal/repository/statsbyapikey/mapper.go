package merchantstatsapikeyrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewMerchantStatsByApiKeyRepository(db *db.Queries) MerchantStatsByApiKeyRepository {

	return &repository{
		MerchantStatsAmountByApiKeyRepository:      NewMerchantStatsAmountByApiKeyRepository(db),
		MerchantStatsMethodByApiKeyRepository:      NewMerchantStatsMethodByApiKeyRepository(db),
		MerchantStatsTotalAmountByApiKeyRepository: NewMerchantStatsTotalAmountByApiKeyRepository(db),
	}
}

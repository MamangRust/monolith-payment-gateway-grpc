package merchantstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewMerchantStatsRepository(db *db.Queries) MerchantStatsRepository {
	return &repository{
		MerchantStatsAmountRepository:      NewMerchantStatsAmountRepository(db),
		MerchantStatsMethodRepository:      NewMerchantStatsMethodRepository(db),
		MerchantStatsTotalAmountRepository: NewMerchantStatsTotalAmountRepository(db),
	}
}

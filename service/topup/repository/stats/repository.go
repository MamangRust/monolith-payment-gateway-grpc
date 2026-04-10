package topupstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type TopupStatsRepository interface {
	TopupStatsAmountRepository
	TopupStatsStatusRepository
	TopupStatsMethodRepository
}

type repository struct {
	TopupStatsAmountRepository
	TopupStatsStatusRepository
	TopupStatsMethodRepository
}

func NewTopupStatsRepository(db *db.Queries) TopupStatsRepository {
	return &repository{
		TopupStatsAmountRepository: NewTopupStatsAmountRepository(db),
		TopupStatsStatusRepository: NewTopupStatsStatusRepository(db),
		TopupStatsMethodRepository: NewTopupStatsMethodRepository(db),
	}
}

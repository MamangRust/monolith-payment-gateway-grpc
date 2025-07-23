package topupstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/stats"
)

type TopupStatsRepository interface {
	TopupStatsAmountRepository
	TopupStatsStatusRepository
	TOpupStatsMethodRepository
}

type repository struct {
	TopupStatsAmountRepository
	TopupStatsStatusRepository
	TOpupStatsMethodRepository
}

func NewTopupStatsRepository(db *db.Queries, mapper recordmapper.TopupStatisticRecordMapper) TopupStatsRepository {

	return &repository{
		TopupStatsAmountRepository: NewTopupStatsAmountRepository(db, mapper),
		TopupStatsStatusRepository: NewTopupStatsStatusRepository(db, mapper),
		TOpupStatsMethodRepository: NewTopupStatsMethodRepository(db, mapper),
	}
}

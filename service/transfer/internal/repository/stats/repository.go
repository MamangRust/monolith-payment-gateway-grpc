package transferstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer/stats"
)

type TransferStatsRepository interface {
	TransferStatsAmountRepository
	TransferStatsStatusRepository
}

type repositories struct {
	TransferStatsAmountRepository
	TransferStatsStatusRepository
}

func NewTransferStatsRepository(db *db.Queries, mapper recordmapper.TransferStatisticRecordMapper) TransferStatsRepository {

	return &repositories{
		TransferStatsAmountRepository: NewTransferStatsAmountRepository(db, mapper),
		TransferStatsStatusRepository: NewTransferStatsStatusRepository(db, mapper),
	}
}

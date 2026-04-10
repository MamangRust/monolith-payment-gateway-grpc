package transferstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type TransferStatsRepository interface {
	TransferStatsAmountRepository
	TransferStatsStatusRepository
}

type repositories struct {
	TransferStatsAmountRepository
	TransferStatsStatusRepository
}

func NewTransferStatsRepository(db *db.Queries) TransferStatsRepository {

	return &repositories{
		TransferStatsAmountRepository: NewTransferStatsAmountRepository(db),
		TransferStatsStatusRepository: NewTransferStatsStatusRepository(db),
	}
}

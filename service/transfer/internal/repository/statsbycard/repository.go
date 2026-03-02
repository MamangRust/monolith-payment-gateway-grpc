package transferstatsbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type TransferStatsByCardRepository interface {
	TransferStatsByCardAmountSenderRepository
	TransferStatsByCardAmountReceiverRepository
	TransferStatsByCardStatusRepository
}

type repositories struct {
	TransferStatsByCardAmountSenderRepository
	TransferStatsByCardAmountReceiverRepository
	TransferStatsByCardStatusRepository
}

func NewTransferStatsByCardRepository(db *db.Queries) TransferStatsByCardRepository {

	return &repositories{
		TransferStatsByCardAmountSenderRepository:   NewTransferStatsAmountSenderRepository(db),
		TransferStatsByCardAmountReceiverRepository: NewTransferStatsAmountReceiverRepository(db),
		TransferStatsByCardStatusRepository:         NewTransferStatsByCardStatusRepository(db),
	}
}

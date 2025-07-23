package transferstatsbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer/statsbycard"
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

func NewTransferStatsByCardRepository(db *db.Queries, mapper recordmapper.TransferStatisticByCardRecordMapper) TransferStatsByCardRepository {

	return &repositories{
		TransferStatsByCardAmountSenderRepository:   NewTransferStatsAmountSenderRepository(db, mapper),
		TransferStatsByCardAmountReceiverRepository: NewTransferStatsAmountReceiverRepository(db, mapper),
		TransferStatsByCardStatusRepository:         NewTransferStatsByCardStatusRepository(db, mapper),
	}
}

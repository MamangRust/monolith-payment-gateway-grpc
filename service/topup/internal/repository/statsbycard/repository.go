package topupstatsbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/statsbycard"
)

type TopupStatsByCardRepository interface {
	TopupStatsByCardAmountRepository
	TopupStatsByCardMethodRepository
	TopupStatsByCardStatusRepository
}

type repository struct {
	TopupStatsByCardAmountRepository
	TopupStatsByCardMethodRepository
	TopupStatsByCardStatusRepository
}

func NewTopupStatsByCardRepository(db *db.Queries, mapper recordmapper.TopupStatisticByCardRecordMapper) TopupStatsByCardRepository {

	return &repository{
		TopupStatsByCardAmountRepository: NewTopupStatsByCardAmountRepository(db, mapper),
		TopupStatsByCardMethodRepository: NewTopupStatsByCardMethodRepository(db, mapper),
		TopupStatsByCardStatusRepository: NewTopupStatsByCardStatusRepository(db, mapper),
	}
}

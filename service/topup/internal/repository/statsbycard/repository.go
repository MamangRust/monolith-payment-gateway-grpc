package topupstatsbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewTopupStatsByCardRepository(db *db.Queries) TopupStatsByCardRepository {

	return &repository{
		TopupStatsByCardAmountRepository: NewTopupStatsByCardAmountRepository(db),
		TopupStatsByCardMethodRepository: NewTopupStatsByCardMethodRepository(db),
		TopupStatsByCardStatusRepository: NewTopupStatsByCardStatusRepository(db),
	}
}

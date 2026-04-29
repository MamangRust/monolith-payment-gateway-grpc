package withdrawstatsbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type WithdrawStatsByCardRepository interface {
	WithdrawStatsByCardAmountRepository
	WithdrawStatsByCardStatusRepository
}

type repositories struct {
	WithdrawStatsByCardAmountRepository
	WithdrawStatsByCardStatusRepository
}

func NewWithdrawStatsByCardRepository(db *db.Queries) WithdrawStatsByCardRepository {

	return &repositories{
		WithdrawStatsByCardAmountRepository: NewWithdrawStatsByCardAmountRepository(db),
		WithdrawStatsByCardStatusRepository: NewWithdrawStatsByCardStatusRepository(db),
	}
}


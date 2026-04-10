package withdrawstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type WithdrawStatsRepository interface {
	WithdrawStatsAmountRepository
	WithdrawStatsStatusRepository
}

type repositories struct {
	WithdrawStatsAmountRepository
	WithdrawStatsStatusRepository
}

func NewWithdrawStatsRepository(db *db.Queries) WithdrawStatsRepository {

	return &repositories{
		WithdrawStatsAmountRepository: NewWithdrawStatsAmountRepository(db),
		WithdrawStatsStatusRepository: NewWithdrawStatsStatusRepository(db),
	}
}

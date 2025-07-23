package withdrawstatsrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw/stats"
)

type WithdrawStatsRepository interface {
	WithdrawStatsAmountRepository
	WithdrawStatsStatusRepository
}

type repositories struct {
	WithdrawStatsAmountRepository
	WithdrawStatsStatusRepository
}

func NewWithdrawStatsRepository(db *db.Queries, mapper recordmapper.WithdrawStatisticRecordMapper) WithdrawStatsRepository {

	return &repositories{
		WithdrawStatsAmountRepository: NewWithdrawStatsAmountRepository(db, mapper),
		WithdrawStatsStatusRepository: NewWithdrawStatsStatusRepository(db, mapper),
	}
}

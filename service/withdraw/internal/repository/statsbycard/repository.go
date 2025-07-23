package withdrawstatsbycardrepository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw/statsbycard"
)

type WithdrawStatsByCardRepository interface {
	WithdrawStatsByCardAmountRepository
	WithdrawStatsByCardStatusRepository
}

type repositories struct {
	WithdrawStatsByCardAmountRepository
	WithdrawStatsByCardStatusRepository
}

func NewWithdrawStatsByCardRepository(db *db.Queries, mapper recordmapper.WithdrawStatisticByCardRecordMapper) WithdrawStatsByCardRepository {

	return &repositories{
		WithdrawStatsByCardAmountRepository: NewWithdrawStatsAmountRepository(db, mapper),
		WithdrawStatsByCardStatusRepository: NewWithdrawStatsStatusRepository(db, mapper),
	}
}

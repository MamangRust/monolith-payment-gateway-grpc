package withdrawstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw/stats"
)

type withdrawStatsAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.WithdrawStatisticAmountRecordMapper
}

func NewWithdrawStatsAmountRepository(db *db.Queries, mapper recordmapper.WithdrawStatisticAmountRecordMapper) WithdrawStatsAmountRepository {
	return &withdrawStatsAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyWithdraws retrieves the total amount of withdraws grouped by month for the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly amounts are requested.
//
// Returns:
//   - []*record.WithdrawMonthlyAmount: List of monthly withdraw amounts.
//   - error: An error if the operation fails.
func (r *withdrawStatsAmountRepository) GetMonthlyWithdraws(ctx context.Context, year int) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(ctx, yearStart)

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsFailed
	}

	so := r.mapper.ToWithdrawsAmountMonthly(res)

	return so, nil
}

// GetYearlyWithdraws retrieves the total amount of withdraws grouped by year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly amounts are requested.
//
// Returns:
//   - []*record.WithdrawYearlyAmount: List of yearly withdraw amounts.
//   - error: An error if the operation fails.
func (r *withdrawStatsAmountRepository) GetYearlyWithdraws(ctx context.Context, year int) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdraws(ctx, year)

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsFailed
	}

	so := r.mapper.ToWithdrawsAmountYearly(res)

	return so, nil
}

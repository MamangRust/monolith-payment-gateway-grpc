package withdrawstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawStatsAmountRepository struct {
	db *db.Queries
}

func NewWithdrawStatsAmountRepository(db *db.Queries) WithdrawStatsAmountRepository {
	return &withdrawStatsAmountRepository{
		db: db,
	}
}

func (r *withdrawStatsAmountRepository) GetMonthlyWithdraws(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(ctx, yearStart)

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsFailed.WithInternal(err)
	}

	return res, nil

}

func (r *withdrawStatsAmountRepository) GetYearlyWithdraws(ctx context.Context, year int) ([]*db.GetYearlyWithdrawsRow, error) {
	res, err := r.db.GetYearlyWithdraws(ctx, year)

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsFailed.WithInternal(err)
	}

	return res, nil

}

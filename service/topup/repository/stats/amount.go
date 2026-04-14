package topupstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsAmountRepository struct {
	db *db.Queries
}

func NewTopupStatsAmountRepository(db *db.Queries) TopupStatsAmountRepository {
	return &topupStatsAmountRepository{
		db: db,
	}
}

func (r *topupStatsAmountRepository) GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmounts(ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsFailed.WithInternal(err)
	}

	return res, nil
}

func (r *topupStatsAmountRepository) GetYearlyTopupAmounts(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountsRow, error) {
	res, err := r.db.GetYearlyTopupAmounts(ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsFailed.WithInternal(err)
	}

	return res, nil
}

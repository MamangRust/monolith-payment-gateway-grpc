package transferstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsAmountRepository struct {
	db *db.Queries
}

func NewTransferStatsAmountRepository(db *db.Queries) TransferStatsAmountRepository {
	return &transferStatsAmountRepository{
		db: db,
	}
}

func (r *transferStatsAmountRepository) GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(ctx, yearStart)

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transferStatsAmountRepository) GetYearlyTransferAmounts(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountsRow, error) {
	res, err := r.db.GetYearlyTransferAmounts(ctx, year)

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsFailed.WithInternal(err)
	}
	return res, nil
}

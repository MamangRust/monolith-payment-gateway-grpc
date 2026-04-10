package merchantstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsAmountRepository struct {
	db *db.Queries
}

func NewMerchantStatsAmountRepository(db *db.Queries) MerchantStatsAmountRepository {
	return &merchantStatsAmountRepository{
		db: db,
	}
}

func (r *merchantStatsAmountRepository) GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyAmountMerchantRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmountMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountMerchantFailed
	}

	return res, nil
}

func (r *merchantStatsAmountRepository) GetYearlyAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyAmountMerchantRow, error) {
	res, err := r.db.GetYearlyAmountMerchant(ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountMerchantFailed
	}

	return res, nil
}

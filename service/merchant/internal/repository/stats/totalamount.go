package merchantstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsTotalAmountRepository struct {
	db *db.Queries
}

func NewMerchantStatsTotalAmountRepository(db *db.Queries) MerchantStatsTotalAmountRepository {
	return &merchantStatsTotalAmountRepository{
		db: db,
	}
}

func (r *merchantStatsTotalAmountRepository) GetMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyTotalAmountMerchantRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountMerchantFailed
	}

	return res, nil
}

func (r *merchantStatsTotalAmountRepository) GetYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyTotalAmountMerchantRow, error) {
	res, err := r.db.GetYearlyTotalAmountMerchant(ctx, int32(year))

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountMerchantFailed
	}

	return res, nil
}

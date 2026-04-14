package merchantstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsMethodRepository struct {
	db *db.Queries
}

func NewMerchantStatsMethodRepository(db *db.Queries) MerchantStatsMethodRepository {
	return &merchantStatsMethodRepository{
		db: db,
	}
}

func (r *merchantStatsMethodRepository) GetMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsMerchantRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodsMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodsMerchantFailed.WithInternal(err)
	}

	return res, nil
}

func (r *merchantStatsMethodRepository) GetYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodMerchantRow, error) {
	res, err := r.db.GetYearlyPaymentMethodMerchant(ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodMerchantFailed.WithInternal(err)
	}

	return res, nil

}

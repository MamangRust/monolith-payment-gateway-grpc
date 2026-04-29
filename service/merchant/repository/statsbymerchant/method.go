package merchantstatsmerchantrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsMethodByMerchantRepository struct {
	db *db.Queries
}

func NewMerchantStatsMethodByMerchantRepository(db *db.Queries) MerchantStatsMethodByMerchantRepository {
	return &merchantStatsMethodByMerchantRepository{
		db: db,
	}
}

func (r *merchantStatsMethodByMerchantRepository) GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetMonthlyPaymentMethodByMerchantsRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodByMerchants(ctx, db.GetMonthlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByMerchantsFailed.WithInternal(err)
	}

	return res, nil
}

func (r *merchantStatsMethodByMerchantRepository) GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetYearlyPaymentMethodByMerchantsRow, error) {
	res, err := r.db.GetYearlyPaymentMethodByMerchants(ctx, db.GetYearlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByMerchantsFailed.WithInternal(err)
	}

	return res, nil
}

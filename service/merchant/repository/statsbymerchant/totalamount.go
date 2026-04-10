package merchantstatsmerchantrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsTotalAmountByMerchantRepository struct {
	db *db.Queries
}

func NewMerchantStatsTotalAmountByMerchantRepository(db *db.Queries) MerchantStatsTotalAmountByMerchantRepository {
	return &merchantStatsTotalAmountByMerchantRepository{
		db: db,
	}
}

func (r *merchantStatsTotalAmountByMerchantRepository) GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetMonthlyTotalAmountByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByMerchant(ctx, db.GetMonthlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByMerchantsFailed
	}

	return res, nil
}

func (r *merchantStatsTotalAmountByMerchantRepository) GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetYearlyTotalAmountByMerchantRow, error) {
	res, err := r.db.GetYearlyTotalAmountByMerchant(ctx, db.GetYearlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByMerchantsFailed
	}

	return res, nil
}

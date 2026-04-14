package merchantstatsapikeyrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsTotalAmountByApiKeyRepository struct {
	db *db.Queries
}

func NewMerchantStatsTotalAmountByApiKeyRepository(db *db.Queries) MerchantStatsTotalAmountByApiKeyRepository {
	return &merchantStatsTotalAmountByApiKeyRepository{
		db: db,
	}
}

func (r *merchantStatsTotalAmountByApiKeyRepository) GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetMonthlyTotalAmountByApikeyRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByApikey(ctx, db.GetMonthlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByApikeyFailed.WithInternal(err)
	}

	return res, nil
}

func (r *merchantStatsTotalAmountByApiKeyRepository) GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetYearlyTotalAmountByApikeyRow, error) {
	res, err := r.db.GetYearlyTotalAmountByApikey(ctx, db.GetYearlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByApikeyFailed.WithInternal(err)
	}

	return res, nil
}

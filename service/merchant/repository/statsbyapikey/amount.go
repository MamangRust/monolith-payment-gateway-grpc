package merchantstatsapikeyrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsAmountByApiKeyRepository struct {
	db *db.Queries
}

func NewMerchantStatsAmountByApiKeyRepository(db *db.Queries) MerchantStatsAmountByApiKeyRepository {
	return &merchantStatsAmountByApiKeyRepository{
		db: db,
	}
}

func (r *merchantStatsAmountByApiKeyRepository) GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetMonthlyAmountByApikeyRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByApikey(ctx, db.GetMonthlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByApikeyFailed
	}

	return res, nil
}

func (r *merchantStatsAmountByApiKeyRepository) GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetYearlyAmountByApikeyRow, error) {
	res, err := r.db.GetYearlyAmountByApikey(ctx, db.GetYearlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByApikeyFailed
	}

	return res, nil
}

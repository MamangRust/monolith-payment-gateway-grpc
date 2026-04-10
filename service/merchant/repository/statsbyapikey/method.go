package merchantstatsapikeyrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantStatsMethodByApiKeyRepository struct {
	db *db.Queries
}

func NewMerchantStatsMethodByApiKeyRepository(db *db.Queries) MerchantStatsMethodByApiKeyRepository {
	return &merchantStatsMethodByApiKeyRepository{
		db: db,
	}
}

func (r *merchantStatsMethodByApiKeyRepository) GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetMonthlyPaymentMethodByApikeyRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByApikey(ctx, db.GetMonthlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByApikeyFailed
	}

	return res, nil
}

func (r *merchantStatsMethodByApiKeyRepository) GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetYearlyPaymentMethodByApikeyRow, error) {
	res, err := r.db.GetYearlyPaymentMethodByApikey(ctx, db.GetYearlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByApikeyFailed
	}

	return res, nil
}

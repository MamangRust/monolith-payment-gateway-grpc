package merchantstatsapikeyrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/statsByApiKey"
)

type merchantStatsMethodByApiKeyRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticMethodByApiKeyMapper
}

func NewMerchantStatsMethodByApiKeyRepository(db *db.Queries, mapper recordmapper.MerchantStatisticMethodByApiKeyMapper) MerchantStatsMethodByApiKeyRepository {
	return &merchantStatsMethodByApiKeyRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyPaymentMethodByApikey retrieves monthly payment method statistics for a specific merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing month, year, and API key.
//
// Returns:
//   - []*record.MerchantMonthlyPaymentMethod: The list of monthly payment method records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsMethodByApiKeyRepository) GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByApikey(ctx, db.GetMonthlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByApikeyFailed
	}

	return r.mapper.ToMerchantMonthlyPaymentMethodsByApikey(res), nil
}

// GetYearlyPaymentMethodByApikey retrieves yearly payment method statistics for a specific merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing year and API key.
//
// Returns:
//   - []*record.MerchantYearlyPaymentMethod: The list of yearly payment method records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsMethodByApiKeyRepository) GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByApikey(ctx, db.GetYearlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByApikeyFailed
	}

	return r.mapper.ToMerchantYearlyPaymentMethodsByApikey(res), nil
}

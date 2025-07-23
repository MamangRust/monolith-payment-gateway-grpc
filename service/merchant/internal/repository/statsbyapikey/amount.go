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

type merchantStatsAmountByApiKeyRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticAmountByApiKeyMapper
}

func NewMerchantStatsAmountByApiKeyRepository(db *db.Queries, mapper recordmapper.MerchantStatisticAmountByApiKeyMapper) MerchantStatsAmountByApiKeyRepository {
	return &merchantStatsAmountByApiKeyRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyAmountByApikey retrieves monthly transaction amount statistics for a specific merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing month, year, and API key.
//
// Returns:
//   - []*record.MerchantMonthlyAmount: The list of monthly amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsAmountByApiKeyRepository) GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByApikey(ctx, db.GetMonthlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByApikeyFailed
	}

	return r.mapper.ToMerchantMonthlyAmountsByApikey(res), nil
}

// GetYearlyAmountByApikey retrieves yearly transaction amount statistics for a specific merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing year and API key.
//
// Returns:
//   - []*record.MerchantYearlyAmount: The list of yearly amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsAmountByApiKeyRepository) GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByApikey(ctx, db.GetYearlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByApikeyFailed
	}

	return r.mapper.ToMerchantYearlyAmountsByApikey(res), nil
}

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

type merchantStatsTotalAmountByApiKeyRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticTotalAmountByApiKeyMapper
}

func NewMerchantStatsTotalAmountByApiKeyRepository(db *db.Queries, mapper recordmapper.MerchantStatisticTotalAmountByApiKeyMapper) MerchantStatsTotalAmountByApiKeyRepository {
	return &merchantStatsTotalAmountByApiKeyRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTotalAmountByApikey retrieves monthly total transaction amount statistics for a specific merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing month, year, and API key.
//
// Returns:
//   - []*record.MerchantMonthlyTotalAmount: The list of monthly total amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsTotalAmountByApiKeyRepository) GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByApikey(ctx, db.GetMonthlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByApikeyFailed
	}

	return r.mapper.ToMerchantMonthlyTotalAmountsByApikey(res), nil
}

// GetYearlyTotalAmountByApikey retrieves yearly total transaction amount statistics for a specific merchant using API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing year and API key.
//
// Returns:
//   - []*record.MerchantYearlyTotalAmount: The list of yearly total amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsTotalAmountByApiKeyRepository) GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByApikey(ctx, db.GetYearlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByApikeyFailed
	}

	return r.mapper.ToMerchantYearlyTotalAmountsByApikey(res), nil
}

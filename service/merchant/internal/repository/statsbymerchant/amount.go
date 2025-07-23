package merchantstatsmerchantrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/statsByMerchant"
)

type merchantStatsAmountByMerchantRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticAmountByMerchantMapper
}

func NewMerchantStatsAmountByMerchantRepository(db *db.Queries, mapper recordmapper.MerchantStatisticAmountByMerchantMapper) MerchantStatsAmountByMerchantRepository {
	return &merchantStatsAmountByMerchantRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyAmountByMerchants retrieves monthly transaction amount statistics for a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing month, year, and merchant ID.
//
// Returns:
//   - []*record.MerchantMonthlyAmount: The list of monthly amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsAmountByMerchantRepository) GetMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByMerchants(ctx, db.GetMonthlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByMerchantsFailed
	}

	return r.mapper.ToMerchantMonthlyAmountsByMerchant(res), nil
}

// GetYearlyAmountByMerchants retrieves yearly transaction amount statistics for a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing year and merchant ID.
//
// Returns:
//   - []*record.MerchantYearlyAmount: The list of yearly amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsAmountByMerchantRepository) GetYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByMerchants(ctx, db.GetYearlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByMerchantsFailed
	}

	return r.mapper.ToMerchantYearlyAmountsByMerchant(res), nil
}

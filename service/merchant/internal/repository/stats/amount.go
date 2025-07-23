package merchantstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/stats"
)

type merchantStatsAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticAmountRecordMapper
}

func NewMerchantStatsAmountRepository(db *db.Queries, mapper recordmapper.MerchantStatisticAmountRecordMapper) MerchantStatsAmountRepository {
	return &merchantStatsAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyAmountMerchant retrieves monthly transaction amount statistics for merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.MerchantMonthlyAmount: The list of monthly amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsAmountRepository) GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmountMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountMerchantFailed
	}

	return r.mapper.ToMerchantMonthlyAmounts(res), nil
}

// GetYearlyAmountMerchant retrieves yearly transaction amount statistics for merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.MerchantYearlyAmount: The list of yearly amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsAmountRepository) GetYearlyAmountMerchant(ctx context.Context, year int) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountMerchant(ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountMerchantFailed
	}

	return r.mapper.ToMerchantYearlyAmounts(res), nil
}

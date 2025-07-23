package merchantstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/stats"
)

type merchantStatsTotalAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticTotalAmountRecordMapper
}

func NewMerchantStatsTotalAmountRepository(db *db.Queries, mapper recordmapper.MerchantStatisticTotalAmountRecordMapper) MerchantStatsTotalAmountRepository {
	return &merchantStatsTotalAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTotalAmountMerchant retrieves monthly total transaction amount statistics across all merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.MerchantMonthlyTotalAmount: The list of monthly total amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsTotalAmountRepository) GetMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountMerchantFailed
	}

	return r.mapper.ToMerchantMonthlyTotalAmounts(res), nil
}

// GetYearlyTotalAmountMerchant retrieves yearly total transaction amount statistics across all merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.MerchantYearlyTotalAmount: The list of yearly total amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsTotalAmountRepository) GetYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountMerchant(ctx, int32(year))

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountMerchantFailed
	}

	return r.mapper.ToMerchantYearlyTotalAmounts(res), nil
}

package merchantstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant/stats"
)

type merchantStatsMethodRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticMethodRecordMapper
}

func NewMerchantStatsMethodRepository(db *db.Queries, mapper recordmapper.MerchantStatisticMethodRecordMapper) MerchantStatsMethodRepository {
	return &merchantStatsMethodRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyPaymentMethodsMerchant retrieves monthly merchant payment method statistics for a given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.MerchantMonthlyPaymentMethod: The list of monthly payment method records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsMethodRepository) GetMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodsMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodsMerchantFailed
	}

	return r.mapper.ToMerchantMonthlyPaymentMethods(res), nil
}

// GetYearlyPaymentMethodMerchant retrieves yearly merchant payment method statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.MerchantYearlyPaymentMethod: The list of yearly payment method records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsMethodRepository) GetYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodMerchant(ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodMerchantFailed
	}

	return r.mapper.ToMerchantYearlyPaymentMethods(res), nil

}

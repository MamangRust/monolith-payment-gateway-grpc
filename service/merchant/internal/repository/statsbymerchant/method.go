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

type merchantStatsMethodByMerchantRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticMethodByMerchantMapper
}

func NewMerchantStatsMethodByMerchantRepository(db *db.Queries, mapper recordmapper.MerchantStatisticMethodByMerchantMapper) MerchantStatsMethodByMerchantRepository {
	return &merchantStatsMethodByMerchantRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyPaymentMethodByMerchants retrieves monthly payment method statistics for a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing month, year, and merchant ID.
//
// Returns:
//   - []*record.MerchantMonthlyPaymentMethod: The list of monthly payment method records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsMethodByMerchantRepository) GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodByMerchants(ctx, db.GetMonthlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByMerchantsFailed
	}

	return r.mapper.ToMerchantMonthlyPaymentMethodsByMerchant(res), nil
}

// GetYearlyPaymentMethodByMerchants retrieves yearly payment method statistics for a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing year and merchant ID.
//
// Returns:
//   - []*record.MerchantYearlyPaymentMethod: The list of yearly payment method records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsMethodByMerchantRepository) GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByMerchants(ctx, db.GetYearlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByMerchantsFailed
	}

	return r.mapper.ToMerchantYearlyPaymentMethodsByMerchant(res), nil
}

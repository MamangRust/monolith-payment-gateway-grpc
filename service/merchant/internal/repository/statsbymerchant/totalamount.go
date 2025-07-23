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

type merchantStatsTotalAmountByMerchantRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantStatisticTotalAmountByMerchantMapper
}

func NewMerchantStatsTotalAmountByMerchantRepository(db *db.Queries, mapper recordmapper.MerchantStatisticTotalAmountByMerchantMapper) MerchantStatsTotalAmountByMerchantRepository {
	return &merchantStatsTotalAmountByMerchantRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTotalAmountByMerchants retrieves monthly total amount statistics for a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing month, year, and merchant ID.
//
// Returns:
//   - []*record.MerchantMonthlyTotalAmount: The list of monthly total amount records.
//   - error: An error if any occurred during the query.

func (r *merchantStatsTotalAmountByMerchantRepository) GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByMerchant(ctx, db.GetMonthlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByMerchantsFailed
	}

	return r.mapper.ToMerchantMonthlyTotalAmountsByMerchant(res), nil
}

// GetYearlyTotalAmountByMerchants retrieves yearly total amount statistics for a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing year and merchant ID.
//
// Returns:
//   - []*record.MerchantYearlyTotalAmount: The list of yearly total amount records.
//   - error: An error if any occurred during the query.
func (r *merchantStatsTotalAmountByMerchantRepository) GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByMerchant(ctx, db.GetYearlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByMerchantsFailed
	}

	return r.mapper.ToMerchantYearlyTotalAmountsByMerchant(res), nil
}

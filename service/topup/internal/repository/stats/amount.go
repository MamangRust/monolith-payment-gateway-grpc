package topupstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/stats"
)

type topupStatsAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupStatisticAmountMapper
}

func NewTopupStatsAmountRepository(db *db.Queries, mapper recordmapper.TopupStatisticAmountMapper) TopupStatsAmountRepository {
	return &topupStatsAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTopupAmounts retrieves monthly statistics of topup amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.TopupMonthAmount: List of monthly topup amount statistics.
//   - error: Error if the query fails.
func (r *topupStatsAmountRepository) GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*record.TopupMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmounts(ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsFailed
	}

	return r.mapper.ToTopupMonthlyAmounts(res), nil
}

// GetYearlyTopupAmounts retrieves yearly statistics of topup amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the data is requested.
//
// Returns:
//   - []*record.TopupYearlyAmount: List of yearly topup amount statistics.
//   - error: Error if the query fails.
func (r *topupStatsAmountRepository) GetYearlyTopupAmounts(ctx context.Context, year int) ([]*record.TopupYearlyAmount, error) {
	res, err := r.db.GetYearlyTopupAmounts(ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsFailed
	}

	return r.mapper.ToTopupYearlyAmounts(res), nil
}

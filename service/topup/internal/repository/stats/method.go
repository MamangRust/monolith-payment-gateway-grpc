package topupstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/stats"
)

type topupStatsMethodRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupStatisticMethodMapper
}

func NewTopupStatsMethodRepository(db *db.Queries, mapper recordmapper.TopupStatisticMethodMapper) TOpupStatsMethodRepository {
	return &topupStatsMethodRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTopupMethods retrieves monthly statistics of topup methods used.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly method statistics are requested.
//
// Returns:
//   - []*record.TopupMonthMethod: List of monthly topup method usage.
//   - error: Error if the query fails.
func (r *topupStatsMethodRepository) GetMonthlyTopupMethods(ctx context.Context, year int) ([]*record.TopupMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethods(ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsFailed
	}

	return r.mapper.ToTopupMonthlyMethods(res), nil
}

// GetYearlyTopupMethods retrieves yearly statistics of topup methods used.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly method statistics are requested.
//
// Returns:
//   - []*record.TopupYearlyMethod: List of yearly topup method usage.
//   - error: Error if the query fails.
func (r *topupStatsMethodRepository) GetYearlyTopupMethods(ctx context.Context, year int) ([]*record.TopupYearlyMethod, error) {
	res, err := r.db.GetYearlyTopupMethods(ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsFailed
	}

	return r.mapper.ToTopupYearlyMethods(res), nil
}

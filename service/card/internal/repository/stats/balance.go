package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/stats"
)

type cardStatsBalanceRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticBalanceRecordMapper
}

func NewCardStatsBalanceRepository(db *db.Queries, mapper recordmapper.CardStatisticBalanceRecordMapper) CardStatsBalanceRepository {
	return &cardStatsBalanceRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyBalance retrieves the monthly balance of all cards for a given year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the balance is requested.
//
// Returns:
//   - A slice of CardMonthBalance containing the balance for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyBalanceFailed.
func (r *cardStatsBalanceRepository) GetMonthlyBalance(ctx context.Context, year int) ([]*record.CardMonthBalance, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalances(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceFailed
	}

	return r.mapper.ToMonthlyBalances(res), nil
}

// GetYearlyBalance retrieves the yearly balance of all cards for a given year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the balance is requested.
//
// Returns:
//   - A slice of CardYearlyBalance containing the balance for each year.
//   - An error if the retrieval fails, of type ErrGetYearlyBalanceFailed.
func (r *cardStatsBalanceRepository) GetYearlyBalance(ctx context.Context, year int) ([]*record.CardYearlyBalance, error) {
	res, err := r.db.GetYearlyBalances(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceFailed
	}

	return r.mapper.ToYearlyBalances(res), nil
}

package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/stats"
)

type cardStatsTopupRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticTopupRecordMapper
}

func NewCardStatsTopupRepository(db *db.Queries, mapper recordmapper.CardStatisticTopupRecordMapper) CardStatsTopupRepository {
	return &cardStatsTopupRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTopupAmount retrieves the monthly top-up amount for all cards for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the top-up amount is requested.
//
// Returns:
//   - A slice of CardMonthAmount containing the top-up amount for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTopupAmountFailed.
func (r *cardStatsTopupRepository) GetMonthlyTopupAmount(ctx context.Context, year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountFailed
	}

	return r.mapper.ToMonthlyTopupAmounts(res), nil
}

// GetYearlyTopupAmount retrieves the yearly top-up amount for all cards for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the top-up amount is requested.
//
// Returns:
//   - A slice of CardYearAmount containing the top-up amount for each year.
//   - An error if the retrieval fails, of type ErrGetYearlyTopupAmountFailed.
func (r *cardStatsTopupRepository) GetYearlyTopupAmount(ctx context.Context, year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTopupAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountFailed
	}

	return r.mapper.ToYearlyTopupAmounts(res), nil
}

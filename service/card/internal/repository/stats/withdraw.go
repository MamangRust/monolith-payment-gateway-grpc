package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/stats"
)

type cardStatsWithdrawRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticWithdrawRecordMapper
}

func NewCardStatsWithdrawRepository(db *db.Queries, mapper recordmapper.CardStatisticWithdrawRecordMapper) CardStatsWithdrawRepository {
	return &cardStatsWithdrawRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyWithdrawAmount retrieves the monthly withdrawal amount for all cards for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the withdrawal amount is requested.
//
// Returns:
//   - A slice of CardMonthAmount containing the withdrawal amount for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyWithdrawAmountFailed.
func (r *cardStatsWithdrawRepository) GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountFailed
	}

	return r.mapper.ToMonthlyWithdrawAmounts(res), nil
}

// GetYearlyWithdrawAmount retrieves the yearly withdrawal amount for all cards for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the withdrawal amount is requested.
//
// Returns:
//   - A slice of CardYearAmount containing the withdrawal amount for each year.
//   - An error if the retrieval fails, of type ErrGetYearlyWithdrawAmountFailed.
func (r *cardStatsWithdrawRepository) GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyWithdrawAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountFailed
	}

	return r.mapper.ToYearlyWithdrawAmounts(res), nil
}

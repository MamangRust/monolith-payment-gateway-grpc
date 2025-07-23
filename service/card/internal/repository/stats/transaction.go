package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/stats"
)

type cardStatsTransactionRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticTransactionRecordMapper
}

func NewCardStatsTransactionRepository(db *db.Queries, mapper recordmapper.CardStatisticTransactionRecordMapper) CardStatsTransactionRepository {
	return &cardStatsTransactionRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransactionAmount retrieves the monthly transaction amounts for all cards for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the transaction amounts are requested.
//
// Returns:
//   - A slice of CardMonthAmount containing the transaction amounts for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTransactionAmountFailed.
func (r *cardStatsTransactionRepository) GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountFailed
	}

	return r.mapper.ToMonthlyTransactionAmounts(res), nil
}

// GetYearlyTransactionAmount retrieves the yearly transaction amounts for all cards for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the transaction amounts are requested.
//
// Returns:
//   - A slice of CardYearAmount containing the transaction amounts for the specified year.
//   - An error if the retrieval fails, of type ErrGetYearlyTransactionAmountFailed.
func (r *cardStatsTransactionRepository) GetYearlyTransactionAmount(ctx context.Context, year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransactionAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountFailed
	}

	return r.mapper.ToYearlyTransactionAmounts(res), nil
}

package repositorystatsbycard

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/statsbycard"
)

type cardStatsBalanceByCardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticBalanceByCardRecordMapper
}

func NewCardStatsBalanceByCardRepository(db *db.Queries, mapper recordmapper.CardStatisticBalanceByCardRecordMapper) CardStatsBalanceByCardRepository {
	return &cardStatsBalanceByCardRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyBalancesByCardNumber retrieves the monthly balance data for a given card number
// and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardMonthBalance containing the balance data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyBalanceByCardFailed.
func (r *cardStatsBalanceByCardRepository) GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthBalance, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalancesByCardNumber(ctx, db.GetMonthlyBalancesByCardNumberParams{
		Column1:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceByCardFailed
	}

	return r.mapper.ToMonthlyBalancesCardNumber(res), nil
}

// GetYearlyBalanceByCardNumber retrieves the yearly balance data for a given card number
// and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardYearlyBalance containing the balance data for the given year.
//   - An error if the retrieval fails, of type ErrGetYearlyBalanceByCardFailed.
func (r *cardStatsBalanceByCardRepository) GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearlyBalance, error) {
	res, err := r.db.GetYearlyBalancesByCardNumber(ctx, db.GetYearlyBalancesByCardNumberParams{
		Column1:    req.Year,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceByCardFailed
	}

	return r.mapper.ToYearlyBalancesCardNumber(res), nil
}

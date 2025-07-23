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

type cardStatsTopupByCardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticTopupByCardRecordMapper
}

func NewCardStatsTopupByCardRepository(db *db.Queries, mapper recordmapper.CardStatisticTopupByCardRecordMapper) CardStatsTopupByCardRepository {
	return &cardStatsTopupByCardRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTopupAmountByCardNumber retrieves the monthly topup amount data for a given card number
// and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the topup amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTopupAmountByCardFailed.
func (r *cardStatsTopupByCardRepository) GetMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountByCardNumber(ctx, db.GetMonthlyTopupAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountByCardFailed
	}

	return r.mapper.ToMonthlyTopupAmountsByCardNumber(res), nil
}

// GetYearlyTopupAmountByCardNumber retrieves the yearly topup amount data for a given card number
// and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the topup amount data for the given year.
//   - An error if the retrieval fails, of type ErrGetYearlyTopupAmountByCardFailed.
func (r *cardStatsTopupByCardRepository) GetYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTopupAmountByCardNumber(ctx, db.GetYearlyTopupAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountByCardFailed
	}

	return r.mapper.ToYearlyTopupAmountsByCardNumber(res), nil
}

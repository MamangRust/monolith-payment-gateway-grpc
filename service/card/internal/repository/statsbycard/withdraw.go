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

type cardStatsWithdrawByCardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticWithdrawByCardRecordMapper
}

func NewCardStatsWithdrawByCardRepository(db *db.Queries, mapper recordmapper.CardStatisticWithdrawByCardRecordMapper) CardStatsWithdrawByCardRepository {
	return &cardStatsWithdrawByCardRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyWithdrawAmountByCardNumber retrieves the monthly withdrawal amount data
// for a given card number and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the withdrawal amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyWithdrawAmountByCardFailed.
func (r *cardStatsWithdrawByCardRepository) GetMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmountByCardNumber(ctx, db.GetMonthlyWithdrawAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountByCardFailed
	}

	return r.mapper.ToMonthlyWithdrawAmountsByCardNumber(res), nil
}

// GetYearlyWithdrawAmountByCardNumber retrieves the yearly withdrawal amount data
// for a given card number and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the withdrawal amount data for the given year.
//   - An error if the retrieval fails, of type ErrGetYearlyWithdrawAmountByCardFailed.
func (r *cardStatsWithdrawByCardRepository) GetYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyWithdrawAmountByCardNumber(ctx, db.GetYearlyWithdrawAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountByCardFailed
	}

	return r.mapper.ToYearlyWithdrawAmountsByCardNumber(res), nil
}

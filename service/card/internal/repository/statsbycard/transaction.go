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

type cardStatsTransactionByCardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticTransactionByCardRecordMapper
}

func NewCardStatsTransactionByCardRepository(db *db.Queries, mapper recordmapper.CardStatisticTransactionByCardRecordMapper) CardStatsTransactionByCardRepository {
	return &cardStatsTransactionByCardRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransactionAmountByCardNumber retrieves the monthly transaction amount data
// for a given card number and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the transaction amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTransactionAmountByCardFailed.
func (r *cardStatsTransactionByCardRepository) GetMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmountByCardNumber(ctx, db.GetMonthlyTransactionAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountByCardFailed
	}

	return r.mapper.ToMonthlyTransactionAmountsByCardNumber(res), nil
}

// GetYearlyTransactionAmountByCardNumber retrieves the yearly transaction amount data
// for a given card number and year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the transaction amount data for the given year.
//   - An error if the retrieval fails, of type ErrGetYearlyTransactionAmountByCardFailed.
func (r *cardStatsTransactionByCardRepository) GetYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransactionAmountByCardNumber(ctx, db.GetYearlyTransactionAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountByCardFailed
	}

	return r.mapper.ToYearlyTransactionAmountsByCardNumber(res), nil
}

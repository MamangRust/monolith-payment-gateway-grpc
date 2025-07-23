package topupstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/statsbycard"
)

type topupStatsByCardAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupStatisticAmountByCardNumberMapper
}

func NewTopupStatsByCardAmountRepository(db *db.Queries, mapper recordmapper.TopupStatisticAmountByCardNumberMapper) TopupStatsByCardAmountRepository {
	return &topupStatsByCardAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTopupAmountsByCardNumber retrieves monthly topup amount statistics for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TopupMonthAmount: List of monthly topup amount data.
//   - error: Error if the query fails.
func (r *topupStatsByCardAmountRepository) GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupMonthAmount, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountsByCardNumber(ctx, db.GetMonthlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsByCardFailed
	}

	return r.mapper.ToTopupMonthlyAmountsByCardNumber(res), nil
}

// GetYearlyTopupAmountsByCardNumber retrieves yearly topup amount statistics for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TopupYearlyAmount: List of yearly topup amount data.
//   - error: Error if the query fails.
func (r *topupStatsByCardAmountRepository) GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupYearlyAmount, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupAmountsByCardNumber(ctx, db.GetYearlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsByCardFailed
	}

	return r.mapper.ToTopupYearlyAmountsByCardNumber(res), nil
}

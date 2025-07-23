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

type topupStatsByCardMethodRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupStatisticMethodByCardNumberMapper
}

func NewTopupStatsByCardMethodRepository(db *db.Queries, mapper recordmapper.TopupStatisticMethodByCardNumberMapper) TopupStatsByCardMethodRepository {
	return &topupStatsByCardMethodRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTopupMethodsByCardNumber retrieves monthly topup method statistics for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TopupMonthMethod: List of monthly topup method usage.
//   - error: Error if the query fails.
func (r *topupStatsByCardMethodRepository) GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupMonthMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethodsByCardNumber(ctx, db.GetMonthlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsByCardFailed
	}

	return r.mapper.ToTopupMonthlyMethodsByCardNumber(res), nil
}

// GetYearlyTopupMethodsByCardNumber retrieves yearly topup method statistics for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TopupYearlyMethod: List of yearly topup method usage.
//   - error: Error if the query fails.
func (r *topupStatsByCardMethodRepository) GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupYearlyMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupMethodsByCardNumber(ctx, db.GetYearlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsByCardFailed
	}

	return r.mapper.ToTopupYearlyMethodsByCardNumber(res), nil
}

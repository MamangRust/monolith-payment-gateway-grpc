package topupstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsByCardAmountRepository struct {
	db *db.Queries
}

func NewTopupStatsByCardAmountRepository(db *db.Queries) TopupStatsByCardAmountRepository {
	return &topupStatsByCardAmountRepository{
		db: db,
	}
}

func (r *topupStatsByCardAmountRepository) GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupAmountsByCardNumberRow, error) {
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

	return res, nil
}

func (r *topupStatsByCardAmountRepository) GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupAmountsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupAmountsByCardNumber(ctx, db.GetYearlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsByCardFailed
	}

	return res, nil
}

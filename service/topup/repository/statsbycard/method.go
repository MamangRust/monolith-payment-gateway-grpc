package topupstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsByCardMethodRepository struct {
	db *db.Queries
}

func NewTopupStatsByCardMethodRepository(db *db.Queries) TopupStatsByCardMethodRepository {
	return &topupStatsByCardMethodRepository{
		db: db,
	}
}

func (r *topupStatsByCardMethodRepository) GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethodsByCardNumber(ctx, db.GetMonthlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *topupStatsByCardMethodRepository) GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupMethodsByCardNumber(ctx, db.GetYearlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsByCardFailed.WithInternal(err)
	}

	return res, nil
}

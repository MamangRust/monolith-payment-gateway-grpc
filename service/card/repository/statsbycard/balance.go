package repositorystatsbycard

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStatsBalanceByCardRepository struct {
	db *db.Queries
}

func NewCardStatsBalanceByCardRepository(db *db.Queries) CardStatsBalanceByCardRepository {
	return &cardStatsBalanceByCardRepository{
		db: db,
	}
}

func (r *cardStatsBalanceByCardRepository) GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyBalancesByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalancesByCardNumber(ctx, db.GetMonthlyBalancesByCardNumberParams{
		Column1:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceByCardFailed
	}

	return res, nil
}

func (r *cardStatsBalanceByCardRepository) GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyBalancesByCardNumberRow, error) {
	res, err := r.db.GetYearlyBalancesByCardNumber(ctx, db.GetYearlyBalancesByCardNumberParams{
		Column1:    req.Year,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceByCardFailed
	}

	return res, nil
}

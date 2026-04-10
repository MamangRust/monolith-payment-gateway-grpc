package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStatsTransactionRepository struct {
	db *db.Queries
}

func NewCardStatsTransactionRepository(db *db.Queries) CardStatsTransactionRepository {
	return &cardStatsTransactionRepository{
		db: db,
	}
}

func (r *cardStatsTransactionRepository) GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*db.GetMonthlyTransactionAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountFailed
	}

	return res, nil
}

func (r *cardStatsTransactionRepository) GetYearlyTransactionAmount(ctx context.Context, year int) ([]*db.GetYearlyTransactionAmountRow, error) {
	res, err := r.db.GetYearlyTransactionAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountFailed
	}

	return res, nil
}

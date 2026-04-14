package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStatsBalanceRepository struct {
	db *db.Queries
}

func NewCardStatsBalanceRepository(db *db.Queries) CardStatsBalanceRepository {
	return &cardStatsBalanceRepository{
		db: db,
	}
}

func (r *cardStatsBalanceRepository) GetMonthlyBalance(ctx context.Context, year int) ([]*db.GetMonthlyBalancesRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalances(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceFailed.WithInternal(err)
	}

	return res, nil
}

func (r *cardStatsBalanceRepository) GetYearlyBalance(ctx context.Context, year int) ([]*db.GetYearlyBalancesRow, error) {
	res, err := r.db.GetYearlyBalances(ctx, year)

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceFailed.WithInternal(err)
	}

	return res, nil
}

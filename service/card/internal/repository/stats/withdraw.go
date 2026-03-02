package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStatsWithdrawRepository struct {
	db *db.Queries
}

func NewCardStatsWithdrawRepository(db *db.Queries) CardStatsWithdrawRepository {
	return &cardStatsWithdrawRepository{
		db: db,
	}
}

func (r *cardStatsWithdrawRepository) GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountFailed
	}

	return res, nil
}

func (r *cardStatsWithdrawRepository) GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetYearlyWithdrawAmountRow, error) {
	res, err := r.db.GetYearlyWithdrawAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountFailed
	}

	return res, nil
}

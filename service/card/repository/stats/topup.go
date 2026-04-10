package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStatsTopupRepository struct {
	db *db.Queries
}

func NewCardStatsTopupRepository(db *db.Queries) CardStatsTopupRepository {
	return &cardStatsTopupRepository{
		db: db,
	}
}

func (r *cardStatsTopupRepository) GetMonthlyTopupAmount(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountFailed
	}

	return res, nil
}

func (r *cardStatsTopupRepository) GetYearlyTopupAmount(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountRow, error) {
	res, err := r.db.GetYearlyTopupAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountFailed
	}

	return res, nil
}

package topupstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsMethodRepository struct {
	db *db.Queries
}

func NewTopupStatsMethodRepository(db *db.Queries) TopupStatsMethodRepository {
	return &topupStatsMethodRepository{
		db: db,
	}
}

func (r *topupStatsMethodRepository) GetMonthlyTopupMethods(ctx context.Context, year int) ([]*db.GetMonthlyTopupMethodsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethods(ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsFailed.WithInternal(err)
	}

	return res, nil
}

func (r *topupStatsMethodRepository) GetYearlyTopupMethods(ctx context.Context, year int) ([]*db.GetYearlyTopupMethodsRow, error) {
	res, err := r.db.GetYearlyTopupMethods(ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsFailed.WithInternal(err)
	}

	return res, nil
}

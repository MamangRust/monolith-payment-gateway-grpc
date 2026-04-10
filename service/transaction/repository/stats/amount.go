package transactionstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsAmountRepository struct {
	db *db.Queries
}

func NewTransactionStatsAmountRepository(db *db.Queries) TransactionStatsAmountRepository {
	return &transactionStatsAmountRepository{
		db: db,
	}
}

func (r *transactionStatsAmountRepository) GetMonthlyAmounts(ctx context.Context, year int) ([]*db.GetMonthlyAmountsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsFailed
	}

	return res, nil
}

func (r *transactionStatsAmountRepository) GetYearlyAmounts(ctx context.Context, year int) ([]*db.GetYearlyAmountsRow, error) {
	res, err := r.db.GetYearlyAmounts(ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsFailed
	}

	return res, nil
}

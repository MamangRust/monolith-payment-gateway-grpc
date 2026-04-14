package transactionstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsMethodRepository struct {
	db *db.Queries
}

func NewTransactionStatsMethodRepository(db *db.Queries) TransactionStatsMethodRepository {
	return &transactionStatsMethodRepository{
		db: db,
	}
}

func (r *transactionStatsMethodRepository) GetMonthlyPaymentMethods(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsMethodRepository) GetYearlyPaymentMethods(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodsRow, error) {
	res, err := r.db.GetYearlyPaymentMethods(ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsFailed.WithInternal(err)
	}

	return res, nil
}

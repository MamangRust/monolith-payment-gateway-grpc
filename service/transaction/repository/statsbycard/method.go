package transactionbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsByCardMethodRepository struct {
	db *db.Queries
}

func NewTransactionStatsByCardMethodRepository(db *db.Queries) TransactionStatsByCardMethodRepository {
	return &transactionStatsByCardMethodRepository{
		db: db,
	}
}

func (r *transactionStatsByCardMethodRepository) GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyPaymentMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetMonthlyPaymentMethodsByCardNumber(ctx, db.GetMonthlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsByCardFailed
	}

	return res, nil
}

func (r *transactionStatsByCardMethodRepository) GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyPaymentMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsByCardFailed
	}

	return res, nil
}

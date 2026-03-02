package transactionbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsByCardAmountRepository struct {
	db *db.Queries
}

func NewTransactionStatsByCardAmountRepository(db *db.Queries) TransactonStatsByCardAmountRepository {
	return &transactionStatsByCardAmountRepository{
		db: db,
	}
}

func (r *transactionStatsByCardAmountRepository) GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyAmountsByCardNumberRow, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetMonthlyAmountsByCardNumber(ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsByCardFailed
	}

	return res, nil
}

func (r *transactionStatsByCardAmountRepository) GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyAmountsByCardNumberRow, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyAmountsByCardNumber(ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsByCardFailed
	}

	return res, nil
}

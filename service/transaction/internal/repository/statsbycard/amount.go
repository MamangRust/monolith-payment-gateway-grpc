package transactionbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/statsbycard"
)

type transactionStatsByCardAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionStatisticByCardAmountMapper
}

func NewTransactionStatsByCardAmountRepository(db *db.Queries, mapper recordmapper.TransactionStatisticByCardAmountMapper) TransactonStatsByCardAmountRepository {
	return &transactionStatsByCardAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyAmountsByCardNumber retrieves monthly transaction amount statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TransactionMonthAmount: List of monthly transaction amounts.
//   - error: Error if any occurs.
func (r *transactionStatsByCardAmountRepository) GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetMonthlyAmountsByCardNumber(ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsByCardFailed
	}

	return r.mapper.ToTransactionMonthlyAmountsByCardNumber(res), nil
}

// GetYearlyAmountsByCardNumber retrieves yearly transaction amount statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TransactionYearlyAmount: List of yearly transaction amounts.
//   - error: Error if any occurs.
func (r *transactionStatsByCardAmountRepository) GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyAmountsByCardNumber(ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsByCardFailed
	}

	return r.mapper.ToTransactionYearlyAmountsByCardNumber(res), nil
}

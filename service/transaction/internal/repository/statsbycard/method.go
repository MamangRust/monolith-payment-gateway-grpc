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

type transactionStatsByCardMethodRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionStatisticByCardMethodMapper
}

func NewTransactionStatsByCardMethodRepository(db *db.Queries, mapper recordmapper.TransactionStatisticByCardMethodMapper) TransactionStatsByCardMethodRepository {
	return &transactionStatsByCardMethodRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyPaymentMethodsByCardNumber retrieves monthly transaction method statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TransactionMonthMethod: List of monthly payment method usage.
//   - error: Error if any occurs.
func (r *transactionStatsByCardMethodRepository) GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetMonthlyPaymentMethodsByCardNumber(ctx, db.GetMonthlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsByCardFailed
	}

	return r.mapper.ToTransactionMonthlyMethodsByCardNumber(res), nil
}

// GetYearlyPaymentMethodsByCardNumber retrieves yearly transaction method statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TransactionYearMethod: List of yearly payment method usage.
//   - error: Error if any occurs.
func (r *transactionStatsByCardMethodRepository) GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsByCardFailed
	}

	return r.mapper.ToTransactionYearlyMethodsByCardNumber(res), nil
}

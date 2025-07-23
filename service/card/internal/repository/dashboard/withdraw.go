package repositorydashboard

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardDashboardWithdrawRepository struct {
	db *db.Queries
}

func NewCardDashboardWithdrawRepository(db *db.Queries) CardDashboardWithdrawRepository {
	return &cardDashboardWithdrawRepository{
		db: db,
	}
}

// GetTotalWithdrawAmount retrieves the total amount of all withdrawals in the database.
//
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A pointer to an int64 containing the total withdrawal amount.
//   - An error if the query fails.
func (r *cardDashboardWithdrawRepository) GetTotalWithdrawAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountFailed
	}

	return &res, nil
}

// GetTotalWithdrawAmountByCardNumber retrieves the total withdraw amount for the card with the given card number.
//
// Parameters:
//   - ctx: the context for the database operation
//   - cardNumber: the card number to retrieve the withdraw amount for.
//
// Returns:
//   - A pointer to an int64 containing the total withdraw amount.
//   - An error if the query fails.
func (r *cardDashboardWithdrawRepository) GetTotalWithdrawAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountByCardFailed
	}

	return &res, nil
}

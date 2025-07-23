package repositorydashboard

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardDashboardTransactionRepository struct {
	db *db.Queries
}

func NewCardDashboardTransactionRepository(db *db.Queries) CardDashboardTransactionRepository {
	return &cardDashboardTransactionRepository{
		db: db,
	}
}

// GetTotalTransactionAmount retrieves the total transaction amount for all cards in the database.
//
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A pointer to an int64 containing the total transaction amount.
//   - An error if the query fails.
func (r *cardDashboardTransactionRepository) GetTotalTransactionAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTransactionAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountFailed
	}

	return &res, nil
}

// GetTotalTransactionAmountByCardNumber retrieves the total transaction amount for the card with the specified card number.
//
// Parameters:
//   - ctx: the context for the database operation
//   - cardNumber: the card number to retrieve the total transaction amount for.
//
// Returns:
//   - A pointer to an int64 containing the total transaction amount.
//   - An error if the query fails.
func (r *cardDashboardTransactionRepository) GetTotalTransactionAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransactionAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountByCardFailed
	}

	return &res, nil
}

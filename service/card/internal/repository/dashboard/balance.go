package repositorydashboard

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardDashboardBalanceRepository struct {
	db *db.Queries
}

func NewCardDashboardBalanceRepository(db *db.Queries) CardDashboardBalanceRepository {
	return &cardDashboardBalanceRepository{
		db: db,
	}
}

// GetTotalBalances retrieves the total balances for all cards in the database.
//
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A pointer to an int64 containing the total balance.
//   - An error if the query fails.
func (r *cardDashboardBalanceRepository) GetTotalBalances(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalBalance(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalancesFailed
	}

	return &res, nil
}

// GetTotalBalanceByCardNumber retrieves the total balance of the card with the given card number.
//
// Parameters:
//   - ctx: the context for the database operation
//   - cardNumber: the card number to retrieve the balance for.
//
// Returns:
//   - A pointer to an int64 containing the total balance.
//   - An error if the query fails.
func (r *cardDashboardBalanceRepository) GetTotalBalanceByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalBalanceByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalanceByCardFailed
	}

	return &res, nil
}

package repositorydashboard

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardDashboardTopupRepository struct {
	db *db.Queries
}

func NewCardDashboardTopupRepository(db *db.Queries) CardDashboardTopupRepository {
	return &cardDashboardTopupRepository{
		db: db,
	}
}

// GetTotalTopAmount retrieves the total amount of all topups in the database.
//
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A pointer to an int64 containing the total topup amount.
//   - An error if the query fails.
func (r *cardDashboardTopupRepository) GetTotalTopAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTopupAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopAmountFailed
	}

	return &res, nil
}

// GetTotalTopupAmountByCardNumber retrieves the total top-up amount for the card with the given card number.
//
// Parameters:
//   - ctx: the context for the database operation
//   - cardNumber: the card number to retrieve the top-up amount for.
//
// Returns:
//   - A pointer to an int64 containing the total top-up amount.
//   - An error if the query fails.
func (r *cardDashboardTopupRepository) GetTotalTopupAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTopupAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopupAmountByCardFailed
	}

	return &res, nil
}

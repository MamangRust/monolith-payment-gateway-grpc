package repositorydashboard

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardDashboardTransferRepository struct {
	db *db.Queries
}

func NewCardDashboardTransferRepository(db *db.Queries) CardDashboardTransferRepository {
	return &cardDashboardTransferRepository{
		db: db,
	}
}

// GetTotalTransferAmount retrieves the total amount of all transfers in the database.
//
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A pointer to an int64 containing the total transfer amount.
//   - An error if the query fails.
func (r *cardDashboardTransferRepository) GetTotalTransferAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTransferAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountFailed
	}

	return &res, nil
}

// GetTotalTransferAmountBySender retrieves the total transfer amount sent by the card with the specified sender card number.
//
// Parameters:
//   - ctx: the context for the database operation
//   - senderCardNumber: the card number of the sender to retrieve the total transfer amount for.
//
// Returns:
//   - A pointer to an int64 containing the total transfer amount.
//   - An error if the query fails.
func (r *cardDashboardTransferRepository) GetTotalTransferAmountBySender(ctx context.Context, senderCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountBySender(ctx, senderCardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountBySenderFailed
	}

	return &res, nil
}

// GetTotalTransferAmountByReceiver retrieves the total transfer amount received by the card with the specified receiver card number.
//
// Parameters:
//   - ctx: the context for the database operation
//   - receiverCardNumber: the card number of the receiver to retrieve the total transfer amount for.
//
// Returns:
//   - A pointer to an int64 containing the total transfer amount.
//   - An error if the query fails.
func (r *cardDashboardTransferRepository) GetTotalTransferAmountByReceiver(ctx context.Context, receiverCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountByReceiver(ctx, receiverCardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountByReceiverFailed
	}

	return &res, nil
}

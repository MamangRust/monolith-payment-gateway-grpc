package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardRepository struct {
	db *db.Queries
}

func NewCardRepository(db *db.Queries) CardRepository {
	return &cardRepository{
		db: db,
	}
}

func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return res, nil
}

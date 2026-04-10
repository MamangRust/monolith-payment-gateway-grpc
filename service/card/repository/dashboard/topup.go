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

func (r *cardDashboardTopupRepository) GetTotalTopAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTopupAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopAmountFailed
	}

	return &res, nil
}

func (r *cardDashboardTopupRepository) GetTotalTopupAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTopupAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopupAmountByCardFailed
	}

	return &res, nil
}

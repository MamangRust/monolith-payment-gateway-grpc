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

func (r *cardDashboardWithdrawRepository) GetTotalWithdrawAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawsFailed.WithInternal(err)
	}

	return &res, nil
}

func (r *cardDashboardWithdrawRepository) GetTotalWithdrawAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountByCardFailed.WithInternal(err)
	}

	return &res, nil
}

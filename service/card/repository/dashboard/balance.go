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

func (r *cardDashboardBalanceRepository) GetTotalBalances(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalBalance(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalancesFailed.WithInternal(err)
	}

	return &res, nil
}

func (r *cardDashboardBalanceRepository) GetTotalBalanceByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalBalanceByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalanceByCardFailed.WithInternal(err)
	}

	return &res, nil
}

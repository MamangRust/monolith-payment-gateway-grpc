package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStaticsTransferRepository struct {
	db *db.Queries
}

func NewCardStatsTransferRepository(db *db.Queries) CardStatsTransferRepository {
	return &cardStaticsTransferRepository{
		db: db,
	}
}

func (r *cardStaticsTransferRepository) GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountSenderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountSender(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountSenderFailed
	}

	return res, nil
}

func (r *cardStaticsTransferRepository) GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountSenderRow, error) {
	res, err := r.db.GetYearlyTransferAmountSender(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountSenderFailed
	}

	return res, nil
}

func (r *cardStaticsTransferRepository) GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountReceiverRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountReceiver(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountReceiverFailed.WithInternal(err)
	}

	return res, nil
}

func (r *cardStaticsTransferRepository) GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountReceiverRow, error) {
	res, err := r.db.GetYearlyTransferAmountReceiver(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountReceiverFailed.WithInternal(err)
	}

	return res, nil
}

package repositorystatsbycard

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardStatsTransferByCardRepository struct {
	db *db.Queries
}

func NewCardStatsTransferByCardRepository(db *db.Queries) CardStatsTransferByCardRepository {
	return &cardStatsTransferByCardRepository{
		db: db,
	}
}

func (r *cardStatsTransferByCardRepository) GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountBySenderRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountBySender(ctx, db.GetMonthlyTransferAmountBySenderParams{
		Column2:      yearStart,
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountBySenderFailed
	}

	return res, nil
}

func (r *cardStatsTransferByCardRepository) GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountBySenderRow, error) {
	res, err := r.db.GetYearlyTransferAmountBySender(ctx, db.GetYearlyTransferAmountBySenderParams{
		Column2:      int32(req.Year),
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountBySenderFailed
	}

	return res, nil
}

func (r *cardStatsTransferByCardRepository) GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountByReceiverRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountByReceiver(ctx, db.GetMonthlyTransferAmountByReceiverParams{
		Column2:    yearStart,
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountByReceiverFailed.WithInternal(err)
	}

	return res, nil
}

func (r *cardStatsTransferByCardRepository) GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountByReceiverRow, error) {
	res, err := r.db.GetYearlyTransferAmountByReceiver(ctx, db.GetYearlyTransferAmountByReceiverParams{
		Column2:    int32(req.Year),
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountByReceiverFailed.WithInternal(err)
	}

	return res, nil
}

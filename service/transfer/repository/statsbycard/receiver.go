package transferstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsAmountReceiverRepository struct {
	db *db.Queries
}

func NewTransferStatsAmountReceiverRepository(db *db.Queries) TransferStatsByCardAmountReceiverRepository {
	return &transferStatsAmountReceiverRepository{
		db: db,
	}
}

func (r *transferStatsAmountReceiverRepository) GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsByReceiverCardNumberRow, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsByReceiverCardFailed
	}
	return res, nil
}

func (r *transferStatsAmountReceiverRepository) GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsByReceiverCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsByReceiverCardFailed
	}

	return res, nil
}

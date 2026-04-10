package transferstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsAmountSenderRepository struct {
	db *db.Queries
}

func NewTransferStatsAmountSenderRepository(db *db.Queries) TransferStatsByCardAmountSenderRepository {
	return &transferStatsAmountSenderRepository{
		db: db,
	}
}

func (r *transferStatsAmountSenderRepository) GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsBySenderCardNumberRow, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsBySenderCardFailed
	}

	return res, nil
}

func (r *transferStatsAmountSenderRepository) GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsBySenderCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsBySenderCardFailed
	}

	return res, nil
}

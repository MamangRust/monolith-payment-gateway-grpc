package transferstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsStatusRepository struct {
	db *db.Queries
}

func NewTransferStatsStatusRepository(db *db.Queries) TransferStatsStatusRepository {
	return &transferStatsStatusRepository{
		db: db,
	}
}

func (r *transferStatsStatusRepository) GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccess(ctx, db.GetMonthTransferStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessFailed
	}

	return res, nil
}

func (r *transferStatsStatusRepository) GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusSuccessRow, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessFailed
	}

	return res, nil
}

func (r *transferStatsStatusRepository) GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailed(ctx, db.GetMonthTransferStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedFailed
	}

	return res, nil
}

func (r *transferStatsStatusRepository) GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusFailedRow, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedFailed
	}

	return res, nil
}

package transferstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsStatusByCardRepository struct {
	db *db.Queries
}

func NewTransferStatsByCardStatusRepository(db *db.Queries) TransferStatsByCardStatusRepository {
	return &transferStatsStatusByCardRepository{
		db: db,
	}
}

func (r *transferStatsStatusByCardRepository) GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transferStatsStatusByCardRepository) GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transferStatsStatusByCardRepository) GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusFailedCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transferStatsStatusByCardRepository) GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedByCardFailed.WithInternal(err)
	}

	return res, nil
}

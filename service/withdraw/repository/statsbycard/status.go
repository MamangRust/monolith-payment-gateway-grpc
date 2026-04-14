package withdrawstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawStatsByCardStatusRepository struct {
	db *db.Queries
}

func NewWithdrawStatsStatusRepository(db *db.Queries) WithdrawStatsByCardStatusRepository {
	return &withdrawStatsByCardStatusRepository{
		db: db,
	}
}

func (r *withdrawStatsByCardStatusRepository) GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccessCardNumber(ctx, db.GetMonthWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *withdrawStatsByCardStatusRepository) GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccessCardNumber(ctx, db.GetYearlyWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *withdrawStatsByCardStatusRepository) GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusFailedCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailedCardNumber(ctx, db.GetMonthWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *withdrawStatsByCardStatusRepository) GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailedCardNumber(ctx, db.GetYearlyWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedByCardFailed.WithInternal(err)
	}

	return res, nil
}

package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type withdrawStatisticsRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.WithdrawRecordMapping
}

func NewWithdrawStatisticRepository(db *db.Queries, ctx context.Context, mapping recordmapper.WithdrawRecordMapping) *withdrawStatisticsRepository {
	return &withdrawStatisticsRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *withdrawStatisticsRepository) GetMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccess(r.ctx, db.GetMonthWithdrawStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessFailed
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *withdrawStatisticsRepository) GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessFailed
	}

	so := r.mapping.ToWithdrawRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *withdrawStatisticsRepository) GetMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailed(r.ctx, db.GetMonthWithdrawStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedFailed
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *withdrawStatisticsRepository) GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedFailed
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailed(res)

	return so, nil
}

func (r *withdrawStatisticsRepository) GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(r.ctx, yearStart)

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsFailed
	}

	return r.mapping.ToWithdrawsAmountMonthly(res), nil

}

func (r *withdrawStatisticsRepository) GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdraws(r.ctx, year)

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsFailed
	}

	return r.mapping.ToWithdrawsAmountYearly(res), nil

}

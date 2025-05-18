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

type withdrawStatisticByCardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.WithdrawRecordMapping
}

func NewWithdrawStatisticByCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.WithdrawRecordMapping) *withdrawStatisticByCardRepository {
	return &withdrawStatisticByCardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *withdrawStatisticByCardRepository) GetMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccessCardNumber(r.ctx, db.GetMonthWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessByCardFailed
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *withdrawStatisticByCardRepository) GetYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccessCardNumber(r.ctx, db.GetYearlyWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessByCardFailed
	}

	so := r.mapping.ToWithdrawRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *withdrawStatisticByCardRepository) GetMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailedCardNumber(r.ctx, db.GetMonthWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedByCardFailed
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawStatisticByCardRepository) GetYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailedCardNumber(r.ctx, db.GetYearlyWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedByCardFailed
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawStatisticByCardRepository) GetMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawsByCardNumber(r.ctx, db.GetMonthlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsByCardFailed
	}

	return r.mapping.ToWithdrawsAmountMonthlyByCardNumber(res), nil

}

func (r *withdrawStatisticByCardRepository) GetYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdrawsByCardNumber(r.ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsByCardFailed
	}

	return r.mapping.ToWithdrawsAmountYearlyByCardNumber(res), nil
}

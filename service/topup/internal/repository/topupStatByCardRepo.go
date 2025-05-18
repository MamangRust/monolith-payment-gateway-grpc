package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type topupStatisticByCardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TopupRecordMapping
}

func NewTopupStatisticByCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TopupRecordMapping) *topupStatisticByCardRepository {
	return &topupStatisticByCardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}
func (r *topupStatisticByCardRepository) GetMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccessCardNumber(r.ctx, db.GetMonthTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessByCardFailed
	}

	so := r.mapping.ToTopupRecordsMonthStatusSuccessByCardNumber(res)

	return so, nil
}

func (r *topupStatisticByCardRepository) GetYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccessCardNumber(r.ctx, db.GetYearlyTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed
	}

	so := r.mapping.ToTopupRecordsYearStatusSuccessByCardNumber(res)

	return so, nil
}

func (r *topupStatisticByCardRepository) GetMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusFailed, error) {
	cardNumber := req.CardNumber
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailedCardNumber(r.ctx, db.GetMonthTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedByCardFailed
	}

	so := r.mapping.ToTopupRecordsMonthStatusFailedByCardNumber(res)

	return so, nil
}

func (r *topupStatisticByCardRepository) GetYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusFailed, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyTopupStatusFailedCardNumber(r.ctx, db.GetYearlyTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusFailedByCardFailed
	}

	so := r.mapping.ToTopupRecordsYearStatusFailedByCardNumber(res)

	return so, nil
}

func (r *topupStatisticByCardRepository) GetMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethodsByCardNumber(r.ctx, db.GetMonthlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsByCardFailed
	}

	return r.mapping.ToTopupMonthlyMethodsByCardNumber(res), nil
}

func (r *topupStatisticByCardRepository) GetYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupMethodsByCardNumber(r.ctx, db.GetYearlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsByCardFailed
	}

	return r.mapping.ToTopupYearlyMethodsByCardNumber(res), nil
}

func (r *topupStatisticByCardRepository) GetMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthAmount, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountsByCardNumber(r.ctx, db.GetMonthlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsByCardFailed
	}

	return r.mapping.ToTopupMonthlyAmountsByCardNumber(res), nil
}

func (r *topupStatisticByCardRepository) GetYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyAmount, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupAmountsByCardNumber(r.ctx, db.GetYearlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsByCardFailed
	}

	return r.mapping.ToTopupYearlyAmountsByCardNumber(res), nil
}

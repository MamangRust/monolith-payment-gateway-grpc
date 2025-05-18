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

type topupStatisticRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TopupRecordMapping
}

func NewTopupStatisticRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TopupRecordMapping) *topupStatisticRepository {
	return &topupStatisticRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *topupStatisticRepository) GetMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccess(r.ctx, db.GetMonthTopupStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessFailed
	}

	so := r.mapping.ToTopupRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *topupStatisticRepository) GetYearlyTopupStatusSuccess(year int) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessFailed
	}
	so := r.mapping.ToTopupRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *topupStatisticRepository) GetMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailed(r.ctx, db.GetMonthTopupStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedFailed
	}

	so := r.mapping.ToTopupRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *topupStatisticRepository) GetYearlyTopupStatusFailed(year int) ([]*record.TopupRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTopupStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusFailedFailed
	}

	so := r.mapping.ToTopupRecordsYearStatusFailed(res)

	return so, nil
}

func (r *topupStatisticRepository) GetMonthlyTopupMethods(year int) ([]*record.TopupMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethods(r.ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsFailed
	}

	return r.mapping.ToTopupMonthlyMethods(res), nil
}

func (r *topupStatisticRepository) GetYearlyTopupMethods(year int) ([]*record.TopupYearlyMethod, error) {
	res, err := r.db.GetYearlyTopupMethods(r.ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsFailed
	}

	return r.mapping.ToTopupYearlyMethods(res), nil
}

func (r *topupStatisticRepository) GetMonthlyTopupAmounts(year int) ([]*record.TopupMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmounts(r.ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsFailed
	}

	return r.mapping.ToTopupMonthlyAmounts(res), nil
}

func (r *topupStatisticRepository) GetYearlyTopupAmounts(year int) ([]*record.TopupYearlyAmount, error) {
	res, err := r.db.GetYearlyTopupAmounts(r.ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsFailed
	}

	return r.mapping.ToTopupYearlyAmounts(res), nil
}

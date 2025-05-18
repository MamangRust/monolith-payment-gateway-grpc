package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transferStatisticRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransferRecordMapping
}

func NewTransferStatisticRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransferRecordMapping) *transferStatisticRepository {
	return &transferStatisticRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transferStatisticRepository) GetMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccess(r.ctx, db.GetMonthTransferStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transferStatisticRepository) GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transferStatisticRepository) GetMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailed(r.ctx, db.GetMonthTransferStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transferStatisticRepository) GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transferStatisticRepository) GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(r.ctx, yearStart)

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsFailed
	}

	return r.mapping.ToTransferMonthAmounts(res), nil
}

func (r *transferStatisticRepository) GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmounts(r.ctx, year)

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsFailed
	}
	return r.mapping.ToTransferYearAmounts(res), nil
}

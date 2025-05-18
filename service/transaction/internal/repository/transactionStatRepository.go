package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transactionStatisticRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionStatisticRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionStatisticRepository {
	return &transactionStatisticRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionStatisticRepository) GetMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccess(r.ctx, db.GetMonthTransactionStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transactionStatisticRepository) GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transactionStatisticRepository) GetMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailed(r.ctx, db.GetMonthTransactionStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transactionStatisticRepository) GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transactionStatisticRepository) GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(r.ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsFailed
	}

	return r.mapping.ToTransactionMonthlyMethods(res), nil
}

func (r *transactionStatisticRepository) GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error) {
	res, err := r.db.GetYearlyPaymentMethods(r.ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsFailed
	}

	return r.mapping.ToTransactionYearlyMethods(res), nil
}

func (r *transactionStatisticRepository) GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(r.ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsFailed
	}

	return r.mapping.ToTransactionMonthlyAmounts(res), nil
}

func (r *transactionStatisticRepository) GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error) {
	res, err := r.db.GetYearlyAmounts(r.ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsFailed
	}

	return r.mapping.ToTransactionYearlyAmounts(res), nil
}

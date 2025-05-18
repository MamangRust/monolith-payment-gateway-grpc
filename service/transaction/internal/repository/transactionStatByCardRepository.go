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

type transactionStatisticByCardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionStatisticByCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionStatisticByCardRepository {
	return &transactionStatisticByCardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionStatisticByCardRepository) GetMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccessCardNumber(r.ctx, db.GetMonthTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transactionStatisticByCardRepository) GetYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccessCardNumber(r.ctx, db.GetYearlyTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transactionStatisticByCardRepository) GetMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailedCardNumber(r.ctx, db.GetMonthTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedByCardFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionStatisticByCardRepository) GetYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailedCardNumber(r.ctx, db.GetYearlyTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedByCardFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionStatisticByCardRepository) GetMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetMonthlyPaymentMethodsByCardNumber(r.ctx, db.GetMonthlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsByCardFailed
	}

	return r.mapping.ToTransactionMonthlyMethodsByCardNumber(res), nil
}

func (r *transactionStatisticByCardRepository) GetYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(r.ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsByCardFailed
	}

	return r.mapping.ToTransactionYearlyMethodsByCardNumber(res), nil
}

func (r *transactionStatisticByCardRepository) GetMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetMonthlyAmountsByCardNumber(r.ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsByCardFailed
	}

	return r.mapping.ToTransactionMonthlyAmountsByCardNumber(res), nil
}

func (r *transactionStatisticByCardRepository) GetYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyAmountsByCardNumber(r.ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsByCardFailed
	}

	return r.mapping.ToTransactionYearlyAmountsByCardNumber(res), nil
}

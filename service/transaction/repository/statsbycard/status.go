package transactionbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsByCardStatusRepository struct {
	db *db.Queries
}

func NewTransactionStatsByCardStatusRepository(db *db.Queries) TransactonStatsByCardStatusRepository {
	return &transactionStatsByCardStatusRepository{
		db: db,
	}
}

func (r *transactionStatsByCardStatusRepository) GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccessCardNumber(ctx, db.GetMonthTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByCardStatusRepository) GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccessCardNumber(ctx, db.GetYearlyTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByCardStatusRepository) GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusFailedCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailedCardNumber(ctx, db.GetMonthTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedByCardFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByCardStatusRepository) GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyTransactionStatusFailedCardNumber(ctx, db.GetYearlyTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedByCardFailed.WithInternal(err)
	}

	return res, nil
}

package transactionbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/statsbycard"
)

type transactionStatsByCardStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionStatisticByCardStatusMapper
}

func NewTransactionStatsByCardStatusRepository(db *db.Queries, mapper recordmapper.TransactionStatisticByCardStatusMapper) TransactonStatsByCardStatusRepository {
	return &transactionStatsByCardStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthTransactionStatusSuccessByCardNumber retrieves monthly successful transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and card number.
//
// Returns:
//   - []*record.TransactionRecordMonthStatusSuccess: List of monthly success transaction stats.
//   - error: Error if any occurs.
func (r *transactionStatsByCardStatusRepository) GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusSuccess, error) {
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
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessByCardFailed
	}

	so := r.mapper.ToTransactionRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

// GetYearlyTransactionStatusSuccessByCardNumber retrieves yearly successful transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TransactionRecordYearStatusSuccess: List of yearly success transaction stats.
//   - error: Error if any occurs.
func (r *transactionStatsByCardStatusRepository) GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccessCardNumber(ctx, db.GetYearlyTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessByCardFailed
	}

	so := r.mapper.ToTransactionRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

// GetMonthTransactionStatusFailedByCardNumber retrieves monthly failed transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and card number.
//
// Returns:
//   - []*record.TransactionRecordMonthStatusFailed: List of monthly failed transaction stats.
//   - error: Error if any occurs.
func (r *transactionStatsByCardStatusRepository) GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusFailed, error) {
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
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedByCardFailed
	}

	so := r.mapper.ToTransactionRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

// GetYearlyTransactionStatusFailedByCardNumber retrieves yearly failed transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TransactionRecordYearStatusFailed: List of yearly failed transaction stats.
//   - error: Error if any occurs.
func (r *transactionStatsByCardStatusRepository) GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailedCardNumber(ctx, db.GetYearlyTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedByCardFailed
	}

	so := r.mapper.ToTransactionRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

package transactionstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/stats"
)

type transactionStatsStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionStatisticStatusRecordMapper
}

func NewTransactionStatsStatusRepository(db *db.Queries, mapper recordmapper.TransactionStatisticStatusRecordMapper) TransactionStatsStatusRepository {
	return &transactionStatsStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthTransactionStatusSuccess retrieves monthly statistics of successful transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and status filter.
//
// Returns:
//   - []*record.TransactionRecordMonthStatusSuccess: List of monthly successful transaction statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsStatusRepository) GetMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccess(ctx, db.GetMonthTransactionStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessFailed
	}

	so := r.mapper.ToTransactionRecordsMonthStatusSuccess(res)

	return so, nil
}

// GetYearlyTransactionStatusSuccess retrieves yearly statistics of successful transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.TransactionRecordYearStatusSuccess: List of yearly successful transaction statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsStatusRepository) GetYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessFailed
	}

	so := r.mapper.ToTransactionRecordsYearStatusSuccess(res)

	return so, nil
}

// GetMonthTransactionStatusFailed retrieves monthly statistics of failed transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and status filter.
//
// Returns:
//   - []*record.TransactionRecordMonthStatusFailed: List of monthly failed transaction statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsStatusRepository) GetMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailed(ctx, db.GetMonthTransactionStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedFailed
	}

	so := r.mapper.ToTransactionRecordsMonthStatusFailed(res)

	return so, nil
}

// GetYearlyTransactionStatusFailed retrieves yearly statistics of failed transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.TransactionRecordYearStatusFailed: List of yearly failed transaction statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsStatusRepository) GetYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedFailed
	}

	so := r.mapper.ToTransactionRecordsYearStatusFailed(res)

	return so, nil
}

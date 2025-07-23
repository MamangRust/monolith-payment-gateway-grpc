package transferstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer/stats"
)

type transferStatsStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferStatisticStatusRecordMapper
}

func NewTransferStatsStatusRepository(db *db.Queries, mapper recordmapper.TransferStatisticStatusRecordMapper) TransferStatsStatusRepository {
	return &transferStatsStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthTransferStatusSuccess retrieves successful transfer statistics per month.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The month and year filter for the statistics.
//
// Returns:
//   - []*record.TransferRecordMonthStatusSuccess: List of monthly successful transfer records.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusRepository) GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccess(ctx, db.GetMonthTransferStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessFailed
	}

	so := r.mapper.ToTransferRecordsMonthStatusSuccess(res)

	return so, nil
}

// GetYearlyTransferStatusSuccess retrieves successful transfer statistics per year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*record.TransferRecordYearStatusSuccess: List of yearly successful transfer records.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusRepository) GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessFailed
	}

	so := r.mapper.ToTransferRecordsYearStatusSuccess(res)

	return so, nil
}

// GetMonthTransferStatusFailed retrieves failed transfer statistics per month.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The month and year filter for the statistics.
//
// Returns:
//   - []*record.TransferRecordMonthStatusFailed: List of monthly failed transfer records.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusRepository) GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailed(ctx, db.GetMonthTransferStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedFailed
	}

	so := r.mapper.ToTransferRecordsMonthStatusFailed(res)

	return so, nil
}

// GetYearlyTransferStatusFailed retrieves failed transfer statistics per year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*record.TransferRecordYearStatusFailed: List of yearly failed transfer records.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusRepository) GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedFailed
	}

	so := r.mapper.ToTransferRecordsYearStatusFailed(res)

	return so, nil
}

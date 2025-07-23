package topupstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/stats"
)

type topupStatsStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupStatisticStatusMapper
}

func NewTopupStatsStatusRepository(db *db.Queries, mapper recordmapper.TopupStatisticStatusMapper) TopupStatsStatusRepository {
	return &topupStatsStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthTopupStatusSuccess retrieves monthly statistics of successful topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and method for filtering.
//
// Returns:
//   - []*record.TopupRecordMonthStatusSuccess: List of monthly successful topup records.
//   - error: Error if the query fails.
func (r *topupStatsStatusRepository) GetMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccess(ctx, db.GetMonthTopupStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessFailed
	}

	so := r.mapper.ToTopupRecordsMonthStatusSuccess(res)

	return so, nil
}

// GetYearlyTopupStatusSuccess retrieves yearly statistics of successful topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*record.TopupRecordYearStatusSuccess: List of yearly successful topup records.
//   - error: Error if the query fails.
func (r *topupStatsStatusRepository) GetYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessFailed
	}
	so := r.mapper.ToTopupRecordsYearStatusSuccess(res)

	return so, nil
}

// GetMonthTopupStatusFailed retrieves monthly statistics of failed topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and method for filtering.
//
// Returns:
//   - []*record.TopupRecordMonthStatusFailed: List of monthly failed topup records.
//   - error: Error if the query fails.
func (r *topupStatsStatusRepository) GetMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailed(ctx, db.GetMonthTopupStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedFailed
	}

	so := r.mapper.ToTopupRecordsMonthStatusFailed(res)

	return so, nil
}

// GetYearlyTopupStatusFailed retrieves yearly statistics of failed topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*record.TopupRecordYearStatusFailed: List of yearly failed topup records.
//   - error: Error if the query fails.
func (r *topupStatsStatusRepository) GetYearlyTopupStatusFailed(ctx context.Context, year int) ([]*record.TopupRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTopupStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusFailedFailed
	}

	so := r.mapper.ToTopupRecordsYearStatusFailed(res)

	return so, nil
}

package withdrawstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw/stats"
)

type withdrawStatsStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.WithdrawStatisticStatusRecordMapper
}

func NewWithdrawStatsStatusRepository(db *db.Queries, mapper recordmapper.WithdrawStatisticStatusRecordMapper) WithdrawStatsStatusRepository {
	return &withdrawStatsStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthWithdrawStatusSuccess retrieves monthly withdraw statistics with status "success".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and additional filters.
//
// Returns:
//   - []*record.WithdrawRecordMonthStatusSuccess: List of successful monthly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsStatusRepository) GetMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccess(ctx, db.GetMonthWithdrawStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessFailed
	}

	so := r.mapper.ToWithdrawRecordsMonthStatusSuccess(res)

	return so, nil
}

// GetYearlyWithdrawStatusSuccess retrieves yearly withdraw statistics with status "success".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.WithdrawRecordYearStatusSuccess: List of successful yearly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsStatusRepository) GetYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessFailed
	}

	so := r.mapper.ToWithdrawRecordsYearStatusSuccess(res)

	return so, nil
}

// GetMonthWithdrawStatusFailed retrieves monthly withdraw statistics with status "failed".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and additional filters.
//
// Returns:
//   - []*record.WithdrawRecordMonthStatusFailed: List of failed monthly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsStatusRepository) GetMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailed(ctx, db.GetMonthWithdrawStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedFailed
	}

	so := r.mapper.ToWithdrawRecordsMonthStatusFailed(res)

	return so, nil
}

// GetYearlyWithdrawStatusFailed retrieves yearly withdraw statistics with status "failed".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.WithdrawRecordYearStatusFailed: List of failed yearly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsStatusRepository) GetYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedFailed
	}

	so := r.mapper.ToWithdrawRecordsYearStatusFailed(res)

	return so, nil
}

package withdrawstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw/statsbycard"
)

type withdrawStatsByCardStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.WithdrawStatisticStatusByCardRecordMapper
}

func NewWithdrawStatsStatusRepository(db *db.Queries, mapper recordmapper.WithdrawStatisticStatusByCardRecordMapper) WithdrawStatsByCardStatusRepository {
	return &withdrawStatsByCardStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthWithdrawStatusSuccessByCardNumber retrieves monthly withdraw success statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number, month, and year.
//
// Returns:
//   - []*record.WithdrawRecordMonthStatusSuccess: List of successful monthly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsByCardStatusRepository) GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccessCardNumber(ctx, db.GetMonthWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessByCardFailed
	}

	so := r.mapper.ToWithdrawRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

// GetYearlyWithdrawStatusSuccessByCardNumber retrieves yearly withdraw success statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and year.
//
// Returns:
//   - []*record.WithdrawRecordYearStatusSuccess: List of successful yearly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsByCardStatusRepository) GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccessCardNumber(ctx, db.GetYearlyWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessByCardFailed
	}

	so := r.mapper.ToWithdrawRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

// GetMonthWithdrawStatusFailedByCardNumber retrieves monthly withdraw failed statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number, month, and year.
//
// Returns:
//   - []*record.WithdrawRecordMonthStatusFailed: List of failed monthly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsByCardStatusRepository) GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailedCardNumber(ctx, db.GetMonthWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedByCardFailed
	}

	so := r.mapper.ToWithdrawRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

// GetYearlyWithdrawStatusFailedByCardNumber retrieves yearly withdraw failed statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and year.
//
// Returns:
//   - []*record.WithdrawRecordYearStatusFailed: List of failed yearly withdraw records.
//   - error: An error if the operation fails.
func (r *withdrawStatsByCardStatusRepository) GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailedCardNumber(ctx, db.GetYearlyWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedByCardFailed
	}

	so := r.mapper.ToWithdrawRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

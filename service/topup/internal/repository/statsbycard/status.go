package topupstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup/statsbycard"
)

type topupStatsByCardStatusRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupStatisticStatusByCardNumberMapper
}

func NewTopupStatsByCardStatusRepository(db *db.Queries, mapper recordmapper.TopupStatisticStatusByCardNumberMapper) TopupStatsByCardStatusRepository {
	return &topupStatsByCardStatusRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthTopupStatusSuccessByCardNumber retrieves monthly statistics of successful topups for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and card number.
//
// Returns:
//   - []*record.TopupRecordMonthStatusSuccess: List of monthly successful topup records.
//   - error: Error if the query fails.
func (r *topupStatsByCardStatusRepository) GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)

	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccessCardNumber(ctx, db.GetMonthTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessByCardFailed
	}

	so := r.mapper.ToTopupRecordsMonthStatusSuccessByCardNumber(res)

	return so, nil
}

// GetYearlyTopupStatusSuccessByCardNumber retrieves yearly statistics of successful topups for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TopupRecordYearStatusSuccess: List of yearly successful topup records.
//   - error: Error if the query fails.
func (r *topupStatsByCardStatusRepository) GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccessCardNumber(ctx, db.GetYearlyTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed
	}

	so := r.mapper.ToTopupRecordsYearStatusSuccessByCardNumber(res)

	return so, nil
}

// GetMonthTopupStatusFailedByCardNumber retrieves monthly statistics of failed topups for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and card number.
//
// Returns:
//   - []*record.TopupRecordMonthStatusFailed: List of monthly failed topup records.
//   - error: Error if the query fails.
func (r *topupStatsByCardStatusRepository) GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusFailed, error) {
	cardNumber := req.CardNumber
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailedCardNumber(ctx, db.GetMonthTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedByCardFailed
	}

	so := r.mapper.ToTopupRecordsMonthStatusFailedByCardNumber(res)

	return so, nil
}

// GetYearlyTopupStatusFailedByCardNumber retrieves yearly statistics of failed topups for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*record.TopupRecordYearStatusFailed: List of yearly failed topup records.
//   - error: Error if the query fails.
func (r *topupStatsByCardStatusRepository) GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusFailed, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyTopupStatusFailedCardNumber(ctx, db.GetYearlyTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusFailedByCardFailed
	}

	so := r.mapper.ToTopupRecordsYearStatusFailedByCardNumber(res)

	return so, nil
}

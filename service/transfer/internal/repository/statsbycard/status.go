package transferstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer/statsbycard"
)

type transferStatsStatusByCardRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferStatisticStatusByCardRecordMapper
}

func NewTransferStatsByCardStatusRepository(db *db.Queries, mapper recordmapper.TransferStatisticStatusByCardRecordMapper) TransferStatsByCardStatusRepository {
	return &transferStatsStatusByCardRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthTransferStatusSuccessByCardNumber retrieves monthly successful transfer statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and date filters.
//
// Returns:
//   - []*record.TransferRecordMonthStatusSuccess: List of monthly successful transfer statistics.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusByCardRepository) GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessByCardFailed
	}

	so := r.mapper.ToTransferRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

// GetYearlyTransferStatusSuccessByCardNumber retrieves yearly successful transfer statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*record.TransferRecordYearStatusSuccess: List of yearly successful transfer statistics.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusByCardRepository) GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessByCardFailed
	}

	so := r.mapper.ToTransferRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

// GetMonthTransferStatusFailedByCardNumber retrieves monthly failed transfer statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and date filters.
//
// Returns:
//   - []*record.TransferRecordMonthStatusFailed: List of monthly failed transfer statistics.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusByCardRepository) GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedByCardFailed
	}

	so := r.mapper.ToTransferRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

// GetYearlyTransferStatusFailedByCardNumber retrieves yearly failed transfer statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*record.TransferRecordYearStatusFailed: List of yearly failed transfer statistics.
//   - error: Any error encountered during the operation.
func (r *transferStatsStatusByCardRepository) GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedByCardFailed
	}

	so := r.mapper.ToTransferRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

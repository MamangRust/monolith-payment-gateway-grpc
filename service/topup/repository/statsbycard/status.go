package topupstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsByCardStatusRepository struct {
	db *db.Queries
}

func NewTopupStatsByCardStatusRepository(db *db.Queries) TopupStatsByCardStatusRepository {
	return &topupStatsByCardStatusRepository{
		db: db,
	}
}

func (r *topupStatsByCardStatusRepository) GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusSuccessCardNumberRow, error) {
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

	return res, nil
}

func (r *topupStatsByCardStatusRepository) GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyTopupStatusSuccessCardNumber(ctx, db.GetYearlyTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *topupStatsByCardStatusRepository) GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusFailedCardNumberRow, error) {
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

	return res, nil
}

func (r *topupStatsByCardStatusRepository) GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyTopupStatusFailedCardNumber(ctx, db.GetYearlyTopupStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed
	}

	return res, nil
}

package withdrawstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawStatsByCardAmountRepository struct {
	db *db.Queries
}

func NewWithdrawStatsByCardAmountRepository(db *db.Queries) WithdrawStatsByCardAmountRepository {
	return &withdrawStatsByCardAmountRepository{
		db: db,
	}
}


func (r *withdrawStatsByCardAmountRepository) GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawsByCardNumber(ctx, db.GetMonthlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsByCardFailed.WithInternal(err)
	}

	return res, nil

}

func (r *withdrawStatsByCardAmountRepository) GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawsByCardNumber(ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsByCardFailed.WithInternal(err)
	}

	return res, nil
}

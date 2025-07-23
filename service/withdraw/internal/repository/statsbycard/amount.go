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

type withdrawStatsByCardAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.WithdrawStatisticAmountByCardRecordMapper
}

func NewWithdrawStatsAmountRepository(db *db.Queries, mapper recordmapper.WithdrawStatisticAmountByCardRecordMapper) WithdrawStatsByCardAmountRepository {
	return &withdrawStatsByCardAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyWithdrawsByCardNumber retrieves total monthly withdraw amounts by card number for a specific year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and year.
//
// Returns:
//   - []*record.WithdrawMonthlyAmount: List of monthly withdraw amount records.
//   - error: An error if the operation fails.
func (r *withdrawStatsByCardAmountRepository) GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawsByCardNumber(ctx, db.GetMonthlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsByCardFailed
	}

	so := r.mapper.ToWithdrawsAmountMonthlyByCardNumber(res)

	return so, nil
}

// GetYearlyWithdrawsByCardNumber retrieves total yearly withdraw amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and year.
//
// Returns:
//   - []*record.WithdrawYearlyAmount: List of yearly withdraw amount records.
//   - error: An error if the operation fails.
func (r *withdrawStatsByCardAmountRepository) GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdrawsByCardNumber(ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsByCardFailed
	}

	so := r.mapper.ToWithdrawsAmountYearlyByCardNumber(res)

	return so, nil
}

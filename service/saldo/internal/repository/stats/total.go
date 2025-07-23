package saldostatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

type saldoStatsTotalBalanceRepository struct {
	db     *db.Queries
	mapper recordmapper.SaldoStatisticRecordMapping
}

func NewSaldoStatsTotalBalanceRepository(db *db.Queries, mapper recordmapper.SaldoStatisticRecordMapping) SaldoStatsTotalSaldoRepository {
	return &saldoStatsTotalBalanceRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTotalSaldoBalance retrieves the total saldo balance grouped by month
// based on the given request (e.g., year, card number).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing the year and other filters.
//
// Returns:
//   - []*record.SaldoMonthTotalBalance: List of saldo totals per month.
//   - error: An error if the query fails.

func (r *saldoStatsTotalBalanceRepository) GetMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*record.SaldoMonthTotalBalance, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSaldoBalance(ctx, db.GetMonthlyTotalSaldoBalanceParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlyTotalSaldoBalanceFailed
	}

	so := r.mapper.ToSaldoMonthTotalBalances(res)
	return so, nil
}

// GetYearTotalSaldoBalance retrieves the total saldo balance for a given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The target year for the statistics.
//
// Returns:
//   - []*record.SaldoYearTotalBalance: List of saldo totals for the year.
//   - error: An error if the query fails.
func (r *saldoStatsTotalBalanceRepository) GetYearTotalSaldoBalance(ctx context.Context, year int) ([]*record.SaldoYearTotalBalance, error) {
	res, err := r.db.GetYearlyTotalSaldoBalances(ctx, int32(year))

	if err != nil {
		return nil, saldo_errors.ErrGetYearTotalSaldoBalanceFailed
	}

	so := r.mapper.ToSaldoYearTotalBalances(res)

	return so, nil
}

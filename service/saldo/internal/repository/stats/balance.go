package saldostatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

type saldoStatsBalanceRepository struct {
	db     *db.Queries
	mapper recordmapper.SaldoStatisticRecordMapping
}

func NewSaldoStatsBalanceRepository(db *db.Queries, mapper recordmapper.SaldoStatisticRecordMapping) SaldoStatsBalanceRepository {
	return &saldoStatsBalanceRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlySaldoBalances retrieves saldo balances for each month in a given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The target year to retrieve monthly balances.
//
// Returns:
//   - []*record.SaldoMonthSaldoBalance: List of saldo balances by month.
//   - error: An error if the query fails.
func (r *saldoStatsBalanceRepository) GetMonthlySaldoBalances(ctx context.Context, year int) ([]*record.SaldoMonthSaldoBalance, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlySaldoBalances(ctx, yearStart)

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlySaldoBalancesFailed
	}

	so := r.mapper.ToSaldoMonthBalances(res)

	return so, nil
}

// GetYearlySaldoBalances retrieves saldo balances aggregated per year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The target year to retrieve yearly balances.
//
// Returns:
//   - []*record.SaldoYearSaldoBalance: List of saldo balances by year.
//   - error: An error if the query fails.
func (r *saldoStatsBalanceRepository) GetYearlySaldoBalances(ctx context.Context, year int) ([]*record.SaldoYearSaldoBalance, error) {
	res, err := r.db.GetYearlySaldoBalances(ctx, year)

	if err != nil {
		return nil, saldo_errors.ErrGetYearlySaldoBalancesFailed
	}

	so := r.mapper.ToSaldoYearSaldoBalances(res)

	return so, nil
}

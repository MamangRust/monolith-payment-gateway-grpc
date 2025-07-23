package saldostatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoStatsBalanceRepository interface {
	// GetMonthlySaldoBalances retrieves saldo balances for each month in a given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The target year to retrieve monthly balances.
	//
	// Returns:
	//   - []*record.SaldoMonthSaldoBalance: List of saldo balances by month.
	//   - error: An error if the query fails.
	GetMonthlySaldoBalances(ctx context.Context, year int) ([]*record.SaldoMonthSaldoBalance, error)

	// GetYearlySaldoBalances retrieves saldo balances aggregated per year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The target year to retrieve yearly balances.
	//
	// Returns:
	//   - []*record.SaldoYearSaldoBalance: List of saldo balances by year.
	//   - error: An error if the query fails.
	GetYearlySaldoBalances(ctx context.Context, year int) ([]*record.SaldoYearSaldoBalance, error)
}

type SaldoStatsTotalSaldoRepository interface {
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
	GetMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*record.SaldoMonthTotalBalance, error)

	// GetYearTotalSaldoBalance retrieves the total saldo balance for a given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The target year for the statistics.
	//
	// Returns:
	//   - []*record.SaldoYearTotalBalance: List of saldo totals for the year.
	//   - error: An error if the query fails.
	GetYearTotalSaldoBalance(ctx context.Context, year int) ([]*record.SaldoYearTotalBalance, error)
}

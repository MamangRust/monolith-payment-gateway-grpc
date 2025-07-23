package saldostatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type SaldoStatsTotalBalanceService interface {
	// FindMonthlyTotalSaldoBalance retrieves the total saldo balance grouped by month
	// based on the provided request, which may include year and optional filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filter criteria (e.g., year, card number).
	//
	// Returns:
	//   - []*response.SaldoMonthTotalBalanceResponse: List of monthly total saldo balances.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse)

	// FindYearTotalSaldoBalance retrieves the total saldo balance aggregated for a specific year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The target year to retrieve the total saldo balance.
	//
	// Returns:
	//   - []*response.SaldoYearTotalBalanceResponse: List of yearly total saldo balances.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindYearTotalSaldoBalance(ctx context.Context, year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse)
}

type SaldoStatsBalanceService interface {
	// FindMonthlySaldoBalances retrieves saldo balances for each month in the specified year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The target year to retrieve monthly saldo balances.
	//
	// Returns:
	//   - []*response.SaldoMonthBalanceResponse: List of saldo balances per month.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindMonthlySaldoBalances(ctx context.Context, year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse)

	// FindYearlySaldoBalances retrieves saldo balances aggregated by year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The target year to retrieve yearly saldo balances.
	//
	// Returns:
	//   - []*response.SaldoYearBalanceResponse: List of saldo balances per year.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindYearlySaldoBalances(ctx context.Context, year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse)
}

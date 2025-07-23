package saldostatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type SaldoStatsTotalCache interface {
	// GetMonthlyTotalSaldoBalanceCache retrieves cached total saldo balance per month based on request filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the year, month, and additional filters.
	//
	// Returns:
	//   - []*response.SaldoMonthTotalBalanceResponse: The list of monthly total saldo balances.
	//   - bool: Whether the cache was found and valid.
	GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, bool)

	// SetMonthlyTotalSaldoCache stores total saldo balance per month in cache based on request filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to be cached.
	SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data []*response.SaldoMonthTotalBalanceResponse)

	// GetYearTotalSaldoBalanceCache retrieves cached total saldo balance for the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year to retrieve total saldo data for.
	//
	// Returns:
	//   - []*response.SaldoYearTotalBalanceResponse: The yearly total saldo data.
	//   - bool: Whether the cache was found and valid.
	GetYearTotalSaldoBalanceCache(ctx context.Context, year int) ([]*response.SaldoYearTotalBalanceResponse, bool)

	// SetYearTotalSaldoBalanceCache stores total saldo balance for a specific year in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as cache key.
	//   - data: The data to be cached.
	SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data []*response.SaldoYearTotalBalanceResponse)
}

type SaldoStatsBalanceCache interface {
	// GetMonthlySaldoBalanceCache retrieves cached saldo balance per month for a specific year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year to retrieve monthly saldo balances for.
	//
	// Returns:
	//   - []*response.SaldoMonthBalanceResponse: The list of monthly saldo balances.
	//   - bool: Whether the cache was found and valid.
	GetMonthlySaldoBalanceCache(ctx context.Context, year int) ([]*response.SaldoMonthBalanceResponse, bool)

	// SetMonthlySaldoBalanceCache stores saldo balance per month in cache for a specific year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as cache key.
	//   - data: The data to be cached.
	SetMonthlySaldoBalanceCache(ctx context.Context, year int, data []*response.SaldoMonthBalanceResponse)

	// GetYearlySaldoBalanceCache retrieves cached yearly saldo balances for the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year to retrieve saldo balances for.
	//
	// Returns:
	//   - []*response.SaldoYearBalanceResponse: The list of yearly saldo balances.
	//   - bool: Whether the cache was found and valid.
	GetYearlySaldoBalanceCache(ctx context.Context, year int) ([]*response.SaldoYearBalanceResponse, bool)

	// SetYearlySaldoBalanceCache stores saldo balances per year in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as cache key.
	//   - data: The data to be cached.
	SetYearlySaldoBalanceCache(ctx context.Context, year int, data []*response.SaldoYearBalanceResponse)
}

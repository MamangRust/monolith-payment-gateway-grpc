package withdrawstatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type WithdrawStatsStatusCache interface {
	// GetCachedMonthWithdrawStatusSuccessCache retrieves cached monthly statistics of successful withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for month and status.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusSuccess: List of monthly successful withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, bool)

	// SetCachedMonthWithdrawStatusSuccessCache stores monthly successful withdraw statistics in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters used for caching.
	//   - data: The successful withdraw statistics to cache.
	SetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusSuccess)

	// GetCachedYearlyWithdrawStatusSuccessCache retrieves cached yearly statistics of successful withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusSuccess: List of yearly successful withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusSuccess, bool)

	// SetCachedYearlyWithdrawStatusSuccessCache stores yearly successful withdraw statistics in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the statistics.
	//   - data: The successful withdraw statistics to cache.
	SetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int, data []*response.WithdrawResponseYearStatusSuccess)

	// GetCachedMonthWithdrawStatusFailedCache retrieves cached monthly statistics of failed withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for month and status.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusFailed: List of monthly failed withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, bool)

	// SetCachedMonthWithdrawStatusFailedCache stores monthly failed withdraw statistics in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters used for caching.
	//   - data: The failed withdraw statistics to cache.
	SetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusFailed)

	// GetCachedYearlyWithdrawStatusFailedCache retrieves cached yearly statistics of failed withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusFailed: List of yearly failed withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusFailed, bool)

	// SetCachedYearlyWithdrawStatusFailedCache stores yearly failed withdraw statistics in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the statistics.
	//   - data: The failed withdraw statistics to cache.
	SetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int, data []*response.WithdrawResponseYearStatusFailed)
}

type WithdrawStatsAmountCache interface {
	// GetCachedMonthlyWithdraws retrieves cached monthly withdraw amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly data is requested.
	//
	// Returns:
	//   - []*response.WithdrawMonthlyAmountResponse: List of monthly withdraw amounts.
	//   - bool: Whether the cache was found.
	GetCachedMonthlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawMonthlyAmountResponse, bool)

	// SetCachedMonthlyWithdraws stores monthly withdraw amounts in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the monthly data.
	//   - data: The monthly withdraw amounts to cache.
	SetCachedMonthlyWithdraws(ctx context.Context, year int, data []*response.WithdrawMonthlyAmountResponse)

	// GetCachedYearlyWithdraws retrieves cached yearly withdraw amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.WithdrawYearlyAmountResponse: List of yearly withdraw amounts.
	//   - bool: Whether the cache was found.
	GetCachedYearlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawYearlyAmountResponse, bool)

	// SetCachedYearlyWithdraws stores yearly withdraw amounts in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the statistics.
	//   - data: The yearly withdraw amounts to cache.
	SetCachedYearlyWithdraws(ctx context.Context, year int, data []*response.WithdrawYearlyAmountResponse)
}

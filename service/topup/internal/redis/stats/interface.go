package topupstatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TopupStatsStatusCache interface {
	// GetMonthTopupStatusSuccessCache retrieves cached monthly topup statistics with status "success".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and optional month filter.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusSuccess: List of monthly successful topup responses.
	//   - bool: Whether the cache was found.
	GetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, bool)

	// SetMonthTopupStatusSuccessCache stores the monthly successful topup statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The original request used as the cache key.
	//   - data: The data to be cached.
	SetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusSuccess)

	// GetYearlyTopupStatusSuccessCache retrieves cached yearly topup statistics with status "success".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the statistics.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusSuccess: List of yearly successful topup responses.
	//   - bool: Whether the cache was found.
	GetYearlyTopupStatusSuccessCache(ctx context.Context, year int) ([]*response.TopupResponseYearStatusSuccess, bool)

	// SetYearlyTopupStatusSuccessCache stores yearly successful topup statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the data.
	//   - data: The data to cache.
	SetYearlyTopupStatusSuccessCache(ctx context.Context, year int, data []*response.TopupResponseYearStatusSuccess)

	// GetMonthTopupStatusFailedCache retrieves cached monthly topup statistics with status "failed".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and optional month filter.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed topup responses.
	//   - bool: Whether the cache was found.
	GetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, bool)

	// SetMonthTopupStatusFailedCache stores monthly failed topup statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The original request used as the cache key.
	//   - data: The data to be cached.
	SetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusFailed)

	// GetYearlyTopupStatusFailedCache retrieves cached yearly topup statistics with status "failed".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the statistics.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusFailed: List of yearly failed topup responses.
	//   - bool: Whether the cache was found.
	GetYearlyTopupStatusFailedCache(ctx context.Context, year int) ([]*response.TopupResponseYearStatusFailed, bool)

	// SetYearlyTopupStatusFailedCache stores yearly failed topup statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the data.
	//   - data: The data to cache.
	SetYearlyTopupStatusFailedCache(ctx context.Context, year int, data []*response.TopupResponseYearStatusFailed)
}

type TopupStatsMethodCache interface {
	// GetMonthlyTopupMethodsCache retrieves cached monthly statistics grouped by topup method.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TopupMonthMethodResponse: List of method-based monthly topup responses.
	//   - bool: Whether the cache was found.
	GetMonthlyTopupMethodsCache(ctx context.Context, year int) ([]*response.TopupMonthMethodResponse, bool)

	// SetMonthlyTopupMethodsCache stores method-based monthly topup statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the data.
	//   - data: The data to cache.
	SetMonthlyTopupMethodsCache(ctx context.Context, year int, data []*response.TopupMonthMethodResponse)

	// GetYearlyTopupMethodsCache retrieves cached yearly statistics grouped by topup method.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TopupYearlyMethodResponse: List of method-based yearly topup responses.
	//   - bool: Whether the cache was found.
	GetYearlyTopupMethodsCache(ctx context.Context, year int) ([]*response.TopupYearlyMethodResponse, bool)

	// SetYearlyTopupMethodsCache stores method-based yearly topup statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the data.
	//   - data: The data to cache.
	SetYearlyTopupMethodsCache(ctx context.Context, year int, data []*response.TopupYearlyMethodResponse)
}

type TopupStatsAmountCache interface {

	// GetMonthlyTopupAmountsCache retrieves cached monthly total amount statistics of topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TopupMonthAmountResponse: List of monthly topup amount statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyTopupAmountsCache(ctx context.Context, year int) ([]*response.TopupMonthAmountResponse, bool)

	// SetMonthlyTopupAmountsCache stores monthly topup amount statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the data.
	//   - data: The data to cache.
	SetMonthlyTopupAmountsCache(ctx context.Context, year int, data []*response.TopupMonthAmountResponse)

	// GetYearlyTopupAmountsCache retrieves cached yearly total amount statistics of topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TopupYearlyAmountResponse: List of yearly topup amount statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTopupAmountsCache(ctx context.Context, year int) ([]*response.TopupYearlyAmountResponse, bool)

	// SetYearlyTopupAmountsCache stores yearly topup amount statistics in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year of the data.
	//   - data: The data to cache.
	SetYearlyTopupAmountsCache(ctx context.Context, year int, data []*response.TopupYearlyAmountResponse)
}

package transferstatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransferStatsAmountCache interface {
	// GetCachedMonthTransferAmounts retrieves cached monthly total transfer amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly amounts are requested.
	//
	// Returns:
	//   - []*response.TransferMonthAmountResponse: List of monthly transfer amount statistics.
	//   - bool: Whether the cache was found.
	GetCachedMonthTransferAmounts(ctx context.Context, year int) ([]*response.TransferMonthAmountResponse, bool)

	// SetCachedMonthTransferAmounts stores monthly transfer amounts into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is cached.
	//   - data: List of monthly transfer amount statistics to cache.
	SetCachedMonthTransferAmounts(ctx context.Context, year int, data []*response.TransferMonthAmountResponse)

	// GetCachedYearlyTransferAmounts retrieves cached yearly total transfer amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransferYearAmountResponse: List of yearly transfer amount statistics.
	//   - bool: Whether the cache was found.
	GetCachedYearlyTransferAmounts(ctx context.Context, year int) ([]*response.TransferYearAmountResponse, bool)

	// SetCachedYearlyTransferAmounts stores yearly transfer amounts into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is cached.
	//   - data: List of yearly transfer amount statistics to cache.
	SetCachedYearlyTransferAmounts(ctx context.Context, year int, data []*response.TransferYearAmountResponse)
}

type TransferStatsStatusCache interface {
	// GetCachedMonthTransferStatusSuccess retrieves cached monthly successful transfer status.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and status filter.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusSuccess: List of monthly successful transfer status.
	//   - bool: Whether the cache was found.
	GetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, bool)

	// SetCachedMonthTransferStatusSuccess stores monthly successful transfer status into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request key used for caching.
	//   - data: List of monthly successful transfer status to cache.
	SetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusSuccess)

	// GetCachedYearlyTransferStatusSuccess retrieves cached yearly successful transfer status.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which statistics are requested.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusSuccess: List of yearly successful transfer status.
	//   - bool: Whether the cache was found.
	GetCachedYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*response.TransferResponseYearStatusSuccess, bool)

	// SetCachedYearlyTransferStatusSuccess stores yearly successful transfer status into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is cached.
	//   - data: List of yearly successful transfer status to cache.
	SetCachedYearlyTransferStatusSuccess(ctx context.Context, year int, data []*response.TransferResponseYearStatusSuccess)

	// GetCachedMonthTransferStatusFailed retrieves cached monthly failed transfer status.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and status filter.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusFailed: List of monthly failed transfer status.
	//   - bool: Whether the cache was found.
	GetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, bool)

	// SetCachedMonthTransferStatusFailed stores monthly failed transfer status into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request key used for caching.
	//   - data: List of monthly failed transfer status to cache.
	SetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusFailed)

	// GetCachedYearlyTransferStatusFailed retrieves cached yearly failed transfer status.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which statistics are requested.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusFailed: List of yearly failed transfer status.
	//   - bool: Whether the cache was found.
	GetCachedYearlyTransferStatusFailed(ctx context.Context, year int) ([]*response.TransferResponseYearStatusFailed, bool)

	// SetCachedYearlyTransferStatusFailed stores yearly failed transfer status into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is cached.
	//   - data: List of yearly failed transfer status to cache.
	SetCachedYearlyTransferStatusFailed(ctx context.Context, year int, data []*response.TransferResponseYearStatusFailed)
}

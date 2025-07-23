package transferstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transferStatsStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsStatusCache(store *sharedcachehelpers.CacheStore) TransferStatsStatusCache {
	return &transferStatsStatusCache{store: store}
}

// GetCachedMonthTransferStatusSuccess retrieves cached monthly successful transfer status.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and status filter.
//
// Returns:
//   - []*response.TransferResponseMonthStatusSuccess: List of monthly successful transfer status.
//   - bool: Whether the cache was found.
func (t *transferStatsStatusCache) GetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseMonthStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthTransferStatusSuccess stores monthly successful transfer status into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request key used for caching.
//   - data: List of monthly successful transfer status to cache.
func (t *transferStatsStatusCache) SetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetCachedYearlyTransferStatusSuccess retrieves cached yearly successful transfer status.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which statistics are requested.
//
// Returns:
//   - []*response.TransferResponseYearStatusSuccess: List of yearly successful transfer status.
//   - bool: Whether the cache was found.
func (t *transferStatsStatusCache) GetCachedYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*response.TransferResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseYearStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyTransferStatusSuccess stores yearly successful transfer status into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is cached.
//   - data: List of yearly successful transfer status to cache.
func (t *transferStatsStatusCache) SetCachedYearlyTransferStatusSuccess(ctx context.Context, year int, data []*response.TransferResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetCachedMonthTransferStatusFailed retrieves cached monthly failed transfer status.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and status filter.
//
// Returns:
//   - []*response.TransferResponseMonthStatusFailed: List of monthly failed transfer status.
//   - bool: Whether the cache was found.
func (t *transferStatsStatusCache) GetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseMonthStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthTransferStatusFailed stores monthly failed transfer status into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request key used for caching.
//   - data: List of monthly failed transfer status to cache.
func (t *transferStatsStatusCache) SetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetCachedYearlyTransferStatusFailed retrieves cached yearly failed transfer status.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which statistics are requested.
//
// Returns:
//   - []*response.TransferResponseYearStatusFailed: List of yearly failed transfer status.
//   - bool: Whether the cache was found.
func (t *transferStatsStatusCache) GetCachedYearlyTransferStatusFailed(ctx context.Context, year int) ([]*response.TransferResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseYearStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyTransferStatusFailed stores yearly failed transfer status into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is cached.
//   - data: List of yearly failed transfer status to cache.
func (t *transferStatsStatusCache) SetCachedYearlyTransferStatusFailed(ctx context.Context, year int, data []*response.TransferResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

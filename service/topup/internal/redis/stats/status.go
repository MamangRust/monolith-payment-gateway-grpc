package topupstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type topupStatsStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsStatusCache(store *sharedcachehelpers.CacheStore) TopupStatsStatusCache {
	return &topupStatsStatusCache{store: store}
}

// GetMonthTopupStatusSuccessCache retrieves cached monthly topup statistics with status "success".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and optional month filter.
//
// Returns:
//   - []*response.TopupResponseMonthStatusSuccess: List of monthly successful topup responses.
//   - bool: Whether the cache was found.
func (c *topupStatsStatusCache) GetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseMonthStatusSuccess](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTopupStatusSuccessCache stores the monthly successful topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The original request used as the cache key.
//   - data: The data to be cached.
func (c *topupStatsStatusCache) SetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetYearlyTopupStatusSuccessCache retrieves cached yearly topup statistics with status "success".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the statistics.
//
// Returns:
//   - []*response.TopupResponseYearStatusSuccess: List of yearly successful topup responses.
//   - bool: Whether the cache was found.
func (c *topupStatsStatusCache) GetYearlyTopupStatusSuccessCache(ctx context.Context, year int) ([]*response.TopupResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseYearStatusSuccess](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupStatusSuccessCache stores yearly successful topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (c *topupStatsStatusCache) SetYearlyTopupStatusSuccessCache(ctx context.Context, year int, data []*response.TopupResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetMonthTopupStatusFailedCache retrieves cached monthly topup statistics with status "failed".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and optional month filter.
//
// Returns:
//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed topup responses.
//   - bool: Whether the cache was found.
func (c *topupStatsStatusCache) GetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseMonthStatusFailed](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTopupStatusFailedCache stores monthly failed topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The original request used as the cache key.
//   - data: The data to be cached.
func (c *topupStatsStatusCache) SetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetYearlyTopupStatusFailedCache retrieves cached yearly topup statistics with status "failed".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the statistics.
//
// Returns:
//   - []*response.TopupResponseYearStatusFailed: List of yearly failed topup responses.
//   - bool: Whether the cache was found.
func (c *topupStatsStatusCache) GetYearlyTopupStatusFailedCache(ctx context.Context, year int) ([]*response.TopupResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseYearStatusFailed](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetYearlyTopupStatusFailedCache stores yearly failed topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (c *topupStatsStatusCache) SetYearlyTopupStatusFailedCache(ctx context.Context, year int, data []*response.TopupResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

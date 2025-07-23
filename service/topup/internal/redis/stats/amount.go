package topupstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type topupStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsAmountCache(store *sharedcachehelpers.CacheStore) TopupStatsAmountCache {
	return &topupStatsAmountCache{store: store}
}

// GetMonthlyTopupAmountsCache retrieves cached monthly total amount statistics of topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TopupMonthAmountResponse: List of monthly topup amount statistics.
//   - bool: Whether the cache was found.
func (c *topupStatsAmountCache) GetMonthlyTopupAmountsCache(ctx context.Context, year int) ([]*response.TopupMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupMonthAmountResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTopupAmountsCache stores monthly topup amount statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (c *topupStatsAmountCache) SetMonthlyTopupAmountsCache(ctx context.Context, year int, data []*response.TopupMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetYearlyTopupAmountsCache retrieves cached yearly total amount statistics of topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TopupYearlyAmountResponse: List of yearly topup amount statistics.
//   - bool: Whether the cache was found.
func (c *topupStatsAmountCache) GetYearlyTopupAmountsCache(ctx context.Context, year int) ([]*response.TopupYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupYearlyAmountResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupAmountsCache stores yearly topup amount statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (c *topupStatsAmountCache) SetYearlyTopupAmountsCache(ctx context.Context, year int, data []*response.TopupYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

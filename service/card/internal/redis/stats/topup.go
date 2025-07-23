package cardstatsmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsTopupCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTopupCache(store *sharedcachehelpers.CacheStore) CardStatsTopupCache {
	return &cardStatsTopupCache{store: store}
}

// GetMonthlyTopupCache retrieves the global monthly top-up statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly top-up data is requested.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly top-up statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTopupCache) GetMonthlyTopupCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTopupCache stores the global monthly top-up statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTopupCache) SetMonthlyTopupCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetYearlyTopupCache retrieves the global yearly top-up statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly top-up data is requested.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly top-up statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTopupCache) GetYearlyTopupCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupCache stores the global yearly top-up statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTopupCache) SetYearlyTopupCache(ctx context.Context, year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

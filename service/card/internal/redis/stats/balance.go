package cardstatsmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsBalanceCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewCardStatsBalanceCache creates a new instance of cardDashboardByCardNumberCache.
//
// Parameters:
//   - store: the underlying cache store to use.
//
// Returns:
//   - A pointer to a newly created instance of cardStatsBalanceCache.
func NewCardStatsBalanceCache(store *sharedcachehelpers.CacheStore) CardStatsBalanceCache {
	return &cardStatsBalanceCache{store: store}
}

// GetMonthlyBalanceCache retrieves the global monthly balance statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly balance data is requested.
//
// Returns:
//   - []*response.CardResponseMonthBalance: Slice of monthly balance statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsBalanceCache) GetMonthlyBalanceCache(ctx context.Context, year int) ([]*response.CardResponseMonthBalance, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthBalance](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyBalanceCache stores the global monthly balance statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsBalanceCache) SetMonthlyBalanceCache(ctx context.Context, year int, data []*response.CardResponseMonthBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetYearlyBalanceCache retrieves the global yearly balance statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly balance data is requested.
//
// Returns:
//   - []*response.CardResponseYearlyBalance: Slice of yearly balance statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsBalanceCache) GetYearlyBalanceCache(ctx context.Context, year int) ([]*response.CardResponseYearlyBalance, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalance, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearlyBalance](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyBalanceCache stores the global yearly balance statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsBalanceCache) SetYearlyBalanceCache(ctx context.Context, year int, data []*response.CardResponseYearlyBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalance, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

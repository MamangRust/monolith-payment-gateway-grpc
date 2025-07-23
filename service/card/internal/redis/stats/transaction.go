package cardstatsmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsTransactionCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransactionCache(store *sharedcachehelpers.CacheStore) CardStatsTransactionCache {
	return &cardStatsTransactionCache{store: store}
}

// GetMonthlyTransactionCache retrieves the global monthly transaction statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly transaction data is requested.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly transaction statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransactionCache) GetMonthlyTransactionCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransactionCache stores the global monthly transaction statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTransactionCache) SetMonthlyTransactionCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetYearlyTransactionCache retrieves the global yearly transaction statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly transaction data is requested.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly transaction statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransactionCache) GetYearlyTransactionCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransactionCache stores the global yearly transaction statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTransactionCache) SetYearlyTransactionCache(ctx context.Context, year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

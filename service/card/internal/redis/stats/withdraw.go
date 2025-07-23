package cardstatsmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsWithdrawCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsWithdrawCache(store *sharedcachehelpers.CacheStore) CardStatsWithdrawCache {
	return &cardStatsWithdrawCache{store: store}
}

// GetMonthlyWithdrawCache retrieves the global monthly withdraw statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly withdraw data is requested.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly withdraw statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsWithdrawCache) GetMonthlyWithdrawCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyWithdrawCache stores the global monthly withdraw statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsWithdrawCache) SetMonthlyWithdrawCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetYearlyWithdrawCache retrieves the global yearly withdraw statistics
// (across all cards) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly withdraw data is requested.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly withdraw statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsWithdrawCache) GetYearlyWithdrawCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyWithdrawCache stores the global yearly withdraw statistics
// (across all cards) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsWithdrawCache) SetYearlyWithdrawCache(ctx context.Context, year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

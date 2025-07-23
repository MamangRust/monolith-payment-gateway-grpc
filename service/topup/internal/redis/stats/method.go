package topupstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type topupStatsMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsMethodCache(store *sharedcachehelpers.CacheStore) TopupStatsMethodCache {
	return &topupStatsMethodCache{store: store}
}

// GetMonthlyTopupMethodsCache retrieves cached monthly statistics grouped by topup method.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TopupMonthMethodResponse: List of method-based monthly topup responses.
//   - bool: Whether the cache was found.
func (c *topupStatsMethodCache) GetMonthlyTopupMethodsCache(ctx context.Context, month int) ([]*response.TopupMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, month)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupMonthMethodResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetMonthlyTopupMethodsCache stores method-based monthly topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (c *topupStatsMethodCache) SetMonthlyTopupMethodsCache(ctx context.Context, month int, data []*response.TopupMonthMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodCacheKey, month)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetYearlyTopupMethodsCache retrieves cached yearly statistics grouped by topup method.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TopupYearlyMethodResponse: List of method-based yearly topup responses.
//   - bool: Whether the cache was found.
func (c *topupStatsMethodCache) GetYearlyTopupMethodsCache(ctx context.Context, year int) ([]*response.TopupYearlyMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupYearlyMethodResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetYearlyTopupMethodsCache stores method-based yearly topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (c *topupStatsMethodCache) SetYearlyTopupMethodsCache(ctx context.Context, year int, data []*response.TopupYearlyMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

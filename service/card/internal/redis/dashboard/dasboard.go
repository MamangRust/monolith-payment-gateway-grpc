package carddashboardmencache

import (
	"context"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// cardDashboardCache is a struct that represents the cache store
type cardDashboardCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewCardDashboardCache creates a new cardDashboardCache instance
//
// Parameters:
//   - store: the underlying cache store to use
//
// Returns:
//   - A pointer to a newly created cardDashboardCache instance
func NewCardDashboardCache(store *sharedcachehelpers.CacheStore) CardDashboardTotalCache {
	return &cardDashboardCache{store: store}
}

// GetDashboardCardCache retrieves the cached aggregated dashboard data for all cards.
//
// Parameters:
//   - ctx: the context for the operation
//
// Returns:
//   - A pointer to DashboardCard if found in cache, or false if not present.data was found in the cache.
func (c *cardDashboardCache) GetDashboardCardCache(ctx context.Context) (*response.DashboardCard, bool) {
	result, found := sharedcachehelpers.GetFromCache[*response.DashboardCard](ctx, c.store, cacheKeyDashboardDefault)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}

// SetDashboardCardCache caches the aggregated dashboard data for all cards.
//
// Parameters:
//   - ctx: the context for the operation
//   - data: the aggregated dashboard data to cache
func (c *cardDashboardCache) SetDashboardCardCache(ctx context.Context, data *response.DashboardCard) {
	if data == nil {
		return
	}

	sharedcachehelpers.SetToCache(ctx, c.store, cacheKeyDashboardDefault, data, ttlDashboardDefault)
}

// DeleteDashboardCardCache removes the aggregated dashboard cache entry.
//
// Parameters:
//   - ctx: the context for the operation
func (c *cardDashboardCache) DeleteDashboardCardCache(ctx context.Context) {
	sharedcachehelpers.DeleteFromCache(ctx, c.store, cacheKeyDashboardDefault)
}

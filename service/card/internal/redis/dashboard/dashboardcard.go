package carddashboardmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardDashboardByCardNumberCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewCardDashboardByCardNumberCache creates a new instance of cardDashboardByCardNumberCache.
//
// Parameters:
//   - store: the underlying cache store to use.
//
// Returns:
//   - A pointer to a newly created instance of CardDashboardByCardNumberCache.
func NewCardDashboardByCardNumberCache(store *sharedcachehelpers.CacheStore) CardDashboardByCardNumberCache {
	return &cardDashboardByCardNumberCache{store: store}
}

// GetDashboardCardCardNumberCache retrieves cached dashboard data for a specific card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - cardNumber: the specific card number to look up
//
// Returns:
//   - A pointer to DashboardCardCardNumber if found, or false if not present.
func (c *cardDashboardByCardNumberCache) GetDashboardCardCardNumberCache(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, bool) {
	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	result, found := sharedcachehelpers.GetFromCache[*response.DashboardCardCardNumber](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetDashboardCardCardNumberCache caches dashboard data for a specific card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - cardNumber: the card number to associate with the cached data
//   - data: the dashboard data to cache
func (c *cardDashboardByCardNumberCache) SetDashboardCardCardNumberCache(ctx context.Context, cardNumber string, data *response.DashboardCardCardNumber) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDashboardDefault)
}

// DeleteDashboardCardCardNumberCache removes cached dashboard data for a specific card number.
//
// Parameters:
//   - ctx: the context for the operation
//   - cardNumber: the card number whose cache should be cleared
func (c *cardDashboardByCardNumberCache) DeleteDashboardCardCardNumberCache(ctx context.Context, cardNumber string) {
	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

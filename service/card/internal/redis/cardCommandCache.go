package mencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

// cardCommandCache is a struct that represents the cache store
type cardCommandCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewCardCommandCache creates a new cardCommandCache instance
func NewCardCommandCache(store *sharedcachehelpers.CacheStore) CardCommandCache {
	return &cardCommandCache{store: store}
}

// DeleteCardCommandCache removes the cache entry associated with the specified card ID.
// It formats the cache key using the card ID and deletes the entry from the cache store.
func (c *cardCommandCache) DeleteCardCommandCache(ctx context.Context, id int) {
	key := fmt.Sprintf(cardByIdCacheKey, id)

	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

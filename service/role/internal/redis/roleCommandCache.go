package mencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

// roleCommandCache is a struct that represents the cache store
type roleCommandCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewRoleCommandCache creates a new roleCommandCache instance
func NewRoleCommandCache(store *sharedcachehelpers.CacheStore) *roleCommandCache {
	return &roleCommandCache{store: store}
}

// DeleteCachedRole removes a cached role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the role to delete from cache.
func (s *roleCommandCache) DeleteCachedRole(ctx context.Context, id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	sharedcachehelpers.DeleteFromCache(ctx, s.store, key)
}

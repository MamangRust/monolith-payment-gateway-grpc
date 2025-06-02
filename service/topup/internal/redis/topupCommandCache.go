package mencache

import "fmt"

type topupCommandCache struct {
	store *CacheStore
}

func NewTopupCommandCache(store *CacheStore) *topupCommandCache {
	return &topupCommandCache{store: store}
}

func (c *topupCommandCache) DeleteCachedTopupCache(id int) {
	key := fmt.Sprintf(topupByIdCacheKey, id)
	DeleteFromCache(c.store, key)
}

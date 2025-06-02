package mencache

import "fmt"

type cardCommandCache struct {
	store *CacheStore
}

func NewCardCommandCache(store *CacheStore) *cardCommandCache {
	return &cardCommandCache{store: store}
}

func (c *cardCommandCache) DeleteCardCommandCache(id int) {
	key := fmt.Sprintf(cardByIdCacheKey, id)

	DeleteFromCache(c.store, key)
}

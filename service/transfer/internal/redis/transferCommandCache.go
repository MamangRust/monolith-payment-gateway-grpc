package mencache

import "fmt"

type transferCommandCache struct {
	store *CacheStore
}

func NewTransferCommandCache(store *CacheStore) *transferCommandCache {
	return &transferCommandCache{store: store}
}

func (t *transferCommandCache) DeleteTransferCache(id int) {
	DeleteFromCache(t.store, fmt.Sprintf(transferByIdCacheKey, id))
}

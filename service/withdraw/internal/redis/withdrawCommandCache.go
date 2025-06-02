package mencache

import "fmt"

type withdrawCommandCache struct {
	store *CacheStore
}

func NewWithdrawCommandCache(store *CacheStore) *withdrawCommandCache {
	return &withdrawCommandCache{store: store}
}

func (wc *withdrawCommandCache) DeleteCachedWithdrawCache(id int) {
	key := fmt.Sprintf(withdrawByIdCacheKey, id)
	DeleteFromCache(wc.store, key)

}

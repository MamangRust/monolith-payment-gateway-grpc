package mencache

import (
	"fmt"
)

type transactionCommandCache struct {
	store *CacheStore
}

func NewTransactionCommandCache(store *CacheStore) *transactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(id int) {
	DeleteFromCache(t.store, fmt.Sprintf(transactionByIdCacheKey, id))
}

package mencache

import "fmt"

type saldoCommandCache struct {
	store *CacheStore
}

func NewSaldoCommandCache(store *CacheStore) *saldoCommandCache {
	return &saldoCommandCache{store: store}
}

func (s *saldoCommandCache) DeleteSaldoCache(saldo_id int) {
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	DeleteFromCache(s.store, key)
}

package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardDashboardCache struct {
	store *CacheStore
}

const (
	cacheKeyDashboardDefault    = "dashboard:card"
	cacheKeyDashboardCardNumber = "dashboard:card:number:%s"
	ttlDashboardDefault         = 5 * time.Minute
)

func NewCardDashboardCache(store *CacheStore) *cardDashboardCache {
	return &cardDashboardCache{store: store}
}

func (c *cardDashboardCache) GetDashboardCardCache() (*response.DashboardCard, bool) {
	result, found := GetFromCache[*response.DashboardCard](c.store, cacheKeyDashboardDefault)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardDashboardCache) SetDashboardCardCache(data *response.DashboardCard) {
	if data == nil {
		return
	}

	SetToCache(c.store, cacheKeyDashboardDefault, data, ttlDashboardDefault)
}

func (c *cardDashboardCache) DeleteDashboardCardCache() {
	DeleteFromCache(c.store, cacheKeyDashboardDefault)
}

func (c *cardDashboardCache) GetDashboardCardCardNumberCache(cardNumber string) (*response.DashboardCardCardNumber, bool) {
	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	result, found := GetFromCache[*response.DashboardCardCardNumber](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardDashboardCache) SetDashboardCardCardNumberCache(cardNumber string, data *response.DashboardCardCardNumber) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	SetToCache(c.store, key, data, ttlDashboardDefault)
}

func (c *cardDashboardCache) DeleteDashboardCardCardNumberCache(cardNumber string) {
	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	DeleteFromCache(c.store, key)
}

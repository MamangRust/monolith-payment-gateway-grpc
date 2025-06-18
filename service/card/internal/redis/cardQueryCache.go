package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	ttlDefault = 5 * time.Minute

	cardAllCacheKey       = "card:all:page:%d:pageSize:%d:search:%s"
	cardByIdCacheKey      = "card:id:%d"
	cardActiveCacheKey    = "card:active:page:%d:pageSize:%d:search:%s"
	cardTrashedCacheKey   = "card:trashed:page:%d:pageSize:%d:search:%s"
	cardByUserIdCacheKey  = "card:user_id:%d"
	cardByCardNumCacheKey = "card:card_number:%s"
)

type cardCachedResponse struct {
	Data         []*response.CardResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type cardCachedResponseDeleteAt struct {
	Data         []*response.CardResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type cardQueryCache struct {
	store *CacheStore
}

func NewCardQueryCache(store *CacheStore) *cardQueryCache {
	return &cardQueryCache{store: store}
}

func (c *cardQueryCache) GetByIdCache(cardID int) (*response.CardResponse, bool) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)

	result, found := GetFromCache[*response.CardResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetByUserIDCache(userID int) (*response.CardResponse, bool) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)

	result, found := GetFromCache[*response.CardResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetByCardNumberCache(cardNumber string) (*response.CardResponse, bool) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	return GetFromCache[response.CardResponse](c.store, key)
}

func (c *cardQueryCache) GetFindAllCache(req *requests.FindAllCards) ([]*response.CardResponse, *int, bool) {
	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cardCachedResponse](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) GetByActiveCache(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cardCachedResponseDeleteAt](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) GetByTrashedCache(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cardCachedResponseDeleteAt](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) SetByIdCache(cardID int, data *response.CardResponse) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	SetToCache(c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetByUserIDCache(userID int, data *response.CardResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	SetToCache(c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetByCardNumberCache(cardNumber string, data *response.CardResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	SetToCache(c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetFindAllCache(req *requests.FindAllCards, data []*response.CardResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CardResponse{}
	}

	payload := &cardCachedResponse{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetByActiveCache(req *requests.FindAllCards, data []*response.CardResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CardResponseDeleteAt{}
	}

	payload := &cardCachedResponseDeleteAt{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetByTrashedCache(req *requests.FindAllCards, data []*response.CardResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CardResponseDeleteAt{}
	}

	payload := &cardCachedResponseDeleteAt{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) DeleteByIdCache(cardID int) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	DeleteFromCache(c.store, key)
}

func (c *cardQueryCache) DeleteByUserIDCache(userID int) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	DeleteFromCache(c.store, key)
}

func (c *cardQueryCache) DeleteByCardNumberCache(cardNumber string) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	DeleteFromCache(c.store, key)
}

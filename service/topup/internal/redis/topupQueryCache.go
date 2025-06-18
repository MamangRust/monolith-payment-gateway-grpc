package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	topupAllCacheKey     = "topup:all:page:%d:pageSize:%d:search:%s"
	topupByCardCacheKey  = "topup:card_number:%s:page:%d:pageSize:%d:search:%s"
	topupByIdCacheKey    = "topup:id:%d"
	topupActiveCacheKey  = "topup:active:page:%d:pageSize:%d:search:%s"
	topupTrashedCacheKey = "topup:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type topupCachedResponse struct {
	Data  []*response.TopupResponse `json:"data"`
	Total *int                      `json:"total_records"`
}

type topupCachedResponseDeleteAt struct {
	Data  []*response.TopupResponseDeleteAt `json:"data"`
	Total *int                              `json:"total_records"`
}

type topupQueryCache struct {
	store *CacheStore
}

func NewTopupQueryCache(store *CacheStore) *topupQueryCache {
	return &topupQueryCache{store: store}
}

func (c *topupQueryCache) GetCachedTopupsCache(req *requests.FindAllTopups) ([]*response.TopupResponse, *int, bool) {
	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[topupCachedResponse](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupsCache(req *requests.FindAllTopups, data []*response.TopupResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponse{}
	}

	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponse{Data: data, Total: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCacheTopupByCardCache(req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, bool) {
	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[topupCachedResponse](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCacheTopupByCardCache(req *requests.FindAllTopupsByCardNumber, data []*response.TopupResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponse{}
	}

	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponse{Data: data, Total: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupActiveCache(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[topupCachedResponseDeleteAt](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupActiveCache(req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponseDeleteAt{}
	}

	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponseDeleteAt{Data: data, Total: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupTrashedCache(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[topupCachedResponseDeleteAt](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupTrashedCache(req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponseDeleteAt{}
	}

	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponseDeleteAt{Data: data, Total: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupCache(id int) (*response.TopupResponse, bool) {
	key := fmt.Sprintf(topupByIdCacheKey, id)

	result, found := GetFromCache[*response.TopupResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupQueryCache) SetCachedTopupCache(data *response.TopupResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(topupByIdCacheKey, data.ID)
	SetToCache(c.store, key, data, ttlDefault)
}

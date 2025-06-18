package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	transferAllCacheKey     = "transfer:all:page:%d:pageSize:%d:search:%s"
	transferByIdCacheKey    = "transfer:id:%d"
	transferActiveCacheKey  = "transfer:active:page:%d:pageSize:%d:search:%s"
	transferTrashedCacheKey = "transfer:trashed:page:%d:pageSize:%d:search:%s"

	transferByFromCacheKey = "transfer:from_card_number:%s:"
	transferByToCacheKey   = "transfer:to_card_number:%s"

	ttlDefault = 5 * time.Minute
)

type transferCacheResponse struct {
	Data         []*response.TransferResponse `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type transferCachedResponseDeleteAt struct {
	Data         []*response.TransferResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type transferQueryCache struct {
	store *CacheStore
}

func NewTransferQueryCache(store *CacheStore) *transferQueryCache {
	return &transferQueryCache{store: store}
}

func (c *transferQueryCache) GetCachedTransfersCache(req *requests.FindAllTranfers) ([]*response.TransferResponse, *int, bool) {
	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transferCacheResponse](c.store, key)

	if !found {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransfersCache(req *requests.FindAllTranfers, data []*response.TransferResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transferCacheResponse{Data: data, TotalRecords: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferActiveCache(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transferCachedResponseDeleteAt](c.store, key)
	if !found {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}
func (c *transferQueryCache) SetCachedTransferActiveCache(req *requests.FindAllTranfers, data []*response.TransferResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transferCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferTrashedCache(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transferCachedResponseDeleteAt](c.store, key)
	if !found {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransferTrashedCache(req *requests.FindAllTranfers, data []*response.TransferResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transferCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferCache(id int) (*response.TransferResponse, bool) {
	key := fmt.Sprintf(transferByIdCacheKey, id)
	result, found := GetFromCache[*response.TransferResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (c *transferQueryCache) SetCachedTransferCache(data *response.TransferResponse) {
	key := fmt.Sprintf(transferByIdCacheKey, data.ID)
	SetToCache(c.store, key, data, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferByFrom(from string) ([]*response.TransferResponse, bool) {
	key := fmt.Sprintf(transferByFromCacheKey, from)
	result, found := GetFromCache[[]*response.TransferResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}
func (c *transferQueryCache) SetCachedTransferByFrom(from string, data []*response.TransferResponse) {
	key := fmt.Sprintf(transferByFromCacheKey, from)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferByTo(to string) ([]*response.TransferResponse, bool) {
	key := fmt.Sprintf(transferByToCacheKey, to)

	result, found := GetFromCache[[]*response.TransferResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}
func (c *transferQueryCache) SetCachedTransferByTo(to string, data []*response.TransferResponse) {
	key := fmt.Sprintf(transferByToCacheKey, to)
	SetToCache(c.store, key, &data, ttlDefault)
}

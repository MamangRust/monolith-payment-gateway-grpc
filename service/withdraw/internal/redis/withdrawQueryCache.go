package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	withdrawAllCacheKey     = "withdraws:all:page:%d:pageSize:%d:search:%s"
	withdrawByCardCacheKey  = "withdraws:card_number:%s:page:%d:pageSize:%d:search:%s"
	withdrawByIdCacheKey    = "withdraws:id:%d"
	withdrawActiveCacheKey  = "withdraws:active:page:%d:pageSize:%d:search:%s"
	withdrawTrashedCacheKey = "withdraws:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type withdrawCachedResponse struct {
	Data         []*response.WithdrawResponse `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type withdrawCachedResponseDeleteAt struct {
	Data         []*response.WithdrawResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type withdrawQueryCache struct {
	store *CacheStore
}

func NewWithdrawQueryCache(store *CacheStore) *withdrawQueryCache {
	return &withdrawQueryCache{store: store}
}

func (w *withdrawQueryCache) GetCachedWithdrawsCache(req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, bool) {
	key := fmt.Sprintf(withdrawAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[withdrawCachedResponse](w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawsCache(req *requests.FindAllWithdraws, data []*response.WithdrawResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponse{}
	}

	key := fmt.Sprintf(withdrawAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &withdrawCachedResponse{Data: data, TotalRecords: total}

	SetToCache(w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawByCardCache(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, bool) {
	key := fmt.Sprintf(withdrawByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[withdrawCachedResponse](w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawByCardCache(req *requests.FindAllWithdrawCardNumber, data []*response.WithdrawResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponse{}
	}

	key := fmt.Sprintf(withdrawByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	payload := &withdrawCachedResponse{Data: data, TotalRecords: total}
	SetToCache(w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawActiveCache(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(withdrawActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[withdrawCachedResponseDeleteAt](w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawActiveCache(req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponseDeleteAt{}
	}

	key := fmt.Sprintf(withdrawActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawTrashedCache(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(withdrawTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[withdrawCachedResponseDeleteAt](w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawTrashedCache(req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponseDeleteAt{}
	}

	key := fmt.Sprintf(withdrawTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawCache(id int) (*response.WithdrawResponse, bool) {
	key := fmt.Sprintf(withdrawByIdCacheKey, id)
	result, found := GetFromCache[*response.WithdrawResponse](w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}

func (w *withdrawQueryCache) SetCachedWithdrawCache(data *response.WithdrawResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(withdrawByIdCacheKey, data.ID)
	SetToCache(w.store, key, data, ttlDefault)
}

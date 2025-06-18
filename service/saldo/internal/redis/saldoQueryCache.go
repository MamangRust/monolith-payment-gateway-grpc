package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	saldoAllCacheKey     = "saldo:all:page:%d:pageSize:%d:search:%s"
	saldoActiveCacheKey  = "saldo:active:page:%d:pageSize:%d:search:%s"
	saldoTrashedCacheKey = "saldo:trashed:page:%d:pageSize:%d:search:%s"
	saldoByIdCacheKey    = "saldo:id:%d"
	saldoByCardNumberKey = "saldo:card_number:%s"

	ttlDefault = 5 * time.Minute
)

type saldoCachedResponse struct {
	Data         []*response.SaldoResponse `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type saldoCachedResponseDeleteAt struct {
	Data         []*response.SaldoResponseDeleteAt `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type saldoQueryCache struct {
	store *CacheStore
}

func NewSaldoQueryCache(store *CacheStore) *saldoQueryCache {
	return &saldoQueryCache{store: store}
}

func (s *saldoQueryCache) GetCachedSaldos(req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, bool) {
	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[saldoCachedResponse](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) GetCachedSaldoByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[saldoCachedResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) GetCachedSaldoByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[saldoCachedResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) GetCachedSaldoById(saldo_id int) (*response.SaldoResponse, bool) {
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	result, found := GetFromCache[*response.SaldoResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *saldoQueryCache) GetCachedSaldoByCardNumber(card_number string) (*response.SaldoResponse, bool) {
	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	result, found := GetFromCache[*response.SaldoResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *saldoQueryCache) SetCachedSaldos(req *requests.FindAllSaldos, data []*response.SaldoResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SaldoResponse{}
	}

	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponse{Data: data, TotalRecords: total}

	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) SetCachedSaldoByActive(req *requests.FindAllSaldos, result []*response.SaldoResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if result == nil {
		result = []*response.SaldoResponseDeleteAt{}
	}

	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &saldoCachedResponseDeleteAt{Data: result, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)

}

func (s *saldoQueryCache) SetCachedSaldoByTrashed(req *requests.FindAllSaldos, data []*response.SaldoResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SaldoResponseDeleteAt{}
	}

	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &saldoCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) SetCachedSaldoById(saldo_id int, result *response.SaldoResponse) {
	if result == nil {
		result = &response.SaldoResponse{}
	}

	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	SetToCache(s.store, key, result, ttlDefault)
}

func (s *saldoQueryCache) SetCachedSaldoByCardNumber(card_number string, result *response.SaldoResponse) {
	if result == nil {
		result = &response.SaldoResponse{}
	}

	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	SetToCache(s.store, key, result, ttlDefault)
}

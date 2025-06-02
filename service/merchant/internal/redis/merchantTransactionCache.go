package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	merchantTransactionsCacheKey = "merchant:transaction:search:%s:page:%d:pageSize:%d"

	merchantTransactionApikeyCacheKey = "merchant:transaction:apikey:%s:search:%s:page:%d:pageSize:%d"

	merchantTransactionCacheKey = "merchant:transaction:merchant:%d:search:%s:page:%d:pageSize:%d"
)

type merchantTransactionCachheResponse struct {
	Data         []*response.MerchantTransactionResponse `json:"data"`
	TotalRecords *int                                    `json:"total_records"`
}

type merchantTransactionCachhe struct {
	store *CacheStore
}

func NewMerchantTransactionCachhe(store *CacheStore) *merchantTransactionCachhe {
	return &merchantTransactionCachhe{store: store}
}

func (m *merchantTransactionCachhe) SetCacheAllMerchantTransactions(req *requests.FindAllMerchantTransactions, data []*response.MerchantTransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(merchantTransactionsCacheKey, req.Search, req.Page, req.PageSize)

	payload := &merchantTransactionCachheResponse{Data: data, TotalRecords: total}

	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantTransactionCachhe) GetCacheAllMerchantTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, bool) {
	key := fmt.Sprintf(merchantTransactionsCacheKey, req.Search, req.Page, req.PageSize)

	result, found := GetFromCache[merchantTransactionCachheResponse](m.store, key)
	if !found {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantTransactionCachhe) SetCacheMerchantTransactions(req *requests.FindAllMerchantTransactionsById, data []*response.MerchantTransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(merchantTransactionCacheKey, req.MerchantID, req.Search, req.Page, req.PageSize)

	payload := &merchantTransactionCachheResponse{Data: data, TotalRecords: total}

	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantTransactionCachhe) GetCacheMerchantTransactions(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, bool) {
	key := fmt.Sprintf(merchantTransactionCacheKey, req.MerchantID, req.Search, req.Page, req.PageSize)

	result, found := GetFromCache[merchantTransactionCachheResponse](m.store, key)
	if !found {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantTransactionCachhe) SetCacheMerchantTransactionApikey(req *requests.FindAllMerchantTransactionsByApiKey, data []*response.MerchantTransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(merchantTransactionApikeyCacheKey, req.ApiKey, req.Search, req.Page, req.PageSize)

	payload := &merchantTransactionCachheResponse{Data: data, TotalRecords: total}

	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantTransactionCachhe) GetCacheMerchantTransactionApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, bool) {
	key := fmt.Sprintf(merchantTransactionApikeyCacheKey, req.ApiKey, req.Search, req.Page, req.PageSize)

	result, found := GetFromCache[merchantTransactionCachheResponse](m.store, key)
	if !found {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

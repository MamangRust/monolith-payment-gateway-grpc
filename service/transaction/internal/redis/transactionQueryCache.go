package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	transactionAllCacheKey          = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey         = "transaction:id:%d"
	transactionActiveCacheKey       = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey      = "transaction:trashed:page:%d:pageSize:%d:search:%s"
	transactionByCardCacheKey       = "transaction:card_number:%s:page:%d:pageSize:%d:search:%s"
	transactionByMerchantIdCacheKey = "transaction:merchant_id:%d"

	ttlDefault = 5 * time.Minute
)

type transactionCachedResponse struct {
	Data         []*response.TransactionResponse `json:"data"`
	TotalRecords *int                            `json:"total_records"`
}

type transactionCachedResponseDeleteAt struct {
	Data         []*response.TransactionResponseDeleteAt `json:"data"`
	TotalRecords *int                                    `json:"total_records"`
}

type transactionQueryCache struct {
	store *CacheStore
}

func NewTransactionQueryCache(store *CacheStore) *transactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCachedResponse](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(req *requests.FindAllTransactions, data []*response.TransactionResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponse{}
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponse{Data: data, TotalRecords: totalRecords}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByCardNumberCache(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, bool) {
	key := fmt.Sprintf(transactionByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCachedResponse](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionByCardNumberCache(req *requests.FindAllTransactionCardNumber, data []*response.TransactionResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponse{}
	}

	key := fmt.Sprintf(transactionByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponse{Data: data, TotalRecords: totalRecords}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCachedResponseDeleteAt](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponseDeleteAt{}
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCachedResponseDeleteAt](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponseDeleteAt{}
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(transactionId int) (*response.TransactionResponse, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, transactionId)
	result, found := GetFromCache[*response.TransactionResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(transaction *response.TransactionResponse) {
	if transaction == nil {
		return
	}

	key := fmt.Sprintf(transactionByIdCacheKey, transaction.ID)
	SetToCache(t.store, key, transaction, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchantIdCache(merchantId int) ([]*response.TransactionResponse, bool) {
	key := fmt.Sprintf(transactionByMerchantIdCacheKey, merchantId)
	result, found := GetFromCache[[]*response.TransactionResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchantIdCache(id int, transaction []*response.TransactionResponse) {
	if transaction == nil {
		return
	}

	key := fmt.Sprintf(transactionByMerchantIdCacheKey, id)
	SetToCache(t.store, key, &transaction, ttlDefault)
}

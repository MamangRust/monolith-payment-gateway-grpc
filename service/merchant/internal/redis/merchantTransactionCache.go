package mencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// Cache keys for merchant transactions.
const (
	// Key for caching all merchant transactions with pagination and search parameters.
	merchantTransactionsCacheKey = "merchant:transaction:search:%s:page:%d:pageSize:%d"

	// Key for caching merchant transactions by API key.
	merchantTransactionApikeyCacheKey = "merchant:transaction:apikey:%s:search:%s:page:%d:pageSize:%d"

	// Key for caching merchant transactions by merchant ID.
	merchantTransactionCacheKey = "merchant:transaction:merchant:%d:search:%s:page:%d:pageSize:%d"
)

// merchantTransactionCachheResponse is a struct that represents the cached response
type merchantTransactionCachheResponse struct {
	Data         []*response.MerchantTransactionResponse `json:"data"`
	TotalRecords *int                                    `json:"total_records"`
}

// merchantTransactionCachhe is a struct that represents the merchant transaction cache
type merchantTransactionCachhe struct {
	store *sharedcachehelpers.CacheStore
}

// NewMerchantTransactionCachhe creates a new instance of merchantTransactionCachhe
// with the provided sharedcachehelpers.CacheStore. This function initializes the cache structure
// that will be used to store and retrieve merchant transactions.
//
// Parameters:
//
//	store: A pointer to sharedcachehelpers.CacheStore which is used for caching operations.
//
// Returns:
//
//	A pointer to merchantTransactionCachhe initialized with the given store.
func NewMerchantTransactionCachhe(store *sharedcachehelpers.CacheStore) MerchantTransactionCache {
	return &merchantTransactionCachhe{store: store}
}

// SetCacheAllMerchantTransactions stores merchant transaction data in the cache using
// search parameters, page, and page size as part of the cache key. If the provided
// total or data is nil, default values are used.
//
// Parameters:
//   - req: The request object containing the search string, page, and page size.
//   - data: A list of merchant transaction responses to be cached.
//   - total: The total number of records.
//
// Returns:
//   - None
func (m *merchantTransactionCachhe) SetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions, data []*response.MerchantTransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantTransactionResponse{}
	}

	key := fmt.Sprintf(merchantTransactionsCacheKey, req.Search, req.Page, req.PageSize)

	payload := &merchantTransactionCachheResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// GetCacheAllMerchantTransactions retrieves all merchant transactions from cache
// based on the search parameters, page, and page size specified in the request.
// If the cache entry is found and valid, it returns the cached transactions and total records.
// Otherwise, it returns nil, nil, false.
//
// Parameters:
//   - req: The request object containing the search string, page, and page size.
//
// Returns:
//   - []*response.MerchantTransactionResponse: The list of cached merchant transactions.
//   - *int: The total number of records.
//   - bool: Whether the cache entry is found and valid.
func (m *merchantTransactionCachhe) GetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, bool) {
	key := fmt.Sprintf(merchantTransactionsCacheKey, req.Search, req.Page, req.PageSize)

	result, found := sharedcachehelpers.GetFromCache[merchantTransactionCachheResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCacheMerchantTransactions stores merchant transaction data in the cache using
// merchant ID, search parameters, page, and page size as part of the cache key.
// If the provided total or data is nil, default values are used.
//
// Parameters:
//   - req: The request object containing the merchant ID, search string, page, and page size.
//   - data: A list of merchant transaction responses to be cached.
//   - total: The total number of records.
//
// Returns:
//   - None
func (m *merchantTransactionCachhe) SetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById, data []*response.MerchantTransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantTransactionResponse{}
	}

	key := fmt.Sprintf(merchantTransactionCacheKey, req.MerchantID, req.Search, req.Page, req.PageSize)

	payload := &merchantTransactionCachheResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// GetCacheMerchantTransactions retrieves a list of merchant transactions from cache
// based on the merchant ID, search parameters, page, and page size specified in the request.
// If the cache entry is found and valid, it returns the cached transactions and total records.
// Otherwise, it returns nil, nil, false.
//
// Parameters:
//   - req: The request object containing the merchant ID, search string, page, and page size.
//
// Returns:
//   - []*response.MerchantTransactionResponse: The list of cached merchant transactions.
//   - *int: The total number of records.
//   - bool: Whether the cache entry is found and valid.
func (m *merchantTransactionCachhe) GetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, bool) {
	key := fmt.Sprintf(merchantTransactionCacheKey, req.MerchantID, req.Search, req.Page, req.PageSize)

	result, found := sharedcachehelpers.GetFromCache[merchantTransactionCachheResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCacheMerchantTransactionApikey stores merchant transaction data in the cache using
// API key, search parameters, page, and page size as part of the cache key.
// If the provided total or data is nil, default values are used.
//
// Parameters:
//   - req: The request object containing the API key, search string, page, and page size.
//   - data: A list of merchant transaction responses to be cached.
//   - total: The total number of records.
//
// Returns:
//   - None
func (m *merchantTransactionCachhe) SetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey, data []*response.MerchantTransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantTransactionResponse{}
	}

	key := fmt.Sprintf(merchantTransactionApikeyCacheKey, req.ApiKey, req.Search, req.Page, req.PageSize)

	payload := &merchantTransactionCachheResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// GetCacheMerchantTransactionApikey retrieves a list of merchant transactions from cache
// based on the API key, search parameters, page, and page size specified in the request.
// If the cache entry is found and valid, it returns the cached transactions and total records.
// Otherwise, it returns nil, nil, false.
//
// Parameters:
//   - req: The request object containing the API key, search string, page, and page size.
//
// Returns:
//   - []*response.MerchantTransactionResponse: The list of cached merchant transactions.
//   - *int: The total number of records.
//   - bool: Whether the cache entry is found and valid.
func (m *merchantTransactionCachhe) GetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, bool) {
	key := fmt.Sprintf(merchantTransactionApikeyCacheKey, req.ApiKey, req.Search, req.Page, req.PageSize)

	result, found := sharedcachehelpers.GetFromCache[merchantTransactionCachheResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

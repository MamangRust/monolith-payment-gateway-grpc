package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// Cache keys and default time-to-live for merchant-related data.
const (
	// Key for caching all merchants with pagination and search parameters.
	merchantAllCacheKey = "merchant:all:page:%d:pageSize:%d:search:%s"

	// Key for caching a specific merchant by ID.
	merchantByIdCacheKey = "merchant:id:%d"

	// Key for caching active merchants with pagination and search parameters.
	merchantActiveCacheKey = "merchant:active:page:%d:pageSize:%d:search:%s"

	// Key for caching trashed merchants with pagination and search parameters.
	merchantTrashedCacheKey = "merchant:trashed:page:%d:pageSize:%d:search:%s"

	// Key for caching a merchant by API key.
	merchantByApiKeyCacheKey = "merchant:api_key:%s"

	// Key for caching merchants by user ID.
	merchantByUserIdCacheKey = "merchant:user_id:%d"

	// Default time-to-live for cache entries.
	ttlDefault = 5 * time.Minute
)

// merchantCachedResponse is a struct that represents the cached response
type merchantCachedResponse struct {
	Data         []*response.MerchantResponse `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

// merchantCachedResponseDeleteAt is a struct that represents the cached response
type merchantCachedResponseDeleteAt struct {
	Data         []*response.MerchantResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

// merchantQueryCache is a struct that represents the cache store for merchant-related data.
type merchantQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewMerchantQueryCache creates a new instance of merchantQueryCache.
//
// Parameters:
//   - store: The cache store to use for caching.
//
// Returns:
//   - *merchantQueryCache: The newly created merchantQueryCache instance.
func NewMerchantQueryCache(store *sharedcachehelpers.CacheStore) MerchantQueryCache {
	return &merchantQueryCache{store: store}
}

// GetCachedMerchants retrieves a list of merchants from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, nil, false.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//
// Returns:
//   - []*response.MerchantResponse: The list of merchants
//   - *int: The total records
//   - bool: Whether the cache is found and valid
func (m *merchantQueryCache) GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantCachedResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedMerchants sets the cache entry associated with the specified request.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//   - data: The list of merchants to be cached
//   - total: The total records
//
// Returns:
//   - None
func (m *merchantQueryCache) SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*response.MerchantResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantResponse{}
	}

	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponse{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// GetCachedMerchantActive retrieves a list of active merchants from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, nil, false.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//
// Returns:
//   - []*response.MerchantResponseDeleteAt: The list of active merchants
//   - *int: The total records
//   - bool: Whether the cache is found and valid
func (m *merchantQueryCache) GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedMerchantActive sets the cache entry for active merchants based on the specified request.
// If the data or total is nil, default values are used.
// Parameters:
//   - req: The request object containing the page, page size, and search string.
//   - data: The list of active merchants to be cached.
//   - total: The total number of active merchants.
//
// Returns:
//   - None
func (m *merchantQueryCache) SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// GetCachedMerchantTrashed retrieves a list of trashed merchants from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, nil, false.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//
// Returns:
//   - []*response.MerchantResponseDeleteAt: The list of trashed merchants
//   - *int: The total records
//   - bool: Whether the cache is found and valid
func (m *merchantQueryCache) GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedMerchantTrashed sets the cache entry for trashed merchants based on the specified request.
// If the data or total is nil, default values are used.
// Parameters:
//   - req: The request object containing the page, page size, and search string.
//   - data: The list of trashed merchants to be cached.
//   - total: The total number of trashed merchants.
//
// Returns:
//   - None
func (m *merchantQueryCache) SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseDeleteAt{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// GetCachedMerchant retrieves a merchant from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, false.
// Parameters:
//   - id: The id of the merchant to retrieve
//
// Returns:
//   - *response.MerchantResponse: The cached merchant
//   - bool: Whether the cache is found and valid
func (m *merchantQueryCache) GetCachedMerchant(ctx context.Context, id int) (*response.MerchantResponse, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*response.MerchantResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMerchant sets the cache entry associated with the specified merchant ID.
// If the data is nil, the function will return without caching anything.
// Parameters:
//   - data: The merchant response to be cached
//
// Returns:
//   - None
func (m *merchantQueryCache) SetCachedMerchant(ctx context.Context, data *response.MerchantResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByIdCacheKey, data.ID)

	sharedcachehelpers.SetToCache(ctx, m.store, key, data, ttlDefault)
}

// GetCachedMerchantsByUserId retrieves a list of merchants associated with a specific user ID from cache.
// If the cache is found and contains a valid response, it will return the cached response.
// Otherwise, it will return nil, false.
// Parameters:
//   - id: The user ID for which to retrieve the associated merchants.
//
// Returns:
//   - []*response.MerchantResponse: The list of merchants associated with the user ID.
//   - bool: Whether the cache is found and valid.
func (m *merchantQueryCache) GetCachedMerchantsByUserId(ctx context.Context, id int) ([]*response.MerchantResponse, bool) {
	key := fmt.Sprintf(merchantByUserIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMerchantsByUserId sets the cache entry associated with the specified user ID.
// If the data is nil, the function will return without caching anything.
// Parameters:
//   - userId: The user ID for which to cache the associated merchants.
//   - data: The list of merchant responses to be cached.
//
// Returns:
//   - None
func (m *merchantQueryCache) SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*response.MerchantResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetCachedMerchantByApiKey retrieves a merchant from cache using the provided API key.
// If the cache is found and contains a valid response, it returns the cached merchant.
// Otherwise, it returns nil, false.
// Parameters:
//   - apiKey: The API key used to identify the merchant in the cache.
//
// Returns:
//   - *response.MerchantResponse: The cached merchant associated with the API key.
//   - bool: Whether the cache is found and valid.
func (m *merchantQueryCache) GetCachedMerchantByApiKey(ctx context.Context, apiKey string) (*response.MerchantResponse, bool) {
	key := fmt.Sprintf(merchantByApiKeyCacheKey, apiKey)

	result, found := sharedcachehelpers.GetFromCache[*response.MerchantResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMerchantByApiKey sets the cache entry associated with the specified API key.
// If the provided data is nil, the function will return without caching anything.
// Parameters:
//   - apiKey: The API key used to identify the merchant in the cache.
//   - data: The merchant response to be cached.
//
// Returns:
//   - None
func (m *merchantQueryCache) SetCachedMerchantByApiKey(ctx context.Context, apiKey string, data *response.MerchantResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByApiKeyCacheKey, apiKey)

	sharedcachehelpers.SetToCache(ctx, m.store, key, data, ttlDefault)
}

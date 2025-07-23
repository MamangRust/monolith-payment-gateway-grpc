package mencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// merchantDocumentQueryCacheKey is a struct that represents the cache key
const (
	merchantDocumentAllCacheKey     = "merchant_document:all:page:%d:pageSize:%d:search:%s"
	merchantDocumentByIdCacheKey    = "merchant_document:id:%d"
	merchantDocumentActiveCacheKey  = "merchant_document:active:page:%d:pageSize:%d:search:%s"
	merchantDocumentTrashedCacheKey = "merchant_document:trashed:page:%d:pageSize:%d:search:%s"
)

// merchantDocumentQueryCachedResponse is a struct that represents the cached response

type merchantDocumentQueryCachedResponse struct {
	Data         []*response.MerchantDocumentResponse `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

// merchantDocumentQueryCachedResponseDeleteAt is a struct that represents the cached response
type merchantDocumentQueryCachedResponseDeleteAt struct {
	Data         []*response.MerchantDocumentResponseDeleteAt `json:"data"`
	TotalRecords *int                                         `json:"total_records"`
}

// merchantDocumentQueryCache is a struct that represents the cache
type merchantDocumentQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewMerchantDocumentQueryCache is a function that returns a new merchantDocumentQueryCache
func NewMerchantDocumentQueryCache(store *sharedcachehelpers.CacheStore) MerchantDocumentQueryCache {
	return &merchantDocumentQueryCache{store: store}
}

// GetCachedMerchantDocuments retrieves a list of merchant documents from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, nil, false.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//
// Returns:
//   - []*response.MerchantDocumentResponse: The list of merchant documents
//   - *int: The total records
//   - bool: Whether the cache is found and valid
func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantDocumentQueryCachedResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedMerchantDocuments sets the cache entry associated with the specified request.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//   - data: The list of merchant documents to be cached
//   - total: The total records
//
// Returns:
//   - None
func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponse{}
	}

	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

// GetCachedMerchantDocumentsActive retrieves a list of active merchant documents from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, nil, false.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//
// Returns:
//   - []*response.MerchantDocumentResponseDeleteAt: The list of active merchant documents
//   - *int: The total records
//   - bool: Whether the cache is found and valid
func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedMerchantDocumentsActive sets the cache entry associated with the specified request.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//   - data: The list of active merchant documents to be cached
//   - total: The total records
//
// Returns:
//   - None
func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

// GetCachedMerchantDocumentsTrashed retrieves a list of trashed merchant documents from cache.
// If the cache is found and contains a valid response, it will return the cached
// response. Otherwise, it will return nil, nil, false.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//
// Returns:
//   - []*response.MerchantDocumentResponseDeleteAt: The list of trashed merchant documents
//   - *int: The total records
//   - bool: Whether the cache is found and valid
func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedMerchantDocumentsTrashed sets the cache entry associated with the specified request.
// Parameters:
//   - req: The request object containing the page, page size, and search string
//   - data: The list of trashed merchant documents to be cached
//   - total: The total records
//
// Returns:
//   - None
func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

// GetCachedMerchantDocument retrieves a merchant document from the cache by its ID.
// If the cache is found and contains a valid response, it will return the cached
// merchant document. Otherwise, it will return nil, false.
// Parameters:
//   - id: The ID of the merchant document to retrieve
//
// Returns:
//   - *response.MerchantDocumentResponse: The cached merchant document
//   - bool: Whether the cache is found and valid
func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(ctx context.Context, id int) (*response.MerchantDocumentResponse, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*response.MerchantDocumentResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMerchantDocument sets the cache entry for a specific merchant document by its ID.
// If the provided data is nil, the function will return without caching anything.
// Parameters:
//   - id: The ID of the merchant document to cache
//   - data: The merchant document response to be cached
//
// Returns:
//   - None
func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(ctx context.Context, id int, data *response.MerchantDocumentResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	sharedcachehelpers.SetToCache(ctx, s.store, key, data, ttlDefault)
}

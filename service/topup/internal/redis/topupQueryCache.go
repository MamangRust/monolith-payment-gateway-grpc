package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// cache keys
const (
	topupAllCacheKey     = "topup:all:page:%d:pageSize:%d:search:%s"
	topupByCardCacheKey  = "topup:card_number:%s:page:%d:pageSize:%d:search:%s"
	topupByIdCacheKey    = "topup:id:%d"
	topupActiveCacheKey  = "topup:active:page:%d:pageSize:%d:search:%s"
	topupTrashedCacheKey = "topup:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

// topupCachedResponse is a struct that represents the cached response
type topupCachedResponse struct {
	Data  []*response.TopupResponse `json:"data"`
	Total *int                      `json:"total_records"`
}

// topupCachedResponseDeleteAt is a struct that represents the cached response
type topupCachedResponseDeleteAt struct {
	Data  []*response.TopupResponseDeleteAt `json:"data"`
	Total *int                              `json:"total_records"`
}

// topupQueryCache is a struct that represents the topup query cache
type topupQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewTopupQueryCache creates a new instance of topupQueryCache with the provided sharedcachehelpers.CacheStore.
func NewTopupQueryCache(store *sharedcachehelpers.CacheStore) TopupQueryCache {
	return &topupQueryCache{store: store}
}

// GetCachedTopupsCache retrieves cached list of topups based on the given filter request.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter including pagination and search.
//
// Returns:
//   - []*response.TopupResponse: Cached topup responses.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (c *topupQueryCache) GetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponse, *int, bool) {
	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

// SetCachedTopupsCache stores the topup responses and total record count in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The original request used as the cache key.
//   - data: The topup response data to cache.
//   - total: The total number of records.
func (c *topupQueryCache) SetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups, data []*response.TopupResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponse{}
	}

	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponse{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCacheTopupByCardCache retrieves cached topups by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and optional filters.
//
// Returns:
//   - []*response.TopupResponse: Cached topups associated with the card.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (c *topupQueryCache) GetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, bool) {
	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

// SetCacheTopupByCardCache stores the topups associated with a card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request used to generate the cache key.
//   - data: Topup response data to cache.
//   - total: Total number of records.
func (c *topupQueryCache) SetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber, data []*response.TopupResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponse{}
	}

	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponse{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCachedTopupActiveCache retrieves cached list of active topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request used to generate the cache key.
//
// Returns:
//   - []*response.TopupResponseDeleteAt: List of active (non-deleted) topups.
//   - *int: Total records.
//   - bool: Whether the cache was found.
func (c *topupQueryCache) GetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponseDeleteAt](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

// SetCachedTopupActiveCache stores the active topups in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The original request used as the cache key.
//   - data: Topup response data to cache.
//   - total: Total number of records.
func (c *topupQueryCache) SetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponseDeleteAt{}
	}

	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponseDeleteAt{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCachedTopupTrashedCache retrieves cached list of trashed (soft-deleted) topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request used to generate the cache key.
//
// Returns:
//   - []*response.TopupResponseDeleteAt: List of trashed topups.
//   - *int: Total records.
//   - bool: Whether the cache was found.
func (c *topupQueryCache) GetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponseDeleteAt](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

// SetCachedTopupTrashedCache stores the trashed (soft-deleted) topups in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request used to generate the cache key.
//   - data: Topup response data to cache.
//   - total: Total number of records.
func (c *topupQueryCache) SetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TopupResponseDeleteAt{}
	}

	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &topupCachedResponseDeleteAt{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCachedTopupCache retrieves a single topup record from the cache by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The unique topup ID.
//
// Returns:
//   - *response.TopupResponse: The cached topup response.
//   - bool: Whether the cache was found.
func (c *topupQueryCache) GetCachedTopupCache(ctx context.Context, id int) (*response.TopupResponse, bool) {
	key := fmt.Sprintf(topupByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*response.TopupResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedTopupCache stores a single topup response in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - data: The topup response to be cached.
func (c *topupQueryCache) SetCachedTopupCache(ctx context.Context, data *response.TopupResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(topupByIdCacheKey, data.ID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

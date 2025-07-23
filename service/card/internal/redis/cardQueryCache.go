package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// cardQueryCache is a struct that represents the cache store
const (
	ttlDefault = 5 * time.Minute

	cardAllCacheKey       = "card:all:page:%d:pageSize:%d:search:%s"
	cardByIdCacheKey      = "card:id:%d"
	cardActiveCacheKey    = "card:active:page:%d:pageSize:%d:search:%s"
	cardTrashedCacheKey   = "card:trashed:page:%d:pageSize:%d:search:%s"
	cardByUserIdCacheKey  = "card:user_id:%d"
	cardByCardNumCacheKey = "card:card_number:%s"
)

// cardCachedResponse is a struct that represents the cached response
type cardCachedResponse struct {
	Data         []*response.CardResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

// cardCachedResponseDeleteAt is a struct that represents the cached response
type cardCachedResponseDeleteAt struct {
	Data         []*response.CardResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

// cardQueryCache is a struct that represents the cache store
type cardQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewCardQueryCache creates a new cardQueryCache instance
func NewCardQueryCache(store *sharedcachehelpers.CacheStore) CardQueryCache {
	return &cardQueryCache{store: store}
}

// GetByIdCache gets the card data from the cache store by card id.
// It formats the cache key using the card id and retrieves the data from the cache store.
// If the data is not found or the cache is empty, it returns false.
// Otherwise, it returns the data and true.
func (c *cardQueryCache) GetByIdCache(ctx context.Context, cardID int) (*response.CardResponse, bool) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)

	result, found := sharedcachehelpers.GetFromCache[*response.CardResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// GetByUserIDCache retrieves card data from the cache store by user ID.
// It formats the cache key using the user ID and retrieves the data from the cache store.
// If the data is not found or the cache is empty, it returns false.
// Otherwise, it returns the data and true.
func (c *cardQueryCache) GetByUserIDCache(ctx context.Context, userID int) (*response.CardResponse, bool) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)

	result, found := sharedcachehelpers.GetFromCache[*response.CardResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// GetByCardNumberCache retrieves card data from the cache store by card number.
// It formats the cache key using the card number and retrieves the data from the cache store.
// If the data is not found or the cache is empty, it returns false.
// Otherwise, it returns the data and true.
func (c *cardQueryCache) GetByCardNumberCache(ctx context.Context, cardNumber string) (*response.CardResponse, bool) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	return sharedcachehelpers.GetFromCache[response.CardResponse](ctx, c.store, key)
}

// GetFindAllCache retrieves card data from the cache store using the given request parameters.
// It formats the cache key using the request parameters and retrieves the data from the cache store.
// If the data is not found or the cache is empty, it returns false.
// Otherwise, it returns the data and true.
func (c *cardQueryCache) GetFindAllCache(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponse, *int, bool) {
	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[cardCachedResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// GetByActiveCache retrieves card data from the cache store using the given request parameters.
// It formats the cache key using the request parameters and retrieves the data from the cache store.
// If the data is not found or the cache is empty, it returns false.
// Otherwise, it returns the data and true.
func (c *cardQueryCache) GetByActiveCache(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[cardCachedResponseDeleteAt](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// GetByTrashedCache retrieves trashed card data from the cache store using the given request parameters.
// It formats the cache key using the request parameters and retrieves the data from the cache store.
// If the data is not found or the cache is empty, it returns false.
// Otherwise, it returns the data and true.

func (c *cardQueryCache) GetByTrashedCache(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[cardCachedResponseDeleteAt](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetByIdCache sets the cached card data for a specific card ID.
// If the data is nil, nothing will be set.
func (c *cardQueryCache) SetByIdCache(ctx context.Context, cardID int, data *response.CardResponse) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

// SetByUserIDCache sets the cached card data for a specific user ID.
// If the data is nil, nothing will be set.
func (c *cardQueryCache) SetByUserIDCache(ctx context.Context, userID int, data *response.CardResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

// SetByCardNumberCache sets the cached card data for a specific card number.
// If the data is nil, nothing will be set.
func (c *cardQueryCache) SetByCardNumberCache(ctx context.Context, cardNumber string, data *response.CardResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

// SetFindAllCache sets the cached card data for a given request.
// It formats the cache key using the request parameters and sets the data to the cache store.
// If the data is nil, an empty list will be set.
// If the total records is nil, 0 will be set.
func (c *cardQueryCache) SetFindAllCache(ctx context.Context, req *requests.FindAllCards, data []*response.CardResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CardResponse{}
	}

	payload := &cardCachedResponse{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// SetByActiveCache sets the cached card data for a given request.
// It formats the cache key using the request parameters and sets the data to the cache store.
// If the data is nil, an empty list will be set.
// If the total records is nil, 0 will be set.
func (c *cardQueryCache) SetByActiveCache(ctx context.Context, req *requests.FindAllCards, data []*response.CardResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CardResponseDeleteAt{}
	}

	payload := &cardCachedResponseDeleteAt{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// SetByTrashedCache sets the cached card data for a given request.
// It formats the cache key using the request parameters and sets the data to the cache store.
// If the data is nil, an empty list will be set.
// If the total records is nil, 0 will be set.
func (c *cardQueryCache) SetByTrashedCache(ctx context.Context, req *requests.FindAllCards, data []*response.CardResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CardResponseDeleteAt{}
	}

	payload := &cardCachedResponseDeleteAt{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// DeleteByIdCache removes the cache entry associated with the specified card ID.
// It formats the cache key using the card ID and deletes the entry from the cache store.
func (c *cardQueryCache) DeleteByIdCache(ctx context.Context, cardID int) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

// DeleteByUserIDCache removes the cache entry associated with the specified user ID.
// It formats the cache key using the user ID and deletes the entry from the cache store.
func (c *cardQueryCache) DeleteByUserIDCache(ctx context.Context, userID int) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

// DeleteByCardNumberCache removes the cache entry associated with the specified card number.
// It formats the cache key using the card number and deletes the entry from the cache store.
func (c *cardQueryCache) DeleteByCardNumberCache(ctx context.Context, cardNumber string) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

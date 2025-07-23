package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// Constants for cache keys
const (
	withdrawAllCacheKey     = "withdraws:all:page:%d:pageSize:%d:search:%s"
	withdrawByCardCacheKey  = "withdraws:card_number:%s:page:%d:pageSize:%d:search:%s"
	withdrawByIdCacheKey    = "withdraws:id:%d"
	withdrawActiveCacheKey  = "withdraws:active:page:%d:pageSize:%d:search:%s"
	withdrawTrashedCacheKey = "withdraws:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

// withdrawCachedResponse represents the structure of the cached withdraw data.
type withdrawCachedResponse struct {
	Data         []*response.WithdrawResponse `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

// withdrawCachedResponseDeleteAt represents the structure of the cached withdraw data.
type withdrawCachedResponseDeleteAt struct {
	Data         []*response.WithdrawResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

// withdrawQueryCache is a struct that represents the withdraw query cache
type withdrawQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewWithdrawQueryCache creates and returns a new instance of withdrawQueryCache
// using the provided sharedcachehelpers.CacheStore.
//
// Parameters:
//   - store: The cache store to be used for storing and retrieving cached data.
//
// Returns:
//   - *withdrawQueryCache: A pointer to the newly created withdrawQueryCache instance.
func NewWithdrawQueryCache(store *sharedcachehelpers.CacheStore) WithdrawQueryCache {
	return &withdrawQueryCache{store: store}
}

// GetCachedWithdrawsCache retrieves cached list of all withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for finding all withdraws.
//
// Returns:
//   - []*response.WithdrawResponse: List of withdraws.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (w *withdrawQueryCache) GetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, bool) {
	key := fmt.Sprintf(withdrawAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[withdrawCachedResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedWithdrawsCache stores a list of withdraws in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters used for caching.
//   - data: The withdraw response data to cache.
//   - total: Total number of records.
func (w *withdrawQueryCache) SetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws, data []*response.WithdrawResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponse{}
	}

	key := fmt.Sprintf(withdrawAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &withdrawCachedResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

// GetCachedWithdrawByCardCache retrieves cached withdraws for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for finding withdraws by card number.
//
// Returns:
//   - []*response.WithdrawResponse: List of withdraws.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (w *withdrawQueryCache) GetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, bool) {
	key := fmt.Sprintf(withdrawByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[withdrawCachedResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedWithdrawByCardCache stores withdraws for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters used for caching.
//   - data: The withdraw response data to cache.
//   - total: Total number of records.
func (w *withdrawQueryCache) SetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber, data []*response.WithdrawResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponse{}
	}

	key := fmt.Sprintf(withdrawByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	payload := &withdrawCachedResponse{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

// GetCachedWithdrawActiveCache retrieves cached active (non-deleted) withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for finding active withdraws.
//
// Returns:
//   - []*response.WithdrawResponseDeleteAt: List of active withdraws.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (w *withdrawQueryCache) GetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(withdrawActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[withdrawCachedResponseDeleteAt](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedWithdrawActiveCache stores active withdraws in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters used for caching.
//   - data: The active withdraw response data to cache.
//   - total: Total number of records.
func (w *withdrawQueryCache) SetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponseDeleteAt{}
	}

	key := fmt.Sprintf(withdrawActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

// GetCachedWithdrawTrashedCache retrieves cached trashed (soft-deleted) withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for finding trashed withdraws.
//
// Returns:
//   - []*response.WithdrawResponseDeleteAt: List of trashed withdraws.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (w *withdrawQueryCache) GetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(withdrawTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[withdrawCachedResponseDeleteAt](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedWithdrawTrashedCache stores trashed withdraws in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters used for caching.
//   - data: The trashed withdraw response data to cache.
//   - total: Total number of records.
func (w *withdrawQueryCache) SetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.WithdrawResponseDeleteAt{}
	}

	key := fmt.Sprintf(withdrawTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

// GetCachedWithdrawCache retrieves cached withdraw by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the withdraw to retrieve.
//
// Returns:
//   - *response.WithdrawResponse: The withdraw response.
//   - bool: Whether the cache was found.
func (w *withdrawQueryCache) GetCachedWithdrawCache(ctx context.Context, id int) (*response.WithdrawResponse, bool) {
	key := fmt.Sprintf(withdrawByIdCacheKey, id)
	result, found := sharedcachehelpers.GetFromCache[*response.WithdrawResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}

// SetCachedWithdrawCache stores a withdraw record in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - data: The withdraw response to cache.
func (w *withdrawQueryCache) SetCachedWithdrawCache(ctx context.Context, data *response.WithdrawResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(withdrawByIdCacheKey, data.ID)
	sharedcachehelpers.SetToCache(ctx, w.store, key, data, ttlDefault)
}

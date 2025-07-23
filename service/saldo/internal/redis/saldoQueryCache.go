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
	saldoAllCacheKey     = "saldo:all:page:%d:pageSize:%d:search:%s"
	saldoActiveCacheKey  = "saldo:active:page:%d:pageSize:%d:search:%s"
	saldoTrashedCacheKey = "saldo:trashed:page:%d:pageSize:%d:search:%s"
	saldoByIdCacheKey    = "saldo:id:%d"
	saldoByCardNumberKey = "saldo:card_number:%s"

	ttlDefault = 5 * time.Minute
)

// saldoCachedResponse is a struct that represents the cached response
type saldoCachedResponse struct {
	Data         []*response.SaldoResponse `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

// saldoCachedResponseDeleteAt is a struct that represents the cached response
type saldoCachedResponseDeleteAt struct {
	Data         []*response.SaldoResponseDeleteAt `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

// saldoQueryCache is a struct that represents the cache store
type saldoQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewSaldoQueryCache creates a new instance of saldoQueryCache.
//
// Parameters:
//   - store: The cache store to use for caching.
//
// Returns:
//   - *saldoQueryCache: The newly created saldoQueryCache instance.
func NewSaldoQueryCache(store *sharedcachehelpers.CacheStore) SaldoQueryCache {
	return &saldoQueryCache{store: store}
}

// GetCachedSaldos retrieves a list of saldos from the cache based on filter parameters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination and search filters.
//
// Returns:
//   - []*response.SaldoResponse: The list of saldos.
//   - *int: The total number of records.
//   - bool: Whether the cache was found and valid.
func (s *saldoQueryCache) GetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, bool) {
	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[saldoCachedResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// GetCachedSaldoByActive retrieves a list of active (non-deleted) saldos from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing filter parameters.
//
// Returns:
//   - []*response.SaldoResponseDeleteAt: The list of active saldos.
//   - *int: The total number of records.
//   - bool: Whether the cache was found and valid.
func (s *saldoQueryCache) GetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[saldoCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// GetCachedSaldoByTrashed retrieves a list of trashed (soft-deleted) saldos from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing filter parameters.
//
// Returns:
//   - []*response.SaldoResponseDeleteAt: The list of trashed saldos.
//   - *int: The total number of records.
//   - bool: Whether the cache was found and valid.
func (s *saldoQueryCache) GetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[saldoCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// GetCachedSaldoById retrieves a saldo by its ID from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo.
//
// Returns:
//   - *response.SaldoResponse: The cached saldo data.
//   - bool: Whether the cache was found and valid.
func (s *saldoQueryCache) GetCachedSaldoById(ctx context.Context, saldo_id int) (*response.SaldoResponse, bool) {
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	result, found := sharedcachehelpers.GetFromCache[*response.SaldoResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// GetCachedSaldoByCardNumber retrieves a saldo by card number from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number.
//
// Returns:
//   - *response.SaldoResponse: The cached saldo data.
//   - bool: Whether the cache was found and valid.
func (s *saldoQueryCache) GetCachedSaldoByCardNumber(ctx context.Context, card_number string) (*response.SaldoResponse, bool) {
	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	result, found := sharedcachehelpers.GetFromCache[*response.SaldoResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedSaldos stores a list of saldos in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object used as the cache key.
//   - data: The list of saldos to be cached.
//   - totalRecords: The total number of records.
func (s *saldoQueryCache) SetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos, data []*response.SaldoResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SaldoResponse{}
	}

	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

// SetCachedSaldoByActive stores a list of active (non-deleted) saldos in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object used as the cache key.
//   - data: The list of active saldos to be cached.
//   - totalRecords: The total number of records.
func (s *saldoQueryCache) SetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos, result []*response.SaldoResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if result == nil {
		result = []*response.SaldoResponseDeleteAt{}
	}

	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &saldoCachedResponseDeleteAt{Data: result, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)

}

// SetCachedSaldoByTrashed stores a list of trashed (soft-deleted) saldos in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object used as the cache key.
//   - data: The list of trashed saldos to be cached.
//   - totalRecords: The total number of records.
func (s *saldoQueryCache) SetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos, data []*response.SaldoResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SaldoResponseDeleteAt{}
	}

	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &saldoCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

// SetCachedSaldoById stores a saldo by its ID in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo.
//   - data: The saldo data to cache.
func (s *saldoQueryCache) SetCachedSaldoById(ctx context.Context, saldo_id int, result *response.SaldoResponse) {
	if result == nil {
		result = &response.SaldoResponse{}
	}

	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	sharedcachehelpers.SetToCache(ctx, s.store, key, result, ttlDefault)
}

// GetCachedSaldoByCardNumber retrieves a saldo by card number from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number.
//
// Returns:
//   - *response.SaldoResponse: The cached saldo data.
//   - bool: Whether the cache was found and valid.
func (s *saldoQueryCache) SetCachedSaldoByCardNumber(ctx context.Context, card_number string, result *response.SaldoResponse) {
	if result == nil {
		result = &response.SaldoResponse{}
	}

	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	sharedcachehelpers.SetToCache(ctx, s.store, key, result, ttlDefault)
}

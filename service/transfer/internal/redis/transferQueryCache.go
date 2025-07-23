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
	transferAllCacheKey     = "transfer:all:page:%d:pageSize:%d:search:%s"
	transferByIdCacheKey    = "transfer:id:%d"
	transferActiveCacheKey  = "transfer:active:page:%d:pageSize:%d:search:%s"
	transferTrashedCacheKey = "transfer:trashed:page:%d:pageSize:%d:search:%s"

	transferByFromCacheKey = "transfer:from_card_number:%s:"
	transferByToCacheKey   = "transfer:to_card_number:%s"

	ttlDefault = 5 * time.Minute
)

// transferCacheResponse represents the structure of the cached transfer data.
type transferCacheResponse struct {
	Data         []*response.TransferResponse `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

// transferCachedResponseDeleteAt represents the structure of the cached transfer data.
type transferCachedResponseDeleteAt struct {
	Data         []*response.TransferResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

// transferQueryCache represents the cache for transfer queries.
type transferQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewTransferQueryCache creates a new instance of transferQueryCache with the provided sharedcachehelpers.CacheStore.
// It returns a pointer to the newly created transferQueryCache.
func NewTransferQueryCache(store *sharedcachehelpers.CacheStore) TransferQueryCache {
	return &transferQueryCache{store: store}
}

// GetCachedTransfersCache retrieves cached list of transfers.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter for transfer data.
//
// Returns:
//   - []*response.TransferResponse: List of transfers.
//   - *int: Total count of transfers.
//   - bool: Whether the cache was found.
func (c *transferQueryCache) GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponse, *int, bool) {
	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transferCacheResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

// SetCachedTransfersCache stores list of transfers into the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter used to cache the data.
//   - data: List of transfers to cache.
//   - total: Total count of transfers.
func (c *transferQueryCache) SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers, data []*response.TransferResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransferResponse{}
	}

	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transferCacheResponse{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCachedTransferActiveCache retrieves cached list of active (non-trashed) transfers.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter for active transfers.
//
// Returns:
//   - []*response.TransferResponseDeleteAt: List of active transfers.
//   - *int: Total count of active transfers.
//   - bool: Whether the cache was found.
func (c *transferQueryCache) GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transferCachedResponseDeleteAt](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedTransferActiveCache stores list of active (non-trashed) transfers into the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter used to cache the data.
//   - data: List of active transfers to cache.
//   - total: Total count of active transfers.
func (c *transferQueryCache) SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers, data []*response.TransferResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransferResponseDeleteAt{}
	}

	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transferCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCachedTransferTrashedCache retrieves cached list of trashed transfers.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter for trashed transfers.
//
// Returns:
//   - []*response.TransferResponseDeleteAt: List of trashed transfers.
//   - *int: Total count of trashed transfers.
//   - bool: Whether the cache was found.
func (c *transferQueryCache) GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transferCachedResponseDeleteAt](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedTransferTrashedCache stores list of trashed transfers into the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request filter used to cache the data.
//   - data: List of trashed transfers to cache.
//   - total: Total count of trashed transfers.
func (c *transferQueryCache) SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers, data []*response.TransferResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransferResponseDeleteAt{}
	}

	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transferCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

// GetCachedTransferCache retrieves a specific transfer by ID from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the transfer to retrieve.
//
// Returns:
//   - *response.TransferResponse: Transfer response.
//   - bool: Whether the cache was found.
func (c *transferQueryCache) GetCachedTransferCache(ctx context.Context, id int) (*response.TransferResponse, bool) {
	key := fmt.Sprintf(transferByIdCacheKey, id)
	result, found := sharedcachehelpers.GetFromCache[*response.TransferResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedTransferCache stores a specific transfer into the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - data: The transfer data to cache.
func (c *transferQueryCache) SetCachedTransferCache(ctx context.Context, data *response.TransferResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByIdCacheKey, data.ID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

// GetCachedTransferByFrom retrieves cached transfers filtered by source card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - from: The card number from which the transfer was made.
//
// Returns:
//   - []*response.TransferResponse: List of transfers.
//   - bool: Whether the cache was found.
func (c *transferQueryCache) GetCachedTransferByFrom(ctx context.Context, from string) ([]*response.TransferResponse, bool) {
	key := fmt.Sprintf(transferByFromCacheKey, from)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponse](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetCachedTransferByFrom stores cached transfers by source card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - from: The card number from which the transfer was made.
//   - data: List of transfers to cache.
func (c *transferQueryCache) SetCachedTransferByFrom(ctx context.Context, from string, data []*response.TransferResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByFromCacheKey, from)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetCachedTransferByTo retrieves cached transfers filtered by destination card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - to: The card number to which the transfer was made.
//
// Returns:
//   - []*response.TransferResponse: List of transfers.
//   - bool: Whether the cache was found.
func (c *transferQueryCache) GetCachedTransferByTo(ctx context.Context, to string) ([]*response.TransferResponse, bool) {
	key := fmt.Sprintf(transferByToCacheKey, to)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponse](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetCachedTransferByTo stores cached transfers by destination card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - to: The card number to which the transfer was made.
//   - data: List of transfers to cache.
func (c *transferQueryCache) SetCachedTransferByTo(ctx context.Context, to string, data []*response.TransferResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByToCacheKey, to)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

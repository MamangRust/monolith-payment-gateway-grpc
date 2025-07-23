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
	transactionAllCacheKey          = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey         = "transaction:id:%d"
	transactionActiveCacheKey       = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey      = "transaction:trashed:page:%d:pageSize:%d:search:%s"
	transactionByCardCacheKey       = "transaction:card_number:%s:page:%d:pageSize:%d:search:%s"
	transactionByMerchantIdCacheKey = "transaction:merchant_id:%d"

	ttlDefault = 5 * time.Minute
)

// transactionCachedResponseDeleteAt is a struct that represents the cached response
type transactionCachedResponse struct {
	Data         []*response.TransactionResponse `json:"data"`
	TotalRecords *int                            `json:"total_records"`
}

// transactionCachedResponseDeleteAt is a struct that represents the cached response
type transactionCachedResponseDeleteAt struct {
	Data         []*response.TransactionResponseDeleteAt `json:"data"`
	TotalRecords *int                                    `json:"total_records"`
}

// transactionQueryCache is a struct that represents the transaction query cache
type transactionQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewTransactionQueryCache creates a new instance of transactionQueryCache with the provided sharedcachehelpers.CacheStore.
// It initializes the cache structure that will be used to store and retrieve transaction queries.
//
// Parameters:
//   - store: The cache store to use for caching transaction queries.
//
// Returns:
//   - A pointer to the newly created transactionQueryCache instance.
func NewTransactionQueryCache(store *sharedcachehelpers.CacheStore) TransactionQueryCache {
	return &transactionQueryCache{store: store}
}

// GetCachedTransactionsCache retrieves cached paginated transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Filter and pagination parameters.
//
// Returns:
//   - []*response.TransactionResponse: List of transactions.
//   - *int: Total number of matching records.
//   - bool: Whether the cache was found.
func (t *transactionQueryCache) GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transactionCachedResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedTransactionsCache stores paginated transaction results in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Original request used as cache key.
//   - data: List of transaction responses to cache.
//   - total: Total number of records for pagination.
func (t *transactionQueryCache) SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransactions, data []*response.TransactionResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponse{}
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponse{Data: data, TotalRecords: totalRecords}
	sharedcachehelpers.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

// GetCachedTransactionByCardNumberCache retrieves cached transactions for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and filters.
//
// Returns:
//   - []*response.TransactionResponse: Transactions matching the card number.
//   - *int: Total number of records.
//   - bool: Whether the cache was found..
func (t *transactionQueryCache) GetCachedTransactionByCardNumberCache(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, bool) {
	key := fmt.Sprintf(transactionByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transactionCachedResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedTransactionByCardNumberCache caches transactions by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and filters.
//   - data: Transactions to cache.
//   - total: Total number of matching records.
func (t *transactionQueryCache) SetCachedTransactionByCardNumberCache(ctx context.Context, req *requests.FindAllTransactionCardNumber, data []*response.TransactionResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponse{}
	}

	key := fmt.Sprintf(transactionByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponse{Data: data, TotalRecords: totalRecords}
	sharedcachehelpers.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

// GetCachedTransactionActiveCache retrieves cached active transactions (not soft-deleted).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Filter and pagination request.
//
// Returns:
//   - []*response.TransactionResponseDeleteAt: List of active transactions.
//   - *int: Total number of records.
//   - bool: Whether the cache was found.
func (t *transactionQueryCache) GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transactionCachedResponseDeleteAt](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedTransactionActiveCache caches active transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Original filter request.
//   - data: Transactions to cache.
//   - total: Total number of records.
func (t *transactionQueryCache) SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponseDeleteAt{}
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	sharedcachehelpers.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

// GetCachedTransactionTrashedCache retrieves cached trashed (soft-deleted) transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Filter and pagination request.
//
// Returns:
//   - []*response.TransactionResponseDeleteAt: List of trashed transactions.
//   - *int: Total records.
//   - bool: Whether the cache was found.
func (t *transactionQueryCache) GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transactionCachedResponseDeleteAt](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedTransactionTrashedCache stores trashed transactions in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Original request with filters.
//   - data: Trashed transactions to cache.
//   - total: Total number of records.
func (t *transactionQueryCache) SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.TransactionResponseDeleteAt{}
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCachedResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	sharedcachehelpers.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

// GetCachedTransactionCache retrieves a transaction by ID from cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The transaction ID.
//
// Returns:
//   - *response.TransactionResponse: The transaction response.
//   - bool: Whether the cache was found.
func (t *transactionQueryCache) GetCachedTransactionCache(ctx context.Context, transactionId int) (*response.TransactionResponse, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, transactionId)
	result, found := sharedcachehelpers.GetFromCache[*response.TransactionResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedTransactionCache caches a transaction by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - data: The transaction response to cache.
func (t *transactionQueryCache) SetCachedTransactionCache(ctx context.Context, transaction *response.TransactionResponse) {
	if transaction == nil {
		return
	}

	key := fmt.Sprintf(transactionByIdCacheKey, transaction.ID)
	sharedcachehelpers.SetToCache(ctx, t.store, key, transaction, ttlDefault)
}

// GetCachedTransactionByMerchantIdCache retrieves cached transactions for a merchant ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: ID of the merchant.
//
// Returns:
//   - []*response.TransactionResponse: Transactions related to the merchant.
//   - bool: Whether the cache was found.
func (t *transactionQueryCache) GetCachedTransactionByMerchantIdCache(ctx context.Context, merchantId int) ([]*response.TransactionResponse, bool) {
	key := fmt.Sprintf(transactionByMerchantIdCacheKey, merchantId)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedTransactionByMerchantIdCache caches transactions by merchant ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: ID of the merchant.
//   - data: Transactions to cache.
func (t *transactionQueryCache) SetCachedTransactionByMerchantIdCache(ctx context.Context, id int, transaction []*response.TransactionResponse) {
	if transaction == nil {
		return
	}

	key := fmt.Sprintf(transactionByMerchantIdCacheKey, id)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &transaction, ttlDefault)
}

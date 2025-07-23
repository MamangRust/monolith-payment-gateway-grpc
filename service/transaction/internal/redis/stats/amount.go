package transactionstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transactionStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsAmountCache(store *sharedcachehelpers.CacheStore) TransactionStatsAmountCache {
	return &transactionStatsAmountCache{store: store}
}

// GetMonthlyAmountsCache retrieves cached monthly transaction amount statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionMonthAmountResponse: Monthly amount statistics.
//   - bool: Whether the cache was found.
func (t *transactionStatsAmountCache) GetMonthlyAmountsCache(ctx context.Context, year int) ([]*response.TransactionMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionMonthAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyAmountsCache caches monthly transaction amount statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as key.
//   - data: Monthly transaction amount data.
func (t *transactionStatsAmountCache) SetMonthlyAmountsCache(ctx context.Context, year int, data []*response.TransactionMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyAmountsCache retrieves cached yearly transaction amount statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionYearlyAmountResponse: Yearly amount statistics.
//   - bool: Whether the cache was found.
func (t *transactionStatsAmountCache) GetYearlyAmountsCache(ctx context.Context, year int) ([]*response.TransactionYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionYearlyAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true

}

// SetYearlyAmountsCache caches yearly transaction amount statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as cache key.
//   - data: Yearly amount statistics to cache.
func (t *transactionStatsAmountCache) SetYearlyAmountsCache(ctx context.Context, year int, data []*response.TransactionYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

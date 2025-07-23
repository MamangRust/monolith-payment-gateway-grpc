package transactionstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transactionStatsMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsMethodCache(store *sharedcachehelpers.CacheStore) TransactionStatsMethodCache {
	return &transactionStatsMethodCache{store: store}
}

// GetMonthlyPaymentMethodsCache retrieves cached monthly payment method statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionMonthMethodResponse: List of monthly payment method statistics.
//   - bool: Whether the cache was found.
func (t *transactionStatsMethodCache) GetMonthlyPaymentMethodsCache(ctx context.Context, year int) ([]*response.TransactionMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionMonthMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyPaymentMethodsCache caches monthly payment method statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as key.
//   - data: Monthly payment method data.
func (t *transactionStatsMethodCache) SetMonthlyPaymentMethodsCache(ctx context.Context, year int, data []*response.TransactionMonthMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyPaymentMethodsCache retrieves cached yearly payment method statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionYearMethodResponse: Yearly payment method stats.
//   - bool: Whether the cache was found.
func (t *transactionStatsMethodCache) GetYearlyPaymentMethodsCache(ctx context.Context, year int) ([]*response.TransactionYearMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionYearMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyPaymentMethodsCache caches yearly payment method statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as cache key.
//   - data: Yearly method statistics.
func (t *transactionStatsMethodCache) SetYearlyPaymentMethodsCache(ctx context.Context, year int, data []*response.TransactionYearMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

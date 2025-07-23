package transferstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transferStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsAmountCache(store *sharedcachehelpers.CacheStore) TransferStatsAmountCache {
	return &transferStatsAmountCache{store: store}
}

// GetCachedMonthTransferAmounts retrieves cached monthly total transfer amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly amounts are requested.
//
// Returns:
//   - []*response.TransferMonthAmountResponse: List of monthly transfer amount statistics.
//   - bool: Whether the cache was found.
func (t *transferStatsAmountCache) GetCachedMonthTransferAmounts(ctx context.Context, year int) ([]*response.TransferMonthAmountResponse, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferMonthAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthTransferAmounts stores monthly transfer amounts into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is cached.
//   - data: List of monthly transfer amount statistics to cache.
func (t *transferStatsAmountCache) SetCachedMonthTransferAmounts(ctx context.Context, year int, data []*response.TransferMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetCachedYearlyTransferAmounts retrieves cached yearly total transfer amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransferYearAmountResponse: List of yearly transfer amount statistics.
//   - bool: Whether the cache was found.
func (t *transferStatsAmountCache) GetCachedYearlyTransferAmounts(ctx context.Context, year int) ([]*response.TransferYearAmountResponse, bool) {
	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferYearAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyTransferAmounts stores yearly transfer amounts into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is cached.
//   - data: List of yearly transfer amount statistics to cache.
func (t *transferStatsAmountCache) SetCachedYearlyTransferAmounts(ctx context.Context, year int, data []*response.TransferYearAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

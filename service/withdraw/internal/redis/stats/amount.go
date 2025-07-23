package withdrawstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type withdrawStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsAmountCache(store *sharedcachehelpers.CacheStore) WithdrawStatsAmountCache {
	return &withdrawStatsAmountCache{store: store}
}

// GetCachedMonthlyWithdraws retrieves cached monthly withdraw amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly data is requested.
//
// Returns:
//   - []*response.WithdrawMonthlyAmountResponse: List of monthly withdraw amounts.
//   - bool: Whether the cache was found.
func (w *withdrawStatsAmountCache) GetCachedMonthlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawMonthlyAmountResponse, bool) {
	key := fmt.Sprintf(montWithdrawAmountKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawMonthlyAmountResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthlyWithdraws stores monthly withdraw amounts in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the monthly data.
//   - data: The monthly withdraw amounts to cache.
func (w *withdrawStatsAmountCache) SetCachedMonthlyWithdraws(ctx context.Context, year int, data []*response.WithdrawMonthlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawAmountKey, year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedYearlyWithdraws retrieves cached yearly withdraw amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.WithdrawYearlyAmountResponse: List of yearly withdraw amounts.
//   - bool: Whether the cache was found.
func (w *withdrawStatsAmountCache) GetCachedYearlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawYearlyAmountResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyWithdraws stores yearly withdraw amounts in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the statistics.
//   - data: The yearly withdraw amounts to cache.
func (w *withdrawStatsAmountCache) SetCachedYearlyWithdraws(ctx context.Context, year int, data []*response.WithdrawYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

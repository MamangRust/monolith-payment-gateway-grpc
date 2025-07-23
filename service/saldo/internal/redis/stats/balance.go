package saldostatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type saldoStatsBalanceCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewSaldoStatsBalanceCache(store *sharedcachehelpers.CacheStore) SaldoStatsBalanceCache {
	return &saldoStatsBalanceCache{store: store}
}

// GetMonthlySaldoBalanceCache retrieves cached saldo balance per month for a specific year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year to retrieve monthly saldo balances for.
//
// Returns:
//   - []*response.SaldoMonthBalanceResponse: The list of monthly saldo balances.
//   - bool: Whether the cache was found and valid.
func (c *saldoStatsBalanceCache) GetMonthlySaldoBalanceCache(ctx context.Context, year int) ([]*response.SaldoMonthBalanceResponse, bool) {
	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.SaldoMonthBalanceResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlySaldoBalanceCache stores saldo balance per month in cache for a specific year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as cache key.
//   - data: The data to be cached.
func (c *saldoStatsBalanceCache) SetMonthlySaldoBalanceCache(ctx context.Context, year int, data []*response.SaldoMonthBalanceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetYearlySaldoBalanceCache retrieves cached yearly saldo balances for the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year to retrieve saldo balances for.
//
// Returns:
//   - []*response.SaldoYearBalanceResponse: The list of yearly saldo balances.
//   - bool: Whether the cache was found and valid.
func (c *saldoStatsBalanceCache) GetYearlySaldoBalanceCache(ctx context.Context, year int) ([]*response.SaldoYearBalanceResponse, bool) {
	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.SaldoYearBalanceResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlySaldoBalanceCache stores saldo balances per year in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as cache key.
//   - data: The data to be cached.
func (c *saldoStatsBalanceCache) SetYearlySaldoBalanceCache(ctx context.Context, year int, data []*response.SaldoYearBalanceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

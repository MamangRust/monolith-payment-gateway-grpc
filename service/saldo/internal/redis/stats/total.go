package saldostatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type saldoStatsTotalCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewSaldoStatsTotalCache(store *sharedcachehelpers.CacheStore) SaldoStatsTotalCache {
	return &saldoStatsTotalCache{store: store}
}

// GetMonthlyTotalSaldoBalanceCache retrieves cached total saldo balance per month based on request filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the year, month, and additional filters.
//
// Returns:
//   - []*response.SaldoMonthTotalBalanceResponse: The list of monthly total saldo balances.
//   - bool: Whether the cache was found and valid.
func (c *saldoStatsTotalCache) GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, bool) {
	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.SaldoMonthTotalBalanceResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTotalSaldoCache stores total saldo balance per month in cache based on request filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object used as cache key.
//   - data: The data to be cached.
func (c *saldoStatsTotalCache) SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data []*response.SaldoMonthTotalBalanceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

// GetYearTotalSaldoBalanceCache retrieves cached total saldo balance for the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year to retrieve total saldo data for.
//
// Returns:
//   - []*response.SaldoYearTotalBalanceResponse: The yearly total saldo data.
//   - bool: Whether the cache was found and valid.
func (c *saldoStatsTotalCache) GetYearTotalSaldoBalanceCache(ctx context.Context, year int) ([]*response.SaldoYearTotalBalanceResponse, bool) {
	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.SaldoYearTotalBalanceResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearTotalSaldoBalanceCache stores total saldo balance for a specific year in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as cache key.
//   - data: The data to be cached.
func (c *saldoStatsTotalCache) SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data []*response.SaldoYearTotalBalanceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

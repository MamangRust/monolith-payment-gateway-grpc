package withdrawstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type withdrawStatsStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsStatusCache(store *sharedcachehelpers.CacheStore) WithdrawStatsStatusCache {
	return &withdrawStatsStatusCache{store: store}
}

// GetCachedMonthWithdrawStatusSuccessCache retrieves cached monthly statistics of successful withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for month and status.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusSuccess: List of monthly successful withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsStatusCache) GetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseMonthStatusSuccess](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthWithdrawStatusSuccessCache stores monthly successful withdraw statistics in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters used for caching.
//   - data: The successful withdraw statistics to cache.
func (w *withdrawStatsStatusCache) SetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedYearlyWithdrawStatusSuccessCache retrieves cached yearly statistics of successful withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusSuccess: List of yearly successful withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsStatusCache) GetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseYearStatusSuccess](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyWithdrawStatusSuccessCache stores yearly successful withdraw statistics in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the statistics.
//   - data: The successful withdraw statistics to cache.
func (w *withdrawStatsStatusCache) SetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int, data []*response.WithdrawResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedMonthWithdrawStatusFailedCache retrieves cached monthly statistics of failed withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for month and status.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusFailed: List of monthly failed withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsStatusCache) GetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseMonthStatusFailed](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthWithdrawStatusFailedCache stores monthly failed withdraw statistics in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters used for caching.
//   - data: The failed withdraw statistics to cache.
func (w *withdrawStatsStatusCache) SetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedYearlyWithdrawStatusFailedCache retrieves cached yearly statistics of failed withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusFailed: List of yearly failed withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsStatusCache) GetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseYearStatusFailed](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyWithdrawStatusFailedCache stores yearly failed withdraw statistics in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the statistics.
//   - data: The failed withdraw statistics to cache.
func (w *withdrawStatsStatusCache) SetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int, data []*response.WithdrawResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

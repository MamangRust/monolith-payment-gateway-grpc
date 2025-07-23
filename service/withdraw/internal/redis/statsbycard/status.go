package withdrawstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type withdrawStatsByCardStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsStatusCache(store *sharedcachehelpers.CacheStore) WithdrawStatsByCardStatusCache {
	return &withdrawStatsByCardStatusCache{store: store}
}

// GetCachedMonthWithdrawStatusSuccessByCardNumber retrieves cached monthly successful withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and card number.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusSuccess: List of monthly successful withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsByCardStatusCache) GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseMonthStatusSuccess](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthWithdrawStatusSuccessByCardNumber stores monthly successful withdraw statistics in the cache by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and card number.
//   - data: The data to cache.
func (w *withdrawStatsByCardStatusCache) SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedYearlyWithdrawStatusSuccessByCardNumber retrieves cached yearly successful withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusSuccess: List of yearly successful withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsByCardStatusCache) GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseYearStatusSuccess](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyWithdrawStatusSuccessByCardNumber stores yearly successful withdraw statistics in the cache by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//   - data: The data to cache.
func (w *withdrawStatsByCardStatusCache) SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusSuccessByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedMonthWithdrawStatusFailedByCardNumber retrieves cached monthly failed withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and card number.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusFailed: List of monthly failed withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsByCardStatusCache) GetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseMonthStatusFailed](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthWithdrawStatusFailedByCardNumber stores monthly failed withdraw statistics in the cache by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and card number.
//   - data: The data to cache.
func (w *withdrawStatsByCardStatusCache) SetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedYearlyWithdrawStatusFailedByCardNumber retrieves cached yearly failed withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusFailed: List of yearly failed withdraw statistics.
//   - bool: Whether the cache was found.
func (w *withdrawStatsByCardStatusCache) GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearWithdrawStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawResponseYearStatusFailed](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyWithdrawStatusFailedByCardNumber stores yearly failed withdraw statistics in the cache by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//   - data: The data to cache.
func (w *withdrawStatsByCardStatusCache) SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusFailedByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

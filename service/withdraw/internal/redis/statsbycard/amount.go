package withdrawstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type withdrawStatsByCardAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsAmountCache(store *sharedcachehelpers.CacheStore) WithdrawStatsByCardAmountCache {
	return &withdrawStatsByCardAmountCache{store: store}
}

// GetCachedMonthlyWithdrawsByCardNumber retrieves cached monthly withdraw amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year, month, and card number.
//
// Returns:
//   - []*response.WithdrawMonthlyAmountResponse: List of monthly withdraw amounts.
//   - bool: Whether the cache was found.
func (w *withdrawStatsByCardAmountCache) GetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, bool) {
	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawMonthlyAmountResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedMonthlyWithdrawsByCardNumber stores monthly withdraw amounts in the cache by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year, month, and card number.
//   - data: The data to cache.
func (w *withdrawStatsByCardAmountCache) SetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*response.WithdrawMonthlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

// GetCachedYearlyWithdrawsByCardNumber retrieves cached yearly withdraw amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.WithdrawYearlyAmountResponse: List of yearly withdraw amounts.
//   - bool: Whether the cache was found.
func (w *withdrawStatsByCardAmountCache) GetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.WithdrawYearlyAmountResponse](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedYearlyWithdrawsByCardNumber stores yearly withdraw amounts in the cache by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//   - data: The data to cache.
func (w *withdrawStatsByCardAmountCache) SetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*response.WithdrawYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

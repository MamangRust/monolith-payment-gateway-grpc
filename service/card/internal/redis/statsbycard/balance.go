package cardstatsbycardmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// cardStatsBalanceByCardCache is a struct that represents the cache store
type cardStatsBalanceByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewCardStatsBalanceByCardCache creates a new instance of cardStatsBalanceByCardCache.
//
// Parameters:
//   - store: The underlying cache store to use.
//
// Returns:
//   - A pointer to a newly created instance of cardStatsBalanceByCardCache.
func NewCardStatsBalanceByCardCache(store *sharedcachehelpers.CacheStore) CardStatsBalanceByCardCache {
	return &cardStatsBalanceByCardCache{store: store}
}

// GetMonthlyBalanceByNumberCache retrieves the cached monthly balance statistics
// for a specific card number based on the given month and year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//
// Returns:
//   - []*response.CardResponseMonthBalance: Slice of monthly balance statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsBalanceByCardCache) GetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalanceByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthBalance](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyBalanceByNumberCache stores the monthly balance statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//   - data: The data to be cached.
func (c *cardStatsBalanceByCardCache) SetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalanceByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetYearlyBalanceByNumberCache retrieves the cached yearly balance statistics
// for a specific card number based on the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//
// Returns:
//   - []*response.CardResponseYearlyBalance: Slice of yearly balance statistics for the specified card.
//   - bool: Whether the data was found in the cache.

func (c *cardStatsBalanceByCardCache) GetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalanceByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearlyBalance](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyBalanceByNumberCache stores the yearly balance statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//   - data: The data to be cached.
func (c *cardStatsBalanceByCardCache) SetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearlyBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalanceByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

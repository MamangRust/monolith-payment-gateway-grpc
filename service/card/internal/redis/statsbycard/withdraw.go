package cardstatsbycardmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsWithdrawByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsWithdrawByCardCache(store *sharedcachehelpers.CacheStore) CardStatsWithdrawByCardCache {
	return &cardStatsWithdrawByCardCache{store: store}
}

// GetMonthlyWithdrawByNumberCache retrieves the cached monthly withdraw statistics
// for a specific card number based on the given month and year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly withdraw statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsWithdrawByCardCache) GetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyWithdrawByNumberCache stores the monthly withdraw statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//   - data: The data to be cached.

func (c *cardStatsWithdrawByCardCache) SetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetYearlyWithdrawByNumberCache retrieves the cached yearly withdraw statistics
// for a specific card number based on the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly withdraw statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsWithdrawByCardCache) GetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyWithdrawByNumberCache stores the yearly withdraw statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//   - data: The data to be cached.
func (c *cardStatsWithdrawByCardCache) SetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

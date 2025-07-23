package cardstatsbycardmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsTopupByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTopupByCardCache(store *sharedcachehelpers.CacheStore) CardStatsTopupByCardCache {
	return &cardStatsTopupByCardCache{store: store}
}

// GetMonthlyTopupByNumberCache retrieves the cached monthly top-up statistics
// for a specific card number based on the given month and year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly top-up statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTopupByCardCache) GetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTopupByNumberCache stores the monthly top-up statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//   - data: The data to be cached.
func (c *cardStatsTopupByCardCache) SetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetYearlyTopupByNumberCache retrieves the cached yearly top-up statistics
// for a specific card number based on the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly top-up statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTopupByCardCache) GetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupByNumberCache stores the yearly top-up statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//   - data: The data to be cached.
func (c *cardStatsTopupByCardCache) SetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

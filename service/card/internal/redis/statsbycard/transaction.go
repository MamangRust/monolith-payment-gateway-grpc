package cardstatsbycardmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsTransactionByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransactionByCardCache(store *sharedcachehelpers.CacheStore) CardStatsTransactionByCardCache {
	return &cardStatsTransactionByCardCache{store: store}
}

// GetMonthlyTransactionByNumberCache retrieves the cached monthly transaction statistics
// for a specific card number based on the given month and year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly transaction statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransactionByCardCache) GetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransactionByNumberCache stores the monthly transaction statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number, month, and year.
//   - data: The data to be cached.
func (c *cardStatsTransactionByCardCache) SetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetYearlyTransactionByNumberCache retrieves the cached yearly transaction statistics
// for a specific card number based on the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly transaction statistics for the specified card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransactionByCardCache) GetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransactionByNumberCache stores the yearly transaction statistics
// for a specific card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the card number and year.
//   - data: The data to be cached.
func (c *cardStatsTransactionByCardCache) SetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

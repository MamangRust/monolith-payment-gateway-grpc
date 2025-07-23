package cardstatsbycardmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsTransferByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransferByCardCache(store *sharedcachehelpers.CacheStore) CardStatsTransferByCardCache {
	return &cardStatsTransferByCardCache{store: store}
}

// GetMonthlyTransferBySenderCache retrieves the cached monthly transfer-out statistics
// for a specific sender card number based on the given month and year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the sender card number, month, and year.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-out statistics for the sender card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferByCardCache) GetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransferBySenderCache stores the monthly transfer-out statistics
// for a specific sender card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the sender card number, month, and year.
//   - data: The data to be cached.
func (c *cardStatsTransferByCardCache) SetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetYearlyTransferBySenderCache retrieves the cached yearly transfer-out statistics
// for a specific sender card number based on the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the sender card number and year.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly transfer-out statistics for the sender card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferByCardCache) GetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransferBySenderCache stores the yearly transfer-out statistics
// for a specific sender card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the sender card number and year.
//   - data: The data to be cached.
func (c *cardStatsTransferByCardCache) SetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetMonthlyTransferByReceiverCache retrieves the cached monthly transfer-in statistics
// for a specific receiver card number based on the given month and year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the receiver card number, month, and year.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-in statistics for the receiver card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferByCardCache) GetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransferByReceiverCache stores the monthly transfer-in statistics
// for a specific receiver card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the receiver card number, month, and year.
//   - data: The data to be cached.
func (c *cardStatsTransferByCardCache) SetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

// GetYearlyTransferByReceiverCache retrieves the cached yearly transfer-in statistics
// for a specific receiver card number based on the given year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the receiver card number and year.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly transfer-in statistics for the receiver card.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferByCardCache) GetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransferByReceiverCache stores the yearly transfer-in statistics
// for a specific receiver card number in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: A request object containing the receiver card number and year.
//   - data: The data to be cached.
func (c *cardStatsTransferByCardCache) SetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

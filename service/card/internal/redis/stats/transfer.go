package cardstatsmencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type cardStatsTransferCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransferCache(store *sharedcachehelpers.CacheStore) CardStatsTransferCache {
	return &cardStatsTransferCache{store: store}
}

// GetMonthlyTransferSenderCache retrieves the global monthly transfer-out statistics
// (across all sender card numbers) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly transfer-out data is requested.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-out statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferCache) GetMonthlyTransferSenderCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransferSenderCache stores the global monthly transfer-out statistics
// (across all sender card numbers) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTransferCache) SetMonthlyTransferSenderCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetYearlyTransferSenderCache retrieves the global yearly transfer-out statistics
// (across all sender card numbers) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly transfer-out data is requested.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly transfer-out statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferCache) GetYearlyTransferSenderCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransferSenderCache stores the global yearly transfer-out statistics
// (across all sender card numbers) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTransferCache) SetYearlyTransferSenderCache(ctx context.Context, year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetMonthlyTransferReceiverCache retrieves the global monthly transfer-in statistics
// (across all receiver card numbers) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly transfer-in data is requested.
//
// Returns:
//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-in statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferCache) GetMonthlyTransferReceiverCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseMonthAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransferReceiverCache stores the global monthly transfer-in statistics
// (across all receiver card numbers) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTransferCache) SetMonthlyTransferReceiverCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

// GetYearlyTransferReceiverCache retrieves the global yearly transfer-in statistics
// (across all receiver card numbers) for a given year from the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which yearly transfer-in data is requested.
//
// Returns:
//   - []*response.CardResponseYearAmount: Slice of yearly transfer-in statistics.
//   - bool: Whether the data was found in the cache.
func (c *cardStatsTransferCache) GetYearlyTransferReceiverCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.CardResponseYearAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransferReceiverCache stores the global yearly transfer-in statistics
// (across all receiver card numbers) for a given year in the cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is being cached.
//   - data: The data to be cached.
func (c *cardStatsTransferCache) SetYearlyTransferReceiverCache(ctx context.Context, year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

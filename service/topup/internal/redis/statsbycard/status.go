package topupstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type topupStatsStatusByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsStatusByCardCache(store *sharedcachehelpers.CacheStore) TopupStatsStatusByCardCache {
	return &topupStatsStatusByCardCache{store: store}
}

// GetMonthTopupStatusSuccessCache retrieves cached monthly topup statistics with status "success".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and optional month filter.
//
// Returns:
//   - []*response.TopupResponseMonthStatusSuccess: List of monthly successful topup responses.
//   - bool: Whether the cache was found.
func (s *topupStatsStatusByCardCache) GetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseMonthStatusSuccess](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTopupStatusSuccessCache stores the monthly successful topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The original request used as the cache key.
//   - data: The data to be cached.
func (s *topupStatsStatusByCardCache) SetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyTopupStatusSuccessCache retrieves cached yearly topup statistics with status "success".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the statistics.
//
// Returns:
//   - []*response.TopupResponseYearStatusSuccess: List of yearly successful topup responses.
//   - bool: Whether the cache was found.
func (s *topupStatsStatusByCardCache) GetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseYearStatusSuccess](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupStatusSuccessCache stores yearly successful topup statistics in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year of the data.
//   - data: The data to cache.
func (s *topupStatsStatusByCardCache) SetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetMonthTopupStatusFailedCache retrieves cached monthly topup statistics with status "failed".
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and optional month filter.
//
// Returns:
//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed topup responses.
//   - bool: Whether the cache was found.
func (s *topupStatsStatusByCardCache) GetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseMonthStatusFailed](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTopupStatusFailedByCardNumberCache stores the monthly topup status failed data
// for a specific card number in the cache. It takes as argument a request containing the
// card number, month, and year, and a slice of TopupResponseMonthStatusFailed.
//
// If the provided data is nil, it returns immediately.
//
// It constructs a cache key using the request's parameters and stores the data in the
// cache with a default TTL.
func (s *topupStatsStatusByCardCache) SetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyTopupStatusFailedByCardNumberCache retrieves the yearly topup status failed data
// for a specific card number from the cache. It takes as an argument a request containing the
// card number and year, and returns a slice of TopupResponseYearStatusFailed and a boolean
// indicating whether the data was found in the cache.
func (s *topupStatsStatusByCardCache) GetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupResponseYearStatusFailed](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupStatusFailedByCardNumberCache stores the yearly topup status failed data
// for a specific card number in the cache. It takes as argument a request containing the
// card number and year, and a slice of TopupResponseYearStatusFailed.
//
// If the provided data is nil, it returns immediately.
//
// It constructs a cache key using the request's parameters and stores the data in the
// cache with a default TTL.
func (s *topupStatsStatusByCardCache) SetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

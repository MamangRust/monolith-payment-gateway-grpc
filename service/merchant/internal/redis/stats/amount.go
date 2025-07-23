package merchantstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsAmountCache(store *sharedcachehelpers.CacheStore) MerchantStatsAmountCache {
	return &merchantStatsAmountCache{store: store}
}

// GetMonthlyAmountMerchantCache retrieves monthly amount statistics from the cache
// using the given year as part of the cache key. It returns the cached data if found
// and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	year: An integer representing the year for which the amount statistics
//	      should be retrieved from the cache.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyAmount: A slice of pointers to the monthly
//	                                           amount statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (s *merchantStatsAmountCache) GetMonthlyAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyAmount](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyAmountMerchantCache stores monthly amount statistics in the cache
// using the given year as part of the cache key. If the data provided is nil,
// the function returns without storing anything.
//
// Parameters:
//
//	year: An integer representing the year for generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyAmount containing the
//	      amount statistics to be cached.
func (s *merchantStatsAmountCache) SetMonthlyAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseMonthlyAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyAmountMerchantCache retrieves yearly amount statistics from the cache
// using the given year as part of the cache key. It returns the cached data
// if found and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	year: An integer representing the year for which the amount statistics
//	      should be retrieved from the cache.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyAmount: A slice of pointers to the yearly
//	                                           amount statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (s *merchantStatsAmountCache) GetYearlyAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseYearlyAmount, bool) {
	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyAmount](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyAmountMerchantCache stores yearly amount statistics in the cache
// using the provided year as part of the cache key. If the data provided is nil,
// the function returns without storing anything.
//
// Parameters:
//
//	year: An integer representing the year for generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyAmount containing the
//	      amount statistics to be cached.
func (s *merchantStatsAmountCache) SetYearlyAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseYearlyAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

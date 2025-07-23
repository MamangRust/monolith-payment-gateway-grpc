package merchantstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsTotalAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsTotalAmountCache(store *sharedcachehelpers.CacheStore) MerchantStatsTotalAmountCache {
	return &merchantStatsTotalAmountCache{store: store}
}

// GetMonthlyTotalAmountMerchantCache retrieves monthly total amount statistics
// from the cache using the given year as part of the cache key. It returns the
// cached data if found and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	year: An integer representing the year for which the total amount statistics
//	      should be retrieved from the cache.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyTotalAmount: A slice of pointers to the
//	                                                  monthly total amount
//	                                                  statistics retrieved from the
//	                                                  cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (s *merchantStatsTotalAmountCache) GetMonthlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyTotalAmount](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTotalAmountMerchantCache stores monthly total amount statistics in the cache
// using the given year as part of the cache key. If the data provided is nil, the function
// returns without storing anything.
//
// Parameters:
//
//	year: An integer representing the year for which the total amount statistics
//	      should be stored in the cache.
//	data: A slice of pointers to MerchantResponseMonthlyTotalAmount containing the
//	      total amount statistics to be cached.
func (s *merchantStatsTotalAmountCache) SetMonthlyTotalAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseMonthlyTotalAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyTotalAmountMerchantCache retrieves yearly total amount statistics from the cache
// using the given year as part of the cache key. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	year: An integer representing the year for which the total amount statistics
//	      should be retrieved from the cache.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyTotalAmount: A slice of pointers to the yearly
//	                                                  total amount statistics retrieved
//	                                                  from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (s *merchantStatsTotalAmountCache) GetYearlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseYearlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyTotalAmount](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTotalAmountMerchantCache stores yearly total amount statistics in the cache
// using the provided year as part of the cache key. If the data provided is nil, the
// function returns without storing anything.
//
// Parameters:
//
//	year: An integer representing the year for generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyTotalAmount containing the
//	      total amount statistics to be cached.
func (s *merchantStatsTotalAmountCache) SetYearlyTotalAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseYearlyTotalAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

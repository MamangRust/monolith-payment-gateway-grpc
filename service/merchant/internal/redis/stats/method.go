package merchantstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsMethodCache(store *sharedcachehelpers.CacheStore) MerchantStatsMethodCache {
	return &merchantStatsMethodCache{store: store}
}

// GetMonthlyPaymentMethodsMerchantCache retrieves monthly payment method statistics
// from the cache using the given year as part of the cache key. It returns the
// cached data if found and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	year: An integer representing the year for which the payment method statistics
//	      should be retrieved from the cache.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyPaymentMethod: A slice of pointers to the
//	                                                  monthly payment method
//	                                                  statistics retrieved from the
//	                                                  cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (s *merchantStatsMethodCache) GetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyPaymentMethod](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyPaymentMethodsMerchantCache stores monthly payment method statistics
// in the cache using the given year as part of the cache key. If the data
// provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	year: An integer representing the year for generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyPaymentMethod containing
//	      the payment method statistics to be cached.
func (s *merchantStatsMethodCache) SetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseMonthlyPaymentMethod) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyPaymentMethodMerchantCache retrieves yearly payment method statistics
// from the cache using the given year as part of the cache key. It returns the
// cached data if found and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	year: An integer representing the year for which the payment method statistics
//	      should be retrieved from the cache.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyPaymentMethod: A slice of pointers to the
//	                                                  yearly payment method
//	                                                  statistics retrieved from the
//	                                                  cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (s *merchantStatsMethodCache) GetYearlyPaymentMethodMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseYearlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyPaymentMethod](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyPaymentMethodMerchantCache stores yearly payment method statistics
// in the cache using the given year as part of the cache key. If the data
// provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	year: An integer representing the year for which the payment method statistics
//	      should be retrieved from the cache.
//	data: A slice of pointers to MerchantResponseYearlyPaymentMethod containing
//	      the payment method statistics to be cached.
func (s *merchantStatsMethodCache) SetYearlyPaymentMethodMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseYearlyPaymentMethod) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

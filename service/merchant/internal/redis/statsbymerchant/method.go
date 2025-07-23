package merchantstatsbymerchant

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsMethodByMerchant struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsMethodByMerchantCache(store *sharedcachehelpers.CacheStore) MerchantStatsMethodByMerchantCache {
	return &merchantStatsMethodByMerchant{store: store}
}

// GetMonthlyPaymentMethodByMerchantsCache retrieves monthly payment method statistics from the cache
// using the merchant ID and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to MonthYearPaymentMethodMerchant containing the merchant ID and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyPaymentMethod: A slice of pointers to the monthly payment
//	                                                   method statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsMethodByMerchant) GetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyPaymentMethod](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyPaymentMethodByMerchantsCache stores monthly payment method statistics
// in the cache using the merchant ID and year as part of the cache key. If the data
// provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to MonthYearPaymentMethodMerchant containing the merchant ID and
//	     year for generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyPaymentMethod containing
//	      the payment method statistics to be cached.
func (m *merchantStatsMethodByMerchant) SetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseMonthlyPaymentMethod) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetYearlyPaymentMethodByMerchantsCache retrieves yearly payment method statistics from the cache
// using the merchant ID and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to MonthYearPaymentMethodMerchant containing the merchant ID and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyPaymentMethod: A slice of pointers to the yearly payment
//	                                                   method statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsMethodByMerchant) GetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyPaymentMethod](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyPaymentMethodByMerchantsCache stores yearly payment method statistics
// in the cache using the merchant ID and year as part of the cache key. If the data
// provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to MonthYearPaymentMethodMerchant containing the merchant ID and
//	     year for generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyPaymentMethod containing
//	      the payment method statistics to be cached.
func (m *merchantStatsMethodByMerchant) SetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseYearlyPaymentMethod) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

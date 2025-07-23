package merchantstatsapikey

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsMethodByApiKeyCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsMethodByApiKeyCache(store *sharedcachehelpers.CacheStore) MerchantStatsMethodByApiKeyCache {
	return &merchantStatsMethodByApiKeyCache{store: store}
}

// GetMonthlyPaymentMethodByApikeysCache retrieves monthly payment method statistics from the cache
// using the API key and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearPaymentMethodApiKey containing the API key and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyPaymentMethod: A slice of pointers to the monthly payment
//	                                                   method statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsMethodByApiKeyCache) GetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyPaymentMethod](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyPaymentMethodByApikeysCache stores monthly payment method statistics
// in the cache using the API key and year as part of the cache key. If the data
// provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearPaymentMethodApiKey containing the API key and
//	     year for generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyPaymentMethod containing
//	      the payment method statistics to be cached.
func (m *merchantStatsMethodByApiKeyCache) SetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseMonthlyPaymentMethod) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetYearlyPaymentMethodByApikeysCache retrieves yearly payment method statistics from the cache
// using the API key and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearPaymentMethodApiKey containing the API key and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyPaymentMethod: A slice of pointers to the yearly payment
//	                                                   method statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsMethodByApiKeyCache) GetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyPaymentMethod](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyPaymentMethodByApikeysCache stores yearly payment method statistics in the cache using
// the API key and year as part of the cache key. If the data provided is nil, the function returns
// without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearPaymentMethodApiKey containing the API key and year for
//	     generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyPaymentMethod containing the payment
//	      method statistics to be cached.
func (m *merchantStatsMethodByApiKeyCache) SetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseYearlyPaymentMethod) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

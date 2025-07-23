package merchantstatsapikey

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsTotalAmountByApiKeyCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsTotalAmountByApiKeyCache(store *sharedcachehelpers.CacheStore) MerchantStatsTotalAmountByApiKeyCache {
	return &merchantStatsTotalAmountByApiKeyCache{store: store}
}

// GetMonthlyTotalAmountByApikeysCache retrieves monthly total amount statistics from the cache
// using the API key and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountApiKey containing the API key and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyTotalAmount: A slice of pointers to the monthly total amount
//	                                                   statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsTotalAmountByApiKeyCache) GetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyTotalAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTotalAmountByApikeysCache stores monthly total amount statistics in the cache using the API key and year as part of the cache key. If the data provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountApiKey containing the API key and year for generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyTotalAmount containing the total amount statistics to be cached.
func (m *merchantStatsTotalAmountByApiKeyCache) SetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseMonthlyTotalAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetYearlyTotalAmountByApikeysCache retrieves yearly total amount statistics from the cache
// using the API key and year provided in the request. It returns the cached data if found
// and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountApiKey containing the API key and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyTotalAmount: A slice of pointers to the yearly total amount
//	                                                   statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsTotalAmountByApiKeyCache) GetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyTotalAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTotalAmountByApikeysCache stores yearly total amount statistics in the cache using the API key and year as part of the cache key. If the data provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountApiKey containing the API key and year for generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyTotalAmount containing the total amount statistics to be cached.
func (m *merchantStatsTotalAmountByApiKeyCache) SetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseYearlyTotalAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

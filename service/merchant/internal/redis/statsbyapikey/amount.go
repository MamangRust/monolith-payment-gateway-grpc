package merchantstatsapikey

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsAmountByApiKeyCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsAmountByApiKeyCache(store *sharedcachehelpers.CacheStore) MerchantStatsAmountByApiKeyCache {
	return &merchantStatsAmountByApiKeyCache{store: store}
}

// GetMonthlyAmountByApikeysCache retrieves monthly amount statistics from the cache
// using the API key and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearAmountApiKey containing the API key and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyAmount: A slice of pointers to the monthly amount
//	                                                   statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsAmountByApiKeyCache) GetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyAmountByApikeysCache stores monthly amount statistics in the cache using the API
// key and year as part of the cache key. If the data provided is nil, the function returns
// without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearAmountApiKey containing the API key and year for generating
//	     the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyAmount containing the amount
//	      statistics to be cached.
func (m *merchantStatsAmountByApiKeyCache) SetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseMonthlyAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetYearlyAmountByApikeysCache retrieves yearly amount statistics from the cache
// using the API key and year provided in the request. It returns the cached data if found
// and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearAmountApiKey containing the API key and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyAmount: A slice of pointers to the yearly amount
//	                                           statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsAmountByApiKeyCache) GetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyAmountByApikeysCache stores yearly amount statistics in the cache using the API
// key and year as part of the cache key. If the data provided is nil, the function returns
// without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearAmountApiKey containing the API key and year for
//	     generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyAmount containing the amount
//	      statistics to be cached.
func (m *merchantStatsAmountByApiKeyCache) SetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseYearlyAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

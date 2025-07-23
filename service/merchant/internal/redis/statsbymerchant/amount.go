package merchantstatsbymerchant

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsAmountByMerchant struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsAmountByMerchantCache(store *sharedcachehelpers.CacheStore) MerchantStatsAmountByMerchantCache {
	return &merchantStatsAmountByMerchant{store: store}
}

// GetMonthlyAmountByMerchantsCache retrieves monthly amount statistics from the cache
// using the merchant ID and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to MonthYearAmountMerchant containing the merchant ID and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyAmount: A slice of pointers to the monthly amount
//	                                                   statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsAmountByMerchant) GetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyAmountByMerchantsCache stores monthly amount statistics
// in the cache using the merchant ID and year as part of the cache key. If the data
// provided is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to MonthYearAmountMerchant containing the merchant ID and
//	     year for generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyAmount containing
//	      the amount statistics to be cached.
func (m *merchantStatsAmountByMerchant) SetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseMonthlyAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetYearlyAmountByMerchantsCache retrieves yearly amount statistics from the cache
// using the merchant ID and year provided in the request. It returns the cached data
// if found and valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to MonthYearAmountMerchant containing the merchant ID and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyAmount: A slice of pointers to the yearly amount
//	                                          statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsAmountByMerchant) GetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyAmountByMerchantsCache stores yearly amount statistics in the cache
// using the merchant ID and year as part of the cache key. If the data provided
// is nil, the function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to MonthYearAmountMerchant containing the merchant ID and
//	     year for generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyAmount containing the
//	      amount statistics to be cached.
func (m *merchantStatsAmountByMerchant) SetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseYearlyAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

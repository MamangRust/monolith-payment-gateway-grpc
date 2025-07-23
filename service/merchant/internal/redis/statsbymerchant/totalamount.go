package merchantstatsbymerchant

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type merchantStatsTotalAmountByMerchant struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsTotalAmountByMerchantCache(store *sharedcachehelpers.CacheStore) MerchantStatsTotalAmountByMerchantCache {
	return &merchantStatsTotalAmountByMerchant{store: store}
}

// SetMonthlyTotalAmountByMerchantsCache stores monthly total amount statistics in the cache
// using the merchant ID and year as part of the cache key. If the data provided is nil, the
// function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountMerchant containing the merchant ID and year for
//	     generating the cache key.
//	data: A slice of pointers to MerchantResponseMonthlyTotalAmount containing the total
//	      amount statistics to be cached.
func (m *merchantStatsTotalAmountByMerchant) SetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseMonthlyTotalAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetMonthlyTotalAmountByMerchantsCache retrieves monthly total amount statistics from the cache
// using the merchant ID and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountMerchant containing the merchant ID and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseMonthlyTotalAmount: A slice of pointers to the monthly total amount
//	                                                   statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsTotalAmountByMerchant) GetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseMonthlyTotalAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTotalAmountByMerchantsCache stores yearly total amount statistics in the cache
// using the merchant ID and year as part of the cache key. If the data provided is nil, the
// function returns without storing anything.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountMerchant containing the merchant ID and year for
//	     generating the cache key.
//	data: A slice of pointers to MerchantResponseYearlyTotalAmount containing the total
//	      amount statistics to be cached.
func (m *merchantStatsTotalAmountByMerchant) SetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseYearlyTotalAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// GetYearlyTotalAmountByMerchantsCache retrieves yearly total amount statistics from the cache
// using the merchant ID and year provided in the request. It returns the cached data if found and
// valid; otherwise, it returns nil and false.
//
// Parameters:
//
//	req: A pointer to a MonthYearTotalAmountMerchant containing the merchant ID and year used to
//	     generate the cache key.
//
// Returns:
//
//	[]*response.MerchantResponseYearlyTotalAmount: A slice of pointers to the yearly total amount
//	                                                   statistics retrieved from the cache.
//	bool: A boolean indicating whether the cache was found and contains valid data.
func (m *merchantStatsTotalAmountByMerchant) GetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.MerchantResponseYearlyTotalAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

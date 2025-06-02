package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	merchantMonthlyPaymentMethodByMerchantCacheKey = "merchant:statistic:monthly:payment-method:merchant-apikey:%d:year:%d"
	merchantYearlyPaymentMethodByMerchantCacheKey  = "merchant:statistic:yearly:payment-method:merchant-apikey:%d:year:%d"

	merchantMonthlyAmountByMerchantCacheKey = "merchant:statistic:monthly:amount:merchant-apikey:%d:year:%d"
	merchantYearlyAmountByMerchantCacheKey  = "merchant:statistic:yearly:amount:merchant-apikey:%d:year:%d"

	merchantMonthlyTotalAmountByMerchantCacheKey = "merchant:statistic:monthly:total-amount:merchant-apikey:%d:year:%d"
	merchantYearlyTotalAmountByMerchantCacheKey  = "merchant:statistic:yearly:total-amount:merchant-apikey:%d:year:%d"
)

type merchantStatisticByMerchantCache struct {
	store *CacheStore
}

func NewMerchantStatisticByMerchantCache(store *CacheStore) *merchantStatisticByMerchantCache {
	return &merchantStatisticByMerchantCache{store: store}
}

func (m *merchantStatisticByMerchantCache) SetMonthlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseMonthlyPaymentMethod) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByMerchantCache) GetMonthlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant) []*response.MerchantResponseMonthlyPaymentMethod {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyPaymentMethod](m.store, key)
	if !found {
		return nil
	}
	return *result
}

func (m *merchantStatisticByMerchantCache) SetYearlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseYearlyPaymentMethod) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByMerchantCache) GetYearlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant) []*response.MerchantResponseYearlyPaymentMethod {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyPaymentMethod](m.store, key)
	if !found {
		return nil
	}
	return *result
}

func (m *merchantStatisticByMerchantCache) SetMonthlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseMonthlyAmount) {
	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByMerchantCache) GetMonthlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant) []*response.MerchantResponseMonthlyAmount {
	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyAmount](m.store, key)
	if !found {
		return nil
	}
	return *result
}

func (m *merchantStatisticByMerchantCache) SetYearlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseYearlyAmount) {
	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByMerchantCache) GetYearlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant) []*response.MerchantResponseYearlyAmount {
	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyAmount](m.store, key)
	if !found {
		return nil
	}
	return *result
}

func (m *merchantStatisticByMerchantCache) SetMonthlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseMonthlyTotalAmount) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByMerchantCache) GetMonthlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant) []*response.MerchantResponseMonthlyTotalAmount {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyTotalAmount](m.store, key)
	if !found {
		return nil
	}
	return *result
}

func (m *merchantStatisticByMerchantCache) SetYearlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseYearlyTotalAmount) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByMerchantCache) GetYearlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant) []*response.MerchantResponseYearlyTotalAmount {
	key := fmt.Sprintf(merchantYearlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyTotalAmount](m.store, key)
	if !found {
		return nil
	}
	return *result
}

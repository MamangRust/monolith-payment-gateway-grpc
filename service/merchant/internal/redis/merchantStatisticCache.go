package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	merchantMonthlyPaymentMethodCacheKey = "merchant:statistic:monthly:payment-method:year:%d"
	merchantYearlyPaymentMethodCacheKey  = "merchant:statistic:yearly:payment-method:year:%d"

	merchantMonthlyAmountCacheKey = "merchant:statistic:monthly:amount:year:%d"
	MerchantYearlyAmountCacheKey  = "merchant:statistic:yearly:amount:year:%d"

	merchantMonthlyTotalAmountCacheKey = "merchant:statistic:monthly:total-amount:year:%d"
	merchantYearlyTotalAmountCacheKey  = "merchant:statistic:yearly:total-amount:year:%d"
)

type merchantStatisticCache struct {
	store *CacheStore
}

func NewMerchantStatisticCache(store *CacheStore) *merchantStatisticCache {
	return &merchantStatisticCache{store: store}
}

func (s *merchantStatisticCache) GetMonthlyPaymentMethodsMerchantCache(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyPaymentMethod](s.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (s *merchantStatisticCache) SetMonthlyPaymentMethodsMerchantCache(year int, data []*response.MerchantResponseMonthlyPaymentMethod) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *merchantStatisticCache) GetYearlyPaymentMethodMerchantCache(year int) ([]*response.MerchantResponseYearlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyPaymentMethod](s.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (s *merchantStatisticCache) SetYearlyPaymentMethodMerchantCache(year int, data []*response.MerchantResponseYearlyPaymentMethod) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *merchantStatisticCache) GetMonthlyAmountMerchantCache(year int) ([]*response.MerchantResponseMonthlyAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyAmount](s.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (s *merchantStatisticCache) SetMonthlyAmountMerchantCache(year int, data []*response.MerchantResponseMonthlyAmount) {
	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *merchantStatisticCache) GetYearlyAmountMerchantCache(year int) ([]*response.MerchantResponseYearlyAmount, bool) {
	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyAmount](s.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (s *merchantStatisticCache) SetYearlyAmountMerchantCache(year int, data []*response.MerchantResponseYearlyAmount) {
	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *merchantStatisticCache) GetMonthlyTotalAmountMerchantCache(year int) ([]*response.MerchantResponseMonthlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyTotalAmount](s.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (s *merchantStatisticCache) SetMonthlyTotalAmountMerchantCache(year int, data []*response.MerchantResponseMonthlyTotalAmount) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *merchantStatisticCache) GetYearlyTotalAmountMerchantCache(year int) ([]*response.MerchantResponseYearlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyTotalAmount](s.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (s *merchantStatisticCache) SetYearlyTotalAmountMerchantCache(year int, data []*response.MerchantResponseYearlyTotalAmount) {
	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	SetToCache(s.store, key, &data, ttlDefault)
}

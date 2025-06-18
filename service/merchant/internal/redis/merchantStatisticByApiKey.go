package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	merchantMonthlyPaymentMethodByApikeyCacheKey = "merchant:statistic:monthly:payment-method:apikey:%s:year:%d"
	merchantYearlyPaymentMethodByApikeyCacheKey  = "merchant:statistic:yearly:payment-method:apikey:%s:year:%d"

	merchantMonthlyAmountByApikeyCacheKey = "merchant:statistic:monthly:amount:apikey:%s:year:%d"
	merchantYearlyAmountByApikeyCacheKey  = "merchant:statistic:yearly:amount:apikey:%s:year:%d"

	merchantMonthlyTotalAmountByApikeyCacheKey = "merchant:statistic:monthly:total-amount:apikey:%s:year:%d"
	merchantYearlyTotalAmountByApikeyCacheKey  = "merchant:statistic:yearly:total-amount:apikey:%s:year:%d"
)

type merchantStatisticByApiKeyCache struct {
	store *CacheStore
}

func NewMerchantStatisticByApiKeyCache(store *CacheStore) *merchantStatisticByApiKeyCache {
	return &merchantStatisticByApiKeyCache{store: store}
}

func (m *merchantStatisticByApiKeyCache) SetMonthlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseMonthlyPaymentMethod) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByApiKeyCache) GetMonthlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyPaymentMethod](m.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (m *merchantStatisticByApiKeyCache) SetYearlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseYearlyPaymentMethod) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByApiKeyCache) GetYearlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyPaymentMethod](m.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (m *merchantStatisticByApiKeyCache) SetMonthlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseMonthlyAmount) {
	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByApiKeyCache) GetMonthlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyAmount](m.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (m *merchantStatisticByApiKeyCache) SetYearlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseYearlyAmount) {
	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByApiKeyCache) GetYearlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyAmount](m.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (m *merchantStatisticByApiKeyCache) SetMonthlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseMonthlyTotalAmount) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByApiKeyCache) GetMonthlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseMonthlyTotalAmount](m.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (m *merchantStatisticByApiKeyCache) SetYearlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseYearlyTotalAmount) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	SetToCache(m.store, key, &data, ttlDefault)
}

func (m *merchantStatisticByApiKeyCache) GetYearlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := GetFromCache[[]*response.MerchantResponseYearlyTotalAmount](m.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	monthTopupStatusSuccessCacheKey = "topup:month:status:success:month:%d:year:%d"
	yearTopupStatusSuccessCacheKey  = "topup:year:status:success:year:%d"
	monthTopupStatusFailedCacheKey  = "topup:month:status:failed:month:%d:year:%d"
	yearTopupStatusFailedCacheKey   = "topup:year:status:failed:year:%d"

	monthTopupAmountCacheKey = "topup:month:amount:year:%d"
	yearTopupAmountCacheKey  = "topup:year:amount:year:%d"

	monthTopupMethodCacheKey = "topup:month:method:year:%d"
	yearTopupMethodCacheKey  = "topup:year:method:year:%d"
)

type topupStatisticCache struct {
	store *CacheStore
}

func NewTopupStatisticCache(store *CacheStore) *topupStatisticCache {
	return &topupStatisticCache{store: store}
}

func (c *topupStatisticCache) GetMonthTopupStatusSuccessCache(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TopupResponseMonthStatusSuccess](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatisticCache) SetMonthTopupStatusSuccessCache(req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetYearlyTopupStatusSuccessCache(year int) ([]*response.TopupResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)

	result, found := GetFromCache[[]*response.TopupResponseYearStatusSuccess](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatisticCache) SetYearlyTopupStatusSuccessCache(year int, data []*response.TopupResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetMonthTopupStatusFailedCache(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TopupResponseMonthStatusFailed](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatisticCache) SetMonthTopupStatusFailedCache(req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetYearlyTopupStatusFailedCache(year int) ([]*response.TopupResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)

	result, found := GetFromCache[[]*response.TopupResponseYearStatusFailed](c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatisticCache) SetYearlyTopupStatusFailedCache(year int, data []*response.TopupResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetMonthlyTopupAmountsCache(year int) ([]*response.TopupMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)

	result, found := GetFromCache[[]*response.TopupMonthAmountResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatisticCache) SetMonthlyTopupAmountsCache(year int, data []*response.TopupMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetYearlyTopupAmountsCache(year int) ([]*response.TopupYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)

	result, found := GetFromCache[[]*response.TopupYearlyAmountResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatisticCache) SetYearlyTopupAmountsCache(year int, data []*response.TopupYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetMonthlyTopupMethodsCache(month int) ([]*response.TopupMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, month)

	result, found := GetFromCache[[]*response.TopupMonthMethodResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatisticCache) SetMonthlyTopupMethodsCache(month int, data []*response.TopupMonthMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodCacheKey, month)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *topupStatisticCache) GetYearlyTopupMethodsCache(year int) ([]*response.TopupYearlyMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)

	result, found := GetFromCache[[]*response.TopupYearlyMethodResponse](c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatisticCache) SetYearlyTopupMethodsCache(year int, data []*response.TopupYearlyMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

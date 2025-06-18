package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	montWithdrawStatusSuccessKey = "withdraws:mont:status:success:month%d:year:%d"
	yearWithdrawStatusSuccessKey = "withdraws:year:status:success:year:%d"

	montWithdrawStatusFailedKey = "withdraws:mont:status:failed:month:%d:year:%d"
	yearWithdrawStatusFailedKey = "withdraws:year:status:failed:year:%d"

	montWithdrawAmountKey = "withdraws:mont:amount:year:%d"
	yearWithdrawAmountKey = "withdraws:year:amount:year:%d"
)

type withdrawStatisticCache struct {
	store *CacheStore
}

func NewWithdrawStatisticCache(store *CacheStore) *withdrawStatisticCache {
	return &withdrawStatisticCache{store: store}
}

func (w *withdrawStatisticCache) GetCachedMonthWithdrawStatusSuccessCache(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.WithdrawResponseMonthStatusSuccess](w.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (w *withdrawStatisticCache) SetCachedMonthWithdrawStatusSuccessCache(req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusSuccess) {
	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticCache) GetCachedYearlyWithdrawStatusSuccessCache(year int) ([]*response.WithdrawResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	result, found := GetFromCache[[]*response.WithdrawResponseYearStatusSuccess](w.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (w *withdrawStatisticCache) SetCachedYearlyWithdrawStatusSuccessCache(year int, data []*response.WithdrawResponseYearStatusSuccess) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticCache) GetCachedMonthWithdrawStatusFailedCache(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	result, found := GetFromCache[[]*response.WithdrawResponseMonthStatusFailed](w.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (w *withdrawStatisticCache) SetCachedMonthWithdrawStatusFailedCache(req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusFailed) {
	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticCache) GetCachedYearlyWithdrawStatusFailedCache(year int) ([]*response.WithdrawResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	result, found := GetFromCache[[]*response.WithdrawResponseYearStatusFailed](w.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (w *withdrawStatisticCache) SetCachedYearlyWithdrawStatusFailedCache(year int, data []*response.WithdrawResponseYearStatusFailed) {
	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticCache) GetCachedMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, bool) {
	key := fmt.Sprintf(montWithdrawAmountKey, year)
	result, found := GetFromCache[[]*response.WithdrawMonthlyAmountResponse](w.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (w *withdrawStatisticCache) SetCachedMonthlyWithdraws(year int, data []*response.WithdrawMonthlyAmountResponse) {
	key := fmt.Sprintf(montWithdrawAmountKey, year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticCache) GetCachedYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	result, found := GetFromCache[[]*response.WithdrawYearlyAmountResponse](w.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (w *withdrawStatisticCache) SetCachedYearlyWithdraws(year int, data []*response.WithdrawYearlyAmountResponse) {
	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	SetToCache(w.store, key, &data, ttlDefault)
}

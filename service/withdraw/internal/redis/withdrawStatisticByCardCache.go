package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	monthWithdrawStatusSuccessByCardKey = "withdraws:month:status:success:card_number:%s:month:%d:year:%d"
	yearWithdrawStatusSuccessByCardKey  = "withdraws:year:status:success:card_number:%s:year:%d"

	monthWithdrawStatusFailedByCardKey = "withdraws:month:status:failed:card_number:%s:month:%d:year:%d"
	yearWithdrawStatusFailedByCardKey  = "withdraws:year:status:failed:card_number:%s:year:%d"

	monthWithdrawAmountByCardKey = "withdraws:month:amount:card_number:%s:year:%d"
	yearWithdrawAmountByCardKey  = "withdraws:year:amount:card_number:%s:year:%d"
)

type withdrawStatisticByCardCache struct {
	store *CacheStore
}

func NewWithdrawStatisticByCardCache(store *CacheStore) *withdrawStatisticByCardCache {
	return &withdrawStatisticByCardCache{store: store}
}

func (w *withdrawStatisticByCardCache) GetCachedMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) []*response.WithdrawResponseMonthStatusSuccess {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := GetFromCache[[]*response.WithdrawResponseMonthStatusSuccess](w.store, key)
	if !found {
		return nil
	}
	return *result
}

func (w *withdrawStatisticByCardCache) SetCachedMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusSuccess) {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticByCardCache) GetCachedYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) []*response.WithdrawResponseYearStatusSuccess {
	key := fmt.Sprintf(yearWithdrawStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.WithdrawResponseYearStatusSuccess](w.store, key)
	if !found {
		return nil
	}
	return *result
}

func (w *withdrawStatisticByCardCache) SetCachedYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusSuccess) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessByCardKey, req.CardNumber, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticByCardCache) GetCachedMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) []*response.WithdrawResponseMonthStatusFailed {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := GetFromCache[[]*response.WithdrawResponseMonthStatusFailed](w.store, key)
	if !found {
		return nil
	}
	return *result
}
func (w *withdrawStatisticByCardCache) SetCachedMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusFailed) {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticByCardCache) GetCachedYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) []*response.WithdrawResponseYearStatusFailed {
	key := fmt.Sprintf(yearWithdrawStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.WithdrawResponseYearStatusFailed](w.store, key)
	if !found {
		return nil
	}
	return *result
}

func (w *withdrawStatisticByCardCache) SetCachedYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusFailed) {
	key := fmt.Sprintf(yearWithdrawStatusFailedByCardKey, req.CardNumber, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticByCardCache) GetCachedMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) []*response.WithdrawMonthlyAmountResponse {
	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.WithdrawMonthlyAmountResponse](w.store, key)
	if !found {
		return nil
	}
	return *result
}

func (w *withdrawStatisticByCardCache) SetCachedMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber, data []*response.WithdrawMonthlyAmountResponse) {
	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

func (w *withdrawStatisticByCardCache) GetCachedYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) []*response.WithdrawYearlyAmountResponse {
	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.WithdrawYearlyAmountResponse](w.store, key)
	if !found {
		return nil
	}
	return *result
}

func (w *withdrawStatisticByCardCache) SetCachedYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber, data []*response.WithdrawYearlyAmountResponse) {
	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	SetToCache(w.store, key, &data, ttlDefault)
}

package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	monthTopupStatusSuccessCacheKey = "transaction:month:status:success:month:%d:year:%d"
	yearTopupStatusSuccessCacheKey  = "transaction:year:status:success:year:%d"
	monthTopupStatusFailedCacheKey  = "transaction:month:status:failed:month:%d:year:%d"
	yearTopupStatusFailedCacheKey   = "transaction:year:status:failed:year:%d"

	monthTopupAmountCacheKey = "transaction:month:amount:year:%d"
	yearTopupAmountCacheKey  = "transaction:year:amount:year:%d"

	monthTopupMethodCacheKey = "transaction:month:method:year:%d"
	yearTopupMethodCacheKey  = "transaction:year:method:year:%d"
)

type transactionStatisticCache struct {
	store *CacheStore
}

func NewTransactionStatisticCache(store *CacheStore) *transactionStatisticCache {
	return &transactionStatisticCache{store: store}
}

func (t *transactionStatisticCache) GetMonthTransactonStatusSuccessCache(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransactionResponseMonthStatusSuccess](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetMonthTransactonStatusSuccessCache(req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusSuccess) {
	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetYearTransactonStatusSuccessCache(year int) ([]*response.TransactionResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	result, found := GetFromCache[[]*response.TransactionResponseYearStatusSuccess](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetYearTransactonStatusSuccessCache(year int, data []*response.TransactionResponseYearStatusSuccess) {
	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetMonthTransactonStatusFailedCache(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransactionResponseMonthStatusFailed](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetMonthTransactonStatusFailedCache(req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusFailed) {
	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetYearTransactonStatusFailedCache(year int) ([]*response.TransactionResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	result, found := GetFromCache[[]*response.TransactionResponseYearStatusFailed](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetYearTransactonStatusFailedCache(year int, data []*response.TransactionResponseYearStatusFailed) {
	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetMonthlyPaymentMethodsCache(year int) ([]*response.TransactionMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, year)
	result, found := GetFromCache[[]*response.TransactionMonthMethodResponse](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetMonthlyPaymentMethodsCache(year int, data []*response.TransactionMonthMethodResponse) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetYearlyPaymentMethodsCache(year int) ([]*response.TransactionYearMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	result, found := GetFromCache[[]*response.TransactionYearMethodResponse](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetYearlyPaymentMethodsCache(year int, data []*response.TransactionYearMethodResponse) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetMonthlyAmountsCache(year int) ([]*response.TransactionMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	result, found := GetFromCache[[]*response.TransactionMonthAmountResponse](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatisticCache) SetMonthlyAmountsCache(year int, data []*response.TransactionMonthAmountResponse) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticCache) GetYearlyAmountsCache(year int) ([]*response.TransactionYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountResponse](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true

}

func (t *transactionStatisticCache) SetYearlyAmountsCache(year int, data []*response.TransactionYearlyAmountResponse) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

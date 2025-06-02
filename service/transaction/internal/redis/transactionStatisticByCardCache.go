package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	monthTopupStatusSuccessByCardCacheKey = "transaction:bycard:month:status:success:card:%s:month:%d:year:%d"
	yearTopupStatusSuccessByCardCacheKey  = "transaction:bycard:year:status:success:card:%s:year:%d"
	monthTopupStatusFailedByCardCacheKey  = "transaction:bycard:month:status:failed:card:%s:month:%d:year:%d"
	yearTopupStatusFailedByCardCacheKey   = "transaction:bycard:year:status:failed:card:%s:year:%d"

	monthTopupAmountByCardCacheKey = "transaction:bycard:month:amount:card:%s:year:%d"
	yearTopupAmountByCardCacheKey  = "transaction:bycard:year:amount:card:%s:year:%d"

	monthTopupMethodByCardCacheKey = "transaction:bycard:month:method:card:%s:year:%d"
	yearTopupMethodByCardCacheKey  = "transaction:bycard:year:method:card:%s:year:%d"
)

type transactionStatisticByCardCache struct {
	store *CacheStore
}

func NewTransactionStatisticByCardCache(store *CacheStore) *transactionStatisticByCardCache {
	return &transactionStatisticByCardCache{store: store}
}

func (t *transactionStatisticByCardCache) GetMonthTransactionStatusSuccessByCardCache(req *requests.MonthStatusTransactionCardNumber) []*response.TransactionResponseMonthStatusSuccess {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransactionResponseMonthStatusSuccess](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetMonthTransactionStatusSuccessByCardCache(req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusSuccess) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetYearTransactionStatusSuccessByCardCache(req *requests.YearStatusTransactionCardNumber) []*response.TransactionResponseYearStatusSuccess {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransactionResponseYearStatusSuccess](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetYearTransactionStatusSuccessByCardCache(req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusSuccess) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetMonthTransactionStatusFailedByCardCache(req *requests.MonthStatusTransactionCardNumber) []*response.TransactionResponseMonthStatusFailed {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransactionResponseMonthStatusFailed](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetMonthTransactionStatusFailedByCardCache(req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusFailed) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetYearTransactionStatusFailedByCardCache(req *requests.YearStatusTransactionCardNumber) []*response.TransactionResponseYearStatusFailed {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransactionResponseYearStatusFailed](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetYearTransactionStatusFailedByCardCache(req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusFailed) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetMonthlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionMonthMethodResponse {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransactionMonthMethodResponse](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetMonthlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthMethodResponse) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetYearlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionYearMethodResponse {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransactionYearMethodResponse](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetYearlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionYearMethodResponse) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetMonthlyAmountsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionMonthAmountResponse {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthAmountResponse](t.store, key)

	if !found {
		return nil
	}
	return *result

}

func (t *transactionStatisticByCardCache) SetMonthlyAmountsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthAmountResponse) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transactionStatisticByCardCache) GetYearlyAmountsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionYearlyAmountResponse {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountResponse](t.store, key)

	if !found {
		return nil
	}
	return *result
}

func (t *transactionStatisticByCardCache) SetYearlyAmountsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionYearlyAmountResponse) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

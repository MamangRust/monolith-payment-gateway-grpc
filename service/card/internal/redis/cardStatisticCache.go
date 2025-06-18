package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	cacheKeyMonthlyBalance = "stat:monthly:balance:%d"
	cacheKeyYearlyBalance  = "stat:yearly:balance:%d"

	cacheKeyMonthlyTopupAmount = "stat:monthly:topup:%d"
	cacheKeyYearlyTopupAmount  = "stat:yearly:topup:%d"

	cacheKeyMonthlyWithdrawAmount = "stat:monthly:withdraw:%d"
	cacheKeyYearlyWithdrawAmount  = "stat:yearly:withdraw:%d"

	cacheKeyMonthlyTransactionAmount = "stat:monthly:transaction:%d"
	cacheKeyYearlyTransactionAmount  = "stat:yearly:transaction:%d"

	cacheKeyMonthlyTransferSender = "stat:monthly:transfer:sender:%d"
	cacheKeyYearlyTransferSender  = "stat:yearly:transfer:sender:%d"

	cacheKeyMonthlyTransferReceiver = "stat:monthly:transfer:receiver:%d"
	cacheKeyYearlyTransferReceiver  = "stat:yearly:transfer:receiver:%d"

	ttlStatistic = 10 * time.Minute
)

type cardStatisticCache struct {
	store *CacheStore
}

func NewCardStatisticCache(store *CacheStore) *cardStatisticCache {
	return &cardStatisticCache{store: store}
}

func (c *cardStatisticCache) GetMonthlyBalanceCache(year int) ([]*response.CardResponseMonthBalance, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)

	result, found := GetFromCache[[]*response.CardResponseMonthBalance](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetMonthlyBalanceCache(year int, data []*response.CardResponseMonthBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetYearlyBalanceCache(year int) ([]*response.CardResponseYearlyBalance, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalance, year)

	result, found := GetFromCache[[]*response.CardResponseYearlyBalance](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetYearlyBalanceCache(year int, data []*response.CardResponseYearlyBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalance, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetMonthlyTopupAmountCache(year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetMonthlyTopupAmountCache(year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetYearlyTopupAmountCache(year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetYearlyTopupAmountCache(year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetMonthlyWithdrawAmountCache(year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetMonthlyWithdrawAmountCache(year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetYearlyWithdrawAmountCache(year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetYearlyWithdrawAmountCache(year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetMonthlyTransactionAmountCache(year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetMonthlyTransactionAmountCache(year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetYearlyTransactionAmountCache(year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetYearlyTransactionAmountCache(year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetMonthlyTransferAmountSenderCache(year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetMonthlyTransferAmountSenderCache(year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetYearlyTransferAmountSenderCache(year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetYearlyTransferAmountSenderCache(year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetMonthlyTransferAmountReceiverCache(year int) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetMonthlyTransferAmountReceiverCache(year int, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

func (c *cardStatisticCache) GetYearlyTransferAmountReceiverCache(year int) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticCache) SetYearlyTransferAmountReceiverCache(year int, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)
	SetToCache(c.store, key, &data, ttlStatistic)
}

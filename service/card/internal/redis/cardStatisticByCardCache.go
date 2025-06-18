package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	expirationCardStatistic = 10 * time.Minute

	cacheKeyMonthlyBalanceByCard  = "stat:monthly_balance_by_card:card_number:%s:year:%d"
	cacheKeyYearlyBalanceByCard   = "stat:yearly_balance_by_card:card_number:%s:year:%d"
	cacheKeyMonthlyTopupByCard    = "stat:monthly_topup_by_card:card_number:%s:year:%d"
	cacheKeyYearlyTopupByCard     = "stat:yearly_topup_by_card:card_number:%s:year:%d"
	cacheKeyMonthlyWithdrawByCard = "stat:monthly_withdraw_by_card:card_number:%s:year:%d"
	cacheKeyYearlyWithdrawByCard  = "stat:yearly_withdraw_by_card:card_number:%s:year:%d"
	cacheKeyMonthlyTxnByCard      = "stat:monthly_txn_by_card:card_number:%s:year:%d"
	cacheKeyYearlyTxnByCard       = "stat:yearly_txn_by_card:card_number:%s:year:%d"
	cacheKeyMonthlySenderByCard   = "stat:monthly_sender_by_card:card_number:%s:year:%d"
	cacheKeyYearlySenderByCard    = "stat:yearly_sender_by_card:card_number:%s:year:%d"
	cacheKeyMonthlyReceiverByCard = "stat:monthly_receiver_by_card:card_number:%s:year:%d"
	cacheKeyYearlyReceiverByCard  = "stat:yearly_receiver_by_card:card_number:%s:year:%d"
)

type cardStatisticByNumberCache struct {
	store *CacheStore
}

func NewCardStatisticByNumberCache(store *CacheStore) *cardStatisticByNumberCache {
	return &cardStatisticByNumberCache{store: store}
}

func (c *cardStatisticByNumberCache) GetMonthlyBalanceCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalanceByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseMonthBalance](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetYearlyBalanceCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalanceByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseYearlyBalance](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetMonthlyTopupAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetYearlyTopupAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetMonthlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetYearlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetMonthlyTransactionAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetYearlyTransactionAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetMonthlyTransferBySenderCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetYearlyTransferBySenderCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetMonthlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseMonthAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) GetYearlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.CardResponseYearAmount](c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatisticByNumberCache) SetMonthlyBalanceCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalanceByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetYearlyBalanceCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearlyBalance) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalanceByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetMonthlyTopupAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetYearlyTopupAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetMonthlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetYearlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetMonthlyTransactionAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetYearlyTransactionAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetMonthlyTransferBySenderCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetYearlyTransferBySenderCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetMonthlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatisticByNumberCache) SetYearlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)
	SetToCache(c.store, key, &data, expirationCardStatistic)
}

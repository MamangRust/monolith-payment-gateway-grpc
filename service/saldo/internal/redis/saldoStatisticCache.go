package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	saldoMonthTotalBalanceCacheKey = "saldo:month_total_balance:month:%d:year:%d"
	saldoYearTotalBalanceCacheKey  = "saldo:year_total_balance:year:%d"
	saldoYearlyBalanceCacheKey     = "saldo:yearly_balance:year:%d"
	saldoMonthBalanceCacheKey      = "saldo:month_balance:year:%d"
)

type saldoStatisticCache struct {
	store *CacheStore
}

func NewSaldoStatisticCache(store *CacheStore) *saldoStatisticCache {
	return &saldoStatisticCache{store: store}
}

func (c *saldoStatisticCache) GetMonthlyTotalSaldoBalanceCache(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, bool) {
	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	result, found := GetFromCache[[]*response.SaldoMonthTotalBalanceResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatisticCache) SetMonthlyTotalSaldoCache(req *requests.MonthTotalSaldoBalance, data []*response.SaldoMonthTotalBalanceResponse) {
	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *saldoStatisticCache) GetYearTotalSaldoBalanceCache(year int) ([]*response.SaldoYearTotalBalanceResponse, bool) {
	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	result, found := GetFromCache[[]*response.SaldoYearTotalBalanceResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatisticCache) SetYearTotalSaldoBalanceCache(year int, data []*response.SaldoYearTotalBalanceResponse) {
	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *saldoStatisticCache) GetMonthlySaldoBalanceCache(year int) ([]*response.SaldoMonthBalanceResponse, bool) {
	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	result, found := GetFromCache[[]*response.SaldoMonthBalanceResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatisticCache) SetMonthlySaldoBalanceCache(year int, data []*response.SaldoMonthBalanceResponse) {
	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

func (c *saldoStatisticCache) GetYearlySaldoBalanceCache(year int) ([]*response.SaldoYearBalanceResponse, bool) {
	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	result, found := GetFromCache[[]*response.SaldoYearBalanceResponse](c.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatisticCache) SetYearlySaldoBalanceCache(year int, data []*response.SaldoYearBalanceResponse) {
	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	SetToCache(c.store, key, &data, ttlDefault)
}

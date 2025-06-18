package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	monthTopupStatusSuccessByCardCacheKey = "topup:month:status:success:card_number:%s:month:%d:year:%d"
	yearTopupStatusSuccessByCardCacheKey  = "topup:year:status:success:card_number:%s:year:%d"
	monthTopupStatusFailedByCardCacheKey  = "topup:month:status:failed:card_number:%s:month:%d:year:%d"
	yearTopupStatusFailedByCardCacheKey   = "topup:year:status:failed:card_number:%s:year:%d"

	monthTopupAmountByCardCacheKey = "topup:month:amount:card_number:%s:year:%d"
	yearTopupAmountByCardCacheKey  = "topup:year:amount:card_number:%s:year:%d"

	monthTopupMethodByCardCacheKey = "topup:month:method:card_number:%s:year:%d"
	yearTopupMethodByCardCacheKey  = "topup:year:method:card_number:%s:year:%d"
)

type topupStatisticByCardCache struct {
	store *CacheStore
}

func NewTopupStatisticByCardCache(store *CacheStore) *topupStatisticByCardCache {
	return &topupStatisticByCardCache{
		store: store,
	}
}

func (s *topupStatisticByCardCache) GetMonthTopupStatusSuccessByCardNumberCache(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TopupResponseMonthStatusSuccess](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetMonthTopupStatusSuccessByCardNumberCache(req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusSuccess) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetYearlyTopupStatusSuccessByCardNumberCache(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TopupResponseYearStatusSuccess](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetYearlyTopupStatusSuccessByCardNumberCache(req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusSuccess) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetMonthTopupStatusFailedByCardNumberCache(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TopupResponseMonthStatusFailed](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetMonthTopupStatusFailedByCardNumberCache(req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusFailed) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetYearlyTopupStatusFailedByCardNumberCache(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TopupResponseYearStatusFailed](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetYearlyTopupStatusFailedByCardNumberCache(req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusFailed) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetMonthlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TopupMonthAmountResponse](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetMonthlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupMonthAmountResponse) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetYearlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TopupYearlyAmountResponse](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetYearlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupYearlyAmountResponse) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetMonthlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TopupMonthMethodResponse](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetMonthlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupMonthMethodResponse) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *topupStatisticByCardCache) GetYearlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := GetFromCache[[]*response.TopupYearlyMethodResponse](s.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (s *topupStatisticByCardCache) SetYearlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupYearlyMethodResponse) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	transferMonthTransferStatusSuccessKey = "transfer:month:transfer_status:success:month:%d:year:%d"
	transferMonthTransferStatusFailedKey  = "transfer:month:transfer_status:failed:month:%d:year:%d"

	transferYearTransferStatusSuccessKey = "transfer:year:transfer_status:success:year:%d"
	transferYearTransferStatusFailedKey  = "transfer:year:transfer_status:failed:year:%d"

	transferMonthTransferAmountKey = "transfer:month:transfer_amount:year:%d"
	transferYearTransferAmountKey  = "transfer:year:transfer_amount:year:%d"
)

type transferStatisticCache struct {
	store *CacheStore
}

func NewTransferStatisticCache(store *CacheStore) *transferStatisticCache {
	return &transferStatisticCache{store: store}
}

func (t *transferStatisticCache) GetCachedMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransferResponseMonthStatusSuccess](t.store, key)

	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transferStatisticCache) SetCachedMonthTransferStatusSuccess(req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusSuccess) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticCache) GetCachedYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)

	result, found := GetFromCache[[]*response.TransferResponseYearStatusSuccess](t.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}
func (t *transferStatisticCache) SetCachedYearlyTransferStatusSuccess(year int, data []*response.TransferResponseYearStatusSuccess) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticCache) GetCachedMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransferResponseMonthStatusFailed](t.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}
func (t *transferStatisticCache) SetCachedMonthTransferStatusFailed(req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusFailed) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticCache) GetCachedYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	result, found := GetFromCache[[]*response.TransferResponseYearStatusFailed](t.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}
func (t *transferStatisticCache) SetCachedYearlyTransferStatusFailed(year int, data []*response.TransferResponseYearStatusFailed) {
	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticCache) GetCachedMonthTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	result, found := GetFromCache[[]*response.TransferMonthAmountResponse](t.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}
func (t *transferStatisticCache) SetCachedMonthTransferAmounts(year int, data []*response.TransferMonthAmountResponse) {
	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticCache) GetCachedYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, bool) {
	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	result, found := GetFromCache[[]*response.TransferYearAmountResponse](t.store, key)
	if !found {
		return nil, false
	}
	return *result, true
}

func (t *transferStatisticCache) SetCachedYearlyTransferAmounts(year int, data []*response.TransferYearAmountResponse) {
	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	SetToCache(t.store, key, &data, ttlDefault)
}

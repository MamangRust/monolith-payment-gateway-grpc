package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

const (
	transferMonthTransferStatusSuccessByCardKey = "transfer:month:transfer_status:success:card_number:%s:month:%d:year:%d"
	transferMonthTransferStatusFailedByCardKey  = "transfer:month:transfer_status:failed:card_number:%s:month:%d:year:%d"

	transferYearTransferStatusSuccessByCardKey = "transfer:year:transfer_status:success:card_number:%s:year:%d"
	transferYearTransferStatusFailedByCardKey  = "transfer:year:transfer_status:failed:card_number:%s:year:%d"

	transferMonthTransferAmountByCardKey = "transfer:month:transfer_amount:card_number:%s:year:%d"

	transferYearTransferAmountByCardKey = "transfer:year:transfer_amount:card_number:%s:year:%d"
)

type transferStatisticByCardCache struct {
	store *CacheStore
}

func NewTransferStatisticByCardCache(store *CacheStore) *transferStatisticByCardCache {
	return &transferStatisticByCardCache{store: store}
}

func (t *transferStatisticByCardCache) GetMonthTransferStatusSuccessByCard(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransferResponseMonthStatusSuccess](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}

func (t *transferStatisticByCardCache) SetMonthTransferStatusSuccessByCard(req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetYearlyTransferStatusSuccessByCard(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransferResponseYearStatusSuccess](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatisticByCardCache) SetYearlyTransferStatusSuccessByCard(req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetMonthTransferStatusFailedByCard(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := GetFromCache[[]*response.TransferResponseMonthStatusFailed](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatisticByCardCache) SetMonthTransferStatusFailedByCard(req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetYearlyTransferStatusFailedByCard(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransferResponseYearStatusFailed](t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transferStatisticByCardCache) SetYearlyTransferStatusFailedByCard(req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetMonthlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransferMonthAmountResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatisticByCardCache) SetMonthlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetMonthlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransferMonthAmountResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatisticByCardCache) SetMonthlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetYearlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool) {
	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransferYearAmountResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatisticByCardCache) SetYearlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

func (t *transferStatisticByCardCache) GetYearlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool) {
	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := GetFromCache[[]*response.TransferYearAmountResponse](t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transferStatisticByCardCache) SetYearlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	SetToCache(t.store, key, &data, ttlDefault)
}

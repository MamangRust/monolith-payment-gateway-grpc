package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransferQueryCache interface {
	GetCachedTransfersCache(req *requests.FindAllTranfers) ([]*response.TransferResponse, *int, bool)
	SetCachedTransfersCache(req *requests.FindAllTranfers, data []*response.TransferResponse, total *int)

	GetCachedTransferActiveCache(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, bool)
	SetCachedTransferActiveCache(req *requests.FindAllTranfers, data []*response.TransferResponseDeleteAt, total *int)

	GetCachedTransferTrashedCache(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, bool)
	SetCachedTransferTrashedCache(req *requests.FindAllTranfers, data []*response.TransferResponseDeleteAt, total *int)

	GetCachedTransferCache(id int) (*response.TransferResponse, bool)
	SetCachedTransferCache(data *response.TransferResponse)

	GetCachedTransferByFrom(from string) ([]*response.TransferResponse, bool)
	SetCachedTransferByFrom(from string, data []*response.TransferResponse)

	GetCachedTransferByTo(to string) ([]*response.TransferResponse, bool)
	SetCachedTransferByTo(to string, data []*response.TransferResponse)
}

type TransferStatisticCache interface {
	GetCachedMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, bool)
	SetCachedMonthTransferStatusSuccess(req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusSuccess)

	GetCachedYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, bool)
	SetCachedYearlyTransferStatusSuccess(year int, data []*response.TransferResponseYearStatusSuccess)

	GetCachedMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, bool)
	SetCachedMonthTransferStatusFailed(req *requests.MonthStatusTransfer, data []*response.TransferResponseMonthStatusFailed)

	GetCachedYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, bool)
	SetCachedYearlyTransferStatusFailed(year int, data []*response.TransferResponseYearStatusFailed)

	GetCachedMonthTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, bool)
	SetCachedMonthTransferAmounts(year int, data []*response.TransferMonthAmountResponse)

	GetCachedYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, bool)
	SetCachedYearlyTransferAmounts(year int, data []*response.TransferYearAmountResponse)
}

type TransferStatisticByCardCache interface {
	GetMonthTransferStatusSuccessByCard(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, bool)
	SetMonthTransferStatusSuccessByCard(req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusSuccess)

	GetYearlyTransferStatusSuccessByCard(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, bool)
	SetYearlyTransferStatusSuccessByCard(req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusSuccess)

	GetMonthTransferStatusFailedByCard(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, bool)
	SetMonthTransferStatusFailedByCard(req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusFailed)

	GetYearlyTransferStatusFailedByCard(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, bool)
	SetYearlyTransferStatusFailedByCard(req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusFailed)

	GetMonthlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool)
	SetMonthlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse)

	GetMonthlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool)
	SetMonthlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse)

	GetYearlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool)
	SetYearlyTransferAmountsBySenderCard(req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse)

	GetYearlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool)
	SetYearlyTransferAmountsByReceiverCard(req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse)
}

type TransferCommandCache interface {
	DeleteTransferCache(id int)
}

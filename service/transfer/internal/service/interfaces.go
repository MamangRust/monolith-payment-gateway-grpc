package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransferQueryService interface {
	FindAll(req *requests.FindAllTranfers) ([]*response.TransferResponse, *int, *response.ErrorResponse)
	FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)
	FindTransferByTransferFrom(transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse)
	FindTransferByTransferTo(transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse)
}

type TransferStatisticsService interface {
	FindMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)
	FindMonthlyTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

type TransferStatisticByCardService interface {
	FindMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

type TransferCommandService interface {
	CreateTransaction(request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse)
	UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse)
	TrashedTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse)
	RestoreTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse)
	DeleteTransferPermanent(transfer_id int) (bool, *response.ErrorResponse)

	RestoreAllTransfer() (bool, *response.ErrorResponse)
	DeleteAllTransferPermanent() (bool, *response.ErrorResponse)
}

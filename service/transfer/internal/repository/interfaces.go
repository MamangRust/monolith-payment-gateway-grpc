package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoRepository interface {
	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
	UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
}

type CardRepository interface {
	FindUserCardByCardNumber(card_number string) (*record.CardEmailRecord, error)
	FindCardByCardNumber(card_number string) (*record.CardRecord, error)
}

type TransferQueryRepository interface {
	FindAll(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error)
	FindByActive(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error)
	FindById(id int) (*record.TransferRecord, error)
	FindByTrashed(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error)
	FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error)
	FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error)
}

type TransferStatisticRepository interface {
	GetMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error)
	GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error)
	GetMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error)
	GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error)

	GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error)
}

type TransferStatisticByCardRepository interface {
	GetMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error)
	GetYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error)
	GetMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error)
	GetYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error)

	GetMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error)
	GetMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error)
}

type TransferCommandRepository interface {
	CreateTransfer(request *requests.CreateTransferRequest) (*record.TransferRecord, error)
	UpdateTransfer(request *requests.UpdateTransferRequest) (*record.TransferRecord, error)
	UpdateTransferAmount(request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error)
	UpdateTransferStatus(request *requests.UpdateTransferStatus) (*record.TransferRecord, error)

	TrashedTransfer(transfer_id int) (*record.TransferRecord, error)
	RestoreTransfer(transfer_id int) (*record.TransferRecord, error)
	DeleteTransferPermanent(topup_id int) (bool, error)

	RestoreAllTransfer() (bool, error)
	DeleteAllTransferPermanent() (bool, error)
}

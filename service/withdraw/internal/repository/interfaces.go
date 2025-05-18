package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoRepository interface {
	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
	UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
	UpdateSaldoWithdraw(request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error)
}

type WithdrawQueryRepository interface {
	FindAll(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)
	FindByActive(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)
	FindByTrashed(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)
	FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*record.WithdrawRecord, *int, error)
	FindById(id int) (*record.WithdrawRecord, error)
}

type WithdrawStatisticRepository interface {
	GetMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusSuccess, error)
	GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error)
	GetMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusFailed, error)
	GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error)

	GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error)
	GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error)
}

type WithdrawStatisticByCardRepository interface {
	GetMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusSuccess, error)
	GetYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusSuccess, error)
	GetMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusFailed, error)
	GetYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusFailed, error)
	GetMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawMonthlyAmount, error)
	GetYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error)
}

type WithdrawCommandRepository interface {
	CreateWithdraw(request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error)
	UpdateWithdraw(request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error)
	UpdateWithdrawStatus(request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error)

	TrashedWithdraw(WithdrawID int) (*record.WithdrawRecord, error)
	RestoreWithdraw(WithdrawID int) (*record.WithdrawRecord, error)
	DeleteWithdrawPermanent(WithdrawID int) (bool, error)

	RestoreAllWithdraw() (bool, error)
	DeleteAllWithdrawPermanent() (bool, error)
}

type CardRepository interface {
	FindUserCardByCardNumber(card_number string) (*record.CardEmailRecord, error)
}

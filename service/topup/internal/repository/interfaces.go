package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoRepository interface {
	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
	UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
}

type TopupQueryRepository interface {
	FindAllTopups(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)
	FindByActive(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)
	FindByTrashed(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)
	FindAllTopupByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*record.TopupRecord, *int, error)

	FindById(topup_id int) (*record.TopupRecord, error)
}

type TopupStatisticRepository interface {
	GetMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusSuccess, error)
	GetYearlyTopupStatusSuccess(year int) ([]*record.TopupRecordYearStatusSuccess, error)
	GetMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusFailed, error)
	GetYearlyTopupStatusFailed(year int) ([]*record.TopupRecordYearStatusFailed, error)

	GetMonthlyTopupMethods(year int) ([]*record.TopupMonthMethod, error)
	GetYearlyTopupMethods(year int) ([]*record.TopupYearlyMethod, error)
	GetMonthlyTopupAmounts(year int) ([]*record.TopupMonthAmount, error)
	GetYearlyTopupAmounts(year int) ([]*record.TopupYearlyAmount, error)
}

type TopupStatisticByCardRepository interface {
	GetMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusSuccess, error)
	GetYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusSuccess, error)

	GetMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusFailed, error)
	GetYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusFailed, error)

	GetMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthMethod, error)
	GetYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyMethod, error)
	GetMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthAmount, error)
	GetYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyAmount, error)
}

type TopupCommandRepository interface {
	CreateTopup(request *requests.CreateTopupRequest) (*record.TopupRecord, error)
	UpdateTopup(request *requests.UpdateTopupRequest) (*record.TopupRecord, error)
	UpdateTopupAmount(request *requests.UpdateTopupAmount) (*record.TopupRecord, error)
	UpdateTopupStatus(request *requests.UpdateTopupStatus) (*record.TopupRecord, error)
	TrashedTopup(topup_id int) (*record.TopupRecord, error)
	RestoreTopup(topup_id int) (*record.TopupRecord, error)
	DeleteTopupPermanent(topup_id int) (bool, error)
	RestoreAllTopup() (bool, error)
	DeleteAllTopupPermanent() (bool, error)
}

type CardRepository interface {
	FindUserCardByCardNumber(card_number string) (*record.CardEmailRecord, error)
	FindCardByCardNumber(card_number string) (*record.CardRecord, error)
	UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error)
}

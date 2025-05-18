package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoQueryRepository interface {
	FindAllSaldos(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)
	FindByActive(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)
	FindByTrashed(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)
	FindById(saldo_id int) (*record.SaldoRecord, error)
	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
}

type SaldoCommandRepository interface {
	CreateSaldo(request *requests.CreateSaldoRequest) (*record.SaldoRecord, error)
	UpdateSaldo(request *requests.UpdateSaldoRequest) (*record.SaldoRecord, error)
	UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
	UpdateSaldoWithdraw(request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error)
	TrashedSaldo(saldoID int) (*record.SaldoRecord, error)
	RestoreSaldo(saldoID int) (*record.SaldoRecord, error)
	DeleteSaldoPermanent(saldo_id int) (bool, error)

	RestoreAllSaldo() (bool, error)
	DeleteAllSaldoPermanent() (bool, error)
}

type SaldoStatisticsRepository interface {
	GetMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*record.SaldoMonthTotalBalance, error)
	GetYearTotalSaldoBalance(year int) ([]*record.SaldoYearTotalBalance, error)
	GetMonthlySaldoBalances(year int) ([]*record.SaldoMonthSaldoBalance, error)
	GetYearlySaldoBalances(year int) ([]*record.SaldoYearSaldoBalance, error)
}

type CardRepository interface {
	FindCardByCardNumber(card_number string) (*record.CardRecord, error)
}

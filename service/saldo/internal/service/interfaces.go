package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type SaldoQueryService interface {
	FindAll(req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, *response.ErrorResponse)
	FindById(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	FindByCardNumber(card_number string) (*response.SaldoResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
}

type SaldoStatisticService interface {
	FindMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse)
	FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse)
	FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse)
	FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse)
}

type SaldoCommandService interface {
	CreateSaldo(request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)
	UpdateSaldo(request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)
	TrashSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	RestoreSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	DeleteSaldoPermanent(saldo_id int) (bool, *response.ErrorResponse)

	RestoreAllSaldo() (bool, *response.ErrorResponse)
	DeleteAllSaldoPermanent() (bool, *response.ErrorResponse)
}

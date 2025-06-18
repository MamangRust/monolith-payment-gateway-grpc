package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type SaldoQueryCache interface {
	GetCachedSaldos(req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, bool)
	SetCachedSaldos(req *requests.FindAllSaldos, data []*response.SaldoResponse, totalRecords *int)

	GetCachedSaldoById(saldo_id int) (*response.SaldoResponse, bool)
	SetCachedSaldoById(saldo_id int, data *response.SaldoResponse)

	GetCachedSaldoByCardNumber(card_number string) (*response.SaldoResponse, bool)
	SetCachedSaldoByCardNumber(card_number string, data *response.SaldoResponse)

	GetCachedSaldoByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool)
	SetCachedSaldoByActive(req *requests.FindAllSaldos, data []*response.SaldoResponseDeleteAt, totalRecords *int)

	GetCachedSaldoByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool)
	SetCachedSaldoByTrashed(req *requests.FindAllSaldos, data []*response.SaldoResponseDeleteAt, totalRecords *int)
}

type SaldoStatisticCache interface {
	GetMonthlyTotalSaldoBalanceCache(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, bool)
	SetMonthlyTotalSaldoCache(req *requests.MonthTotalSaldoBalance, data []*response.SaldoMonthTotalBalanceResponse)

	GetYearTotalSaldoBalanceCache(year int) ([]*response.SaldoYearTotalBalanceResponse, bool)
	SetYearTotalSaldoBalanceCache(year int, data []*response.SaldoYearTotalBalanceResponse)

	GetMonthlySaldoBalanceCache(year int) ([]*response.SaldoMonthBalanceResponse, bool)
	SetMonthlySaldoBalanceCache(year int, data []*response.SaldoMonthBalanceResponse)

	GetYearlySaldoBalanceCache(year int) ([]*response.SaldoYearBalanceResponse, bool)
	SetYearlySaldoBalanceCache(year int, data []*response.SaldoYearBalanceResponse)
}

type SaldoCommandCache interface {
	DeleteSaldoCache(saldo_id int)
}

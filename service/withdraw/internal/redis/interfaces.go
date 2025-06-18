package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type WithdrawQueryCache interface {
	GetCachedWithdrawsCache(req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, bool)
	SetCachedWithdrawsCache(req *requests.FindAllWithdraws, data []*response.WithdrawResponse, total *int)

	GetCachedWithdrawByCardCache(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, bool)
	SetCachedWithdrawByCardCache(req *requests.FindAllWithdrawCardNumber, data []*response.WithdrawResponse, total *int)

	GetCachedWithdrawActiveCache(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool)
	SetCachedWithdrawActiveCache(req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int)

	GetCachedWithdrawTrashedCache(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool)
	SetCachedWithdrawTrashedCache(req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int)

	GetCachedWithdrawCache(id int) (*response.WithdrawResponse, bool)
	SetCachedWithdrawCache(data *response.WithdrawResponse)
}

type WithdrawCommandCache interface {
	DeleteCachedWithdrawCache(id int)
}

type WithdrawStatisticCache interface {
	GetCachedMonthWithdrawStatusSuccessCache(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, bool)
	SetCachedMonthWithdrawStatusSuccessCache(req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusSuccess)

	GetCachedYearlyWithdrawStatusSuccessCache(year int) ([]*response.WithdrawResponseYearStatusSuccess, bool)
	SetCachedYearlyWithdrawStatusSuccessCache(year int, data []*response.WithdrawResponseYearStatusSuccess)

	GetCachedMonthWithdrawStatusFailedCache(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, bool)
	SetCachedMonthWithdrawStatusFailedCache(req *requests.MonthStatusWithdraw, data []*response.WithdrawResponseMonthStatusFailed)

	GetCachedYearlyWithdrawStatusFailedCache(year int) ([]*response.WithdrawResponseYearStatusFailed, bool)
	SetCachedYearlyWithdrawStatusFailedCache(year int, data []*response.WithdrawResponseYearStatusFailed)

	GetCachedMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, bool)
	SetCachedMonthlyWithdraws(year int, data []*response.WithdrawMonthlyAmountResponse)

	GetCachedYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, bool)
	SetCachedYearlyWithdraws(year int, data []*response.WithdrawYearlyAmountResponse)
}

type WithdrawStasticByCardCache interface {
	GetCachedMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, bool)
	SetCachedMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusSuccess)

	GetCachedYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, bool)
	SetCachedYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusSuccess)

	GetCachedMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, bool)
	SetCachedMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusFailed)

	GetCachedYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, bool)
	SetCachedYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusFailed)

	GetCachedMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, bool)
	SetCachedMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber, data []*response.WithdrawMonthlyAmountResponse)

	GetCachedYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, bool)
	SetCachedYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber, data []*response.WithdrawYearlyAmountResponse)
}

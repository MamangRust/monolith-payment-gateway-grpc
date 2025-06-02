package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TopupQueryCache interface {
	GetCachedTopupsCache(req *requests.FindAllTopups) ([]*response.TopupResponse, *int, bool)
	SetCachedTopupsCache(req *requests.FindAllTopups, data []*response.TopupResponse, total *int)

	GetCacheTopupByCardCache(req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, bool)
	SetCacheTopupByCardCache(req *requests.FindAllTopupsByCardNumber, data []*response.TopupResponse, total *int)

	GetCachedTopupActiveCache(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool)
	SetCachedTopupActiveCache(req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int)

	GetCachedTopupTrashedCache(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool)
	SetCachedTopupTrashedCache(req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int)

	GetCachedTopupCache(id int) *response.TopupResponse
	SetCachedTopupCache(data *response.TopupResponse)
}

type TopupStatisticCache interface {
	GetMonthTopupStatusSuccessCache(req *requests.MonthTopupStatus) []*response.TopupResponseMonthStatusSuccess
	SetMonthTopupStatusSuccessCache(req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusSuccess)

	GetYearlyTopupStatusSuccessCache(year int) []*response.TopupResponseYearStatusSuccess
	SetYearlyTopupStatusSuccessCache(year int, data []*response.TopupResponseYearStatusSuccess)

	GetMonthTopupStatusFailedCache(req *requests.MonthTopupStatus) []*response.TopupResponseMonthStatusFailed
	SetMonthTopupStatusFailedCache(req *requests.MonthTopupStatus, data []*response.TopupResponseMonthStatusFailed)

	GetYearlyTopupStatusFailedCache(year int) []*response.TopupResponseYearStatusFailed
	SetYearlyTopupStatusFailedCache(year int, data []*response.TopupResponseYearStatusFailed)

	GetMonthlyTopupMethodsCache(year int) []*response.TopupMonthMethodResponse
	SetMonthlyTopupMethodsCache(year int, data []*response.TopupMonthMethodResponse)

	GetYearlyTopupMethodsCache(year int) []*response.TopupYearlyMethodResponse
	SetYearlyTopupMethodsCache(year int, data []*response.TopupYearlyMethodResponse)

	GetMonthlyTopupAmountsCache(year int) []*response.TopupMonthAmountResponse
	SetMonthlyTopupAmountsCache(year int, data []*response.TopupMonthAmountResponse)

	GetYearlyTopupAmountsCache(year int) []*response.TopupYearlyAmountResponse
	SetYearlyTopupAmountsCache(year int, data []*response.TopupYearlyAmountResponse)
}

type TopupStatisticByCardCache interface {
	GetMonthTopupStatusSuccessByCardNumberCache(req *requests.MonthTopupStatusCardNumber) []*response.TopupResponseMonthStatusSuccess
	SetMonthTopupStatusSuccessByCardNumberCache(req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusSuccess)

	GetYearlyTopupStatusSuccessByCardNumberCache(req *requests.YearTopupStatusCardNumber) []*response.TopupResponseYearStatusSuccess
	SetYearlyTopupStatusSuccessByCardNumberCache(req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusSuccess)

	GetMonthTopupStatusFailedByCardNumberCache(req *requests.MonthTopupStatusCardNumber) []*response.TopupResponseMonthStatusFailed
	SetMonthTopupStatusFailedByCardNumberCache(req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusFailed)

	GetYearlyTopupStatusFailedByCardNumberCache(req *requests.YearTopupStatusCardNumber) []*response.TopupResponseYearStatusFailed
	SetYearlyTopupStatusFailedByCardNumberCache(req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusFailed)

	GetMonthlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod) []*response.TopupMonthMethodResponse
	SetMonthlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupMonthMethodResponse)

	GetYearlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod) []*response.TopupYearlyMethodResponse
	SetYearlyTopupMethodsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupYearlyMethodResponse)

	GetMonthlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod) []*response.TopupMonthAmountResponse
	SetMonthlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupMonthAmountResponse)

	GetYearlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod) []*response.TopupYearlyAmountResponse
	SetYearlyTopupAmountsByCardNumberCache(req *requests.YearMonthMethod, data []*response.TopupYearlyAmountResponse)
}

type TopupCommandCache interface {
	DeleteCachedTopupCache(id int)
}

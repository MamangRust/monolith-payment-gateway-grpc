package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TopupQueryService interface {
	FindAll(req *requests.FindAllTopups) ([]*response.TopupResponse, *int, *response.ErrorResponse)
	FindAllByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, *response.ErrorResponse)
	FindById(topupID int) (*response.TopupResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
}

type TopupStatisticService interface {
	FindMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTopupStatusSuccess(year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTopupStatusFailed(year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTopupMethods(year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	FindYearlyTopupMethods(year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyTopupAmounts(year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTopupAmounts(year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

type TopupStatisticByCardService interface {
	FindMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	FindYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

type TopupCommandService interface {
	CreateTopup(request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)
	UpdateTopup(request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)
	TrashedTopup(topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	RestoreTopup(topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	DeleteTopupPermanent(topup_id int) (bool, *response.ErrorResponse)

	RestoreAllTopup() (bool, *response.ErrorResponse)
	DeleteAllTopupPermanent() (bool, *response.ErrorResponse)
}

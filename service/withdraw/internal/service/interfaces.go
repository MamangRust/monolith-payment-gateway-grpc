package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type WithdrawQueryService interface {
	FindAll(req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)
	FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)
	FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)
}

type WithdrawStatisticService interface {
	FindMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyWithdrawStatusSuccess(year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyWithdrawStatusFailed(year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)
	FindMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	FindYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
}

type WithdrawStatisticByCardService interface {
	FindMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	FindYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
}

type WithdrawCommandService interface {
	Create(request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)
	Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)
	TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)
	RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)
	DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse)

	RestoreAllWithdraw() (bool, *response.ErrorResponse)
	DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse)
}

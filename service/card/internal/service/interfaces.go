package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type CardQueryService interface {
	FindAll(req *requests.FindAllCards) ([]*response.CardResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(cardID int) (*response.CardResponse, *response.ErrorResponse)
	FindByUserID(userID int) (*response.CardResponse, *response.ErrorResponse)
	FindByCardNumber(cardNumber string) (*response.CardResponse, *response.ErrorResponse)
}

type CardDashboardService interface {
	DashboardCard() (*response.DashboardCard, *response.ErrorResponse)
	DashboardCardCardNumber(cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse)
}

type CardStatisticService interface {
	FindMonthlyBalance(year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)
	FindYearlyBalance(year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)

	FindMonthlyTopupAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTopupAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyWithdrawAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyWithdrawAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyTransactionAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransactionAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyTransferAmountSender(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountSender(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyTransferAmountReceiver(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountReceiver(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardStatisticByNumberService interface {
	FindMonthlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)
	FindYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)

	FindMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardCommandService interface {
	CreateCard(request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse)
	UpdateCard(request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse)
	TrashedCard(cardId int) (*response.CardResponse, *response.ErrorResponse)
	RestoreCard(cardId int) (*response.CardResponse, *response.ErrorResponse)
	DeleteCardPermanent(cardId int) (bool, *response.ErrorResponse)
	RestoreAllCard() (bool, *response.ErrorResponse)
	DeleteAllCardPermanent() (bool, *response.ErrorResponse)
}

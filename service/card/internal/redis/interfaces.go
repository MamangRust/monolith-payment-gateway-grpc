package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type CardQueryCache interface {
	GetByIdCache(cardID int) (*response.CardResponse, bool)
	GetByUserIDCache(userID int) (*response.CardResponse, bool)
	GetByCardNumberCache(cardNumber string) (*response.CardResponse, bool)
	GetFindAllCache(req *requests.FindAllCards) ([]*response.CardResponse, *int, bool)
	GetByActiveCache(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool)
	GetByTrashedCache(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool)
	SetByIdCache(cardID int, data *response.CardResponse)
	SetByUserIDCache(userID int, data *response.CardResponse)
	SetByCardNumberCache(cardNumber string, data *response.CardResponse)
	SetFindAllCache(req *requests.FindAllCards, data []*response.CardResponse, totalRecords *int)
	SetByActiveCache(req *requests.FindAllCards, data []*response.CardResponseDeleteAt, totalRecords *int)
	SetByTrashedCache(req *requests.FindAllCards, data []*response.CardResponseDeleteAt, totalRecords *int)
	DeleteByIdCache(cardID int)
	DeleteByUserIDCache(userID int)
	DeleteByCardNumberCache(cardNumber string)
}

type CardDashboardCache interface {
	GetDashboardCardCache() (*response.DashboardCard, bool)
	SetDashboardCardCache(data *response.DashboardCard)
	DeleteDashboardCardCache()
	GetDashboardCardCardNumberCache(cardNumber string) (*response.DashboardCardCardNumber, bool)
	SetDashboardCardCardNumberCache(cardNumber string, data *response.DashboardCardCardNumber)
	DeleteDashboardCardCardNumberCache(cardNumber string)
}

type CardCommandCache interface {
	DeleteCardCommandCache(id int)
}

type CardStatisticByNumberCache interface {
	GetMonthlyBalanceCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, bool)
	GetYearlyBalanceCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, bool)
	GetMonthlyTopupAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)
	GetYearlyTopupAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)
	GetMonthlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)
	GetYearlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)
	GetMonthlyTransactionAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)
	GetYearlyTransactionAmountCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)
	GetMonthlyTransferBySenderCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)
	GetYearlyTransferBySenderCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)
	GetMonthlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)
	GetYearlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)
	SetMonthlyBalanceCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthBalance)
	SetYearlyBalanceCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearlyBalance)
	SetMonthlyTopupAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)
	SetYearlyTopupAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
	SetMonthlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)
	SetYearlyWithdrawAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
	SetMonthlyTransactionAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)
	SetYearlyTransactionAmountCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
	SetMonthlyTransferBySenderCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)
	SetYearlyTransferBySenderCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
	SetMonthlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)
	SetYearlyTransferByReceiverCache(req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
}

type CardStatisticCache interface {
	GetMonthlyBalanceCache(year int) ([]*response.CardResponseMonthBalance, bool)
	SetMonthlyBalanceCache(year int, data []*response.CardResponseMonthBalance)
	GetYearlyBalanceCache(year int) ([]*response.CardResponseYearlyBalance, bool)
	SetYearlyBalanceCache(year int, data []*response.CardResponseYearlyBalance)
	GetMonthlyTopupAmountCache(year int) ([]*response.CardResponseMonthAmount, bool)
	SetMonthlyTopupAmountCache(year int, data []*response.CardResponseMonthAmount)
	GetYearlyTopupAmountCache(year int) ([]*response.CardResponseYearAmount, bool)
	SetYearlyTopupAmountCache(year int, data []*response.CardResponseYearAmount)
	GetMonthlyWithdrawAmountCache(year int) ([]*response.CardResponseMonthAmount, bool)
	SetMonthlyWithdrawAmountCache(year int, data []*response.CardResponseMonthAmount)
	SetYearlyWithdrawAmountCache(year int, data []*response.CardResponseYearAmount)
	GetMonthlyTransactionAmountCache(year int) ([]*response.CardResponseMonthAmount, bool)
	SetMonthlyTransactionAmountCache(year int, data []*response.CardResponseMonthAmount)
	GetYearlyTransactionAmountCache(year int) ([]*response.CardResponseYearAmount, bool)
	SetYearlyTransactionAmountCache(year int, data []*response.CardResponseYearAmount)
	GetMonthlyTransferAmountSenderCache(year int) ([]*response.CardResponseMonthAmount, bool)
	SetMonthlyTransferAmountSenderCache(year int, data []*response.CardResponseMonthAmount)
	GetYearlyTransferAmountSenderCache(year int) ([]*response.
		CardResponseYearAmount, bool)
	SetYearlyTransferAmountSenderCache(year int, data []*response.CardResponseYearAmount)
	GetMonthlyTransferAmountReceiverCache(year int) ([]*response.CardResponseMonthAmount, bool)
	SetMonthlyTransferAmountReceiverCache(year int, data []*response.CardResponseMonthAmount)
	GetYearlyTransferAmountReceiverCache(year int) ([]*response.CardResponseYearAmount, bool)
	SetYearlyTransferAmountReceiverCache(year int, data []*response.CardResponseYearAmount)
}

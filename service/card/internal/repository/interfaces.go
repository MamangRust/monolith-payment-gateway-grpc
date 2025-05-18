package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type CardCommandRepository interface {
	CreateCard(request *requests.CreateCardRequest) (*record.CardRecord, error)
	UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error)
	TrashedCard(cardId int) (*record.CardRecord, error)
	RestoreCard(cardId int) (*record.CardRecord, error)
	DeleteCardPermanent(card_id int) (bool, error)
	RestoreAllCard() (bool, error)
	DeleteAllCardPermanent() (bool, error)
}

type CardQueryRepository interface {
	FindAllCards(req *requests.FindAllCards) ([]*record.CardRecord, *int, error)
	FindByActive(req *requests.FindAllCards) ([]*record.CardRecord, *int, error)
	FindByTrashed(req *requests.FindAllCards) ([]*record.CardRecord, *int, error)
	FindById(card_id int) (*record.CardRecord, error)
	FindCardByUserId(user_id int) (*record.CardRecord, error)
	FindCardByCardNumber(card_number string) (*record.CardRecord, error)
}

type CardDashboardRepository interface {
	GetTotalBalances() (*int64, error)
	GetTotalTopAmount() (*int64, error)
	GetTotalWithdrawAmount() (*int64, error)
	GetTotalTransactionAmount() (*int64, error)
	GetTotalTransferAmount() (*int64, error)

	GetTotalBalanceByCardNumber(cardNumber string) (*int64, error)
	GetTotalTopupAmountByCardNumber(cardNumber string) (*int64, error)
	GetTotalWithdrawAmountByCardNumber(cardNumber string) (*int64, error)
	GetTotalTransactionAmountByCardNumber(cardNumber string) (*int64, error)
	GetTotalTransferAmountBySender(senderCardNumber string) (*int64, error)
	GetTotalTransferAmountByReceiver(receiverCardNumber string) (*int64, error)
}

type CardStatisticRepository interface {
	GetMonthlyBalance(year int) ([]*record.CardMonthBalance, error)
	GetYearlyBalance(year int) ([]*record.CardYearlyBalance, error)
	GetMonthlyTopupAmount(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTopupAmount(year int) ([]*record.CardYearAmount, error)
	GetMonthlyWithdrawAmount(year int) ([]*record.CardMonthAmount, error)
	GetYearlyWithdrawAmount(year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransactionAmount(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransactionAmount(year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountSender(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountSender(year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountReceiver(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountReceiver(year int) ([]*record.CardYearAmount, error)
}

type CardStatisticByCardRepository interface {
	GetMonthlyBalancesByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthBalance, error)
	GetYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearlyBalance, error)
	GetMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
}

type UserRepository interface {
	FindById(user_id int) (*record.UserRecord, error)
}

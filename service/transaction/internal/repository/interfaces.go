package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantRepository interface {
	FindByApiKey(api_key string) (*record.MerchantRecord, error)
}

type SaldoRepository interface {
	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
	UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
}

type CardRepository interface {
	FindCardByUserId(user_id int) (*record.CardRecord, error)
	FindUserCardByCardNumber(card_number string) (*record.CardEmailRecord, error)
	FindCardByCardNumber(card_number string) (*record.CardRecord, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)
	FindByActive(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)
	FindByTrashed(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)
	FindAllTransactionByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*record.TransactionRecord, *int, error)
	FindById(transaction_id int) (*record.TransactionRecord, error)
	FindTransactionByMerchantId(merchant_id int) ([]*record.TransactionRecord, error)
}

type TransactionStatisticsRepository interface {
	GetMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusSuccess, error)
	GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error)
	GetMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusFailed, error)
	GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error)

	GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error)
	GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error)
	GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error)
	GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error)
}

type TransactionStatisticByCardRepository interface {
	GetMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusSuccess, error)
	GetYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusSuccess, error)
	GetMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusFailed, error)
	GetYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusFailed, error)

	GetMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthMethod, error)
	GetYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error)
	GetMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error)
	GetYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error)
}

type TransactionCommandRepository interface {
	CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransactionStatus(request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error)
	TrashedTransaction(transaction_id int) (*record.TransactionRecord, error)
	RestoreTransaction(topup_id int) (*record.TransactionRecord, error)
	DeleteTransactionPermanent(topup_id int) (bool, error)

	RestoreAllTransaction() (bool, error)
	DeleteAllTransactionPermanent() (bool, error)
}

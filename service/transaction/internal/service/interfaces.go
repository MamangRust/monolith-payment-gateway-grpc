package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransactionQueryService interface {
	FindAll(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindAllByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindTransactionByMerchantId(merchant_id int) ([]*response.TransactionResponse, *response.ErrorResponse)
}

type TransactionStatisticService interface {
	FindMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransactionStatusSuccess(year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransactionStatusFailed(year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyPaymentMethods(year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	FindYearlyPaymentMethods(year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
	FindMonthlyAmounts(year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	FindYearlyAmounts(year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

type TransactionsStatistcByCardService interface {
	FindMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	FindYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
	FindMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	FindYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

type TransactionCommandService interface {
	Create(apiKey string, request *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	Update(apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	TrashedTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse)
	RestoreTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse)
	DeleteTransactionPermanent(transaction_id int) (bool, *response.ErrorResponse)

	RestoreAllTransaction() (bool, *response.ErrorResponse)
	DeleteAllTransactionPermanent() (bool, *response.ErrorResponse)
}

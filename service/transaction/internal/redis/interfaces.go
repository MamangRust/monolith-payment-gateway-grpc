package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransactinQueryCache interface {
	GetCachedTransactionsCache(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, bool)
	SetCachedTransactionsCache(req *requests.FindAllTransactions, data []*response.TransactionResponse, total *int)
	GetCachedTransactionByCardNumberCache(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, bool)
	SetCachedTransactionByCardNumberCache(req *requests.FindAllTransactionCardNumber, data []*response.TransactionResponse, total *int)

	GetCachedTransactionActiveCache(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool)
	SetCachedTransactionActiveCache(req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, total *int)
	GetCachedTransactionTrashedCache(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool)
	SetCachedTransactionTrashedCache(req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, total *int)

	GetCachedTransactionByMerchantIdCache(merchant_id int) []*response.TransactionResponse
	SetCachedTransactionByMerchantIdCache(merchant_id int, data []*response.TransactionResponse)

	GetCachedTransactionCache(id int) *response.TransactionResponse
	SetCachedTransactionCache(data *response.TransactionResponse)
}

type TransactonStatistcCache interface {
	GetMonthTransactonStatusSuccessCache(req *requests.MonthStatusTransaction) []*response.TransactionResponseMonthStatusSuccess
	SetMonthTransactonStatusSuccessCache(req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusSuccess)
	GetYearTransactonStatusSuccessCache(year int) []*response.TransactionResponseYearStatusSuccess
	SetYearTransactonStatusSuccessCache(year int, data []*response.TransactionResponseYearStatusSuccess)

	GetMonthTransactonStatusFailedCache(req *requests.MonthStatusTransaction) []*response.TransactionResponseMonthStatusFailed
	SetMonthTransactonStatusFailedCache(req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusFailed)
	GetYearTransactonStatusFailedCache(year int) []*response.TransactionResponseYearStatusFailed
	SetYearTransactonStatusFailedCache(year int, data []*response.TransactionResponseYearStatusFailed)

	GetMonthlyPaymentMethodsCache(year int) []*response.TransactionMonthMethodResponse
	SetMonthlyPaymentMethodsCache(year int, data []*response.TransactionMonthMethodResponse)
	GetYearlyPaymentMethodsCache(year int) []*response.TransactionYearMethodResponse
	SetYearlyPaymentMethodsCache(year int, data []*response.TransactionYearMethodResponse)

	GetMonthlyAmountsCache(year int) []*response.TransactionMonthAmountResponse
	SetMonthlyAmountsCache(year int, data []*response.TransactionMonthAmountResponse)

	GetYearlyAmountsCache(year int) []*response.TransactionYearlyAmountResponse
	SetYearlyAmountsCache(year int, data []*response.TransactionYearlyAmountResponse)
}

type TransactionStatisticByCardCache interface {
	GetMonthTransactionStatusSuccessByCardCache(req *requests.MonthStatusTransactionCardNumber) []*response.TransactionResponseMonthStatusSuccess
	SetMonthTransactionStatusSuccessByCardCache(req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusSuccess)

	GetYearTransactionStatusSuccessByCardCache(req *requests.YearStatusTransactionCardNumber) []*response.TransactionResponseYearStatusSuccess
	SetYearTransactionStatusSuccessByCardCache(req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusSuccess)

	GetMonthTransactionStatusFailedByCardCache(req *requests.MonthStatusTransactionCardNumber) []*response.TransactionResponseMonthStatusFailed
	SetMonthTransactionStatusFailedByCardCache(req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusFailed)

	GetYearTransactionStatusFailedByCardCache(req *requests.YearStatusTransactionCardNumber) []*response.TransactionResponseYearStatusFailed
	SetYearTransactionStatusFailedByCardCache(req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusFailed)

	GetMonthlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionMonthMethodResponse
	SetMonthlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthMethodResponse)

	GetYearlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionYearMethodResponse
	SetYearlyPaymentMethodsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionYearMethodResponse)

	GetMonthlyAmountsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionMonthAmountResponse
	SetMonthlyAmountsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthAmountResponse)

	GetYearlyAmountsByCardCache(req *requests.MonthYearPaymentMethod) []*response.TransactionYearlyAmountResponse
	SetYearlyAmountsByCardCache(req *requests.MonthYearPaymentMethod, data []*response.TransactionYearlyAmountResponse)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(id int)
}

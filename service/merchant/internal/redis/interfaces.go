package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantQueryCache interface {
	GetCachedMerchants(req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, bool)
	SetCachedMerchants(req *requests.FindAllMerchants, data []*response.MerchantResponse, total *int)
	GetCachedMerchantActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool)
	SetCachedMerchantActive(req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int)
	GetCachedMerchantTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool)
	SetCachedMerchantTrashed(req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int)
	GetCachedMerchant(id int) (*response.MerchantResponse, bool)
	SetCachedMerchant(data *response.MerchantResponse)
	GetCachedMerchantsByUserId(id int) ([]*response.MerchantResponse, bool)
	SetCachedMerchantsByUserId(userId int, data []*response.MerchantResponse)
	GetCachedMerchantByApiKey(apiKey string) *response.MerchantResponse
	SetCachedMerchantByApiKey(apiKey string, data *response.MerchantResponse)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool)
	SetCachedMerchantDocuments(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int)
	SetCachedMerchantDocumentsTrashed(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)
	GetCachedMerchantDocumentsActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)
	SetCachedMerchantDocumentsActive(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)
	GetCachedMerchantDocumentsTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)
	GetCachedMerchantDocument(id int) (*response.MerchantDocumentResponse, bool)
	SetCachedMerchantDocument(id int, data *response.MerchantDocumentResponse)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(id int)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(id int)
}

type MerchantTransactionCache interface {
	SetCacheAllMerchantTransactions(req *requests.FindAllMerchantTransactions, data []*response.MerchantTransactionResponse, total *int)
	GetCacheAllMerchantTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, bool)
	SetCacheMerchantTransactions(req *requests.FindAllMerchantTransactionsById, data []*response.MerchantTransactionResponse, total *int)
	GetCacheMerchantTransactions(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, bool)
	SetCacheMerchantTransactionApikey(req *requests.FindAllMerchantTransactionsByApiKey, data []*response.MerchantTransactionResponse, total *int)
	GetCacheMerchantTransactionApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, bool)
}

type MerchantStatisticCache interface {
	GetMonthlyPaymentMethodsMerchantCache(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, bool)
	SetMonthlyPaymentMethodsMerchantCache(year int, data []*response.MerchantResponseMonthlyPaymentMethod)

	GetYearlyPaymentMethodMerchantCache(year int) ([]*response.MerchantResponseYearlyPaymentMethod, bool)
	SetYearlyPaymentMethodMerchantCache(year int, data []*response.MerchantResponseYearlyPaymentMethod)

	GetMonthlyAmountMerchantCache(year int) ([]*response.MerchantResponseMonthlyAmount, bool)
	SetMonthlyAmountMerchantCache(year int, data []*response.MerchantResponseMonthlyAmount)

	GetYearlyAmountMerchantCache(year int) ([]*response.MerchantResponseYearlyAmount, bool)
	SetYearlyAmountMerchantCache(year int, data []*response.MerchantResponseYearlyAmount)

	GetMonthlyTotalAmountMerchantCache(year int) ([]*response.MerchantResponseMonthlyTotalAmount, bool)
	SetMonthlyTotalAmountMerchantCache(year int, data []*response.MerchantResponseMonthlyTotalAmount)

	GetYearlyTotalAmountMerchantCache(year int) ([]*response.MerchantResponseYearlyTotalAmount, bool)
	SetYearlyTotalAmountMerchantCache(year int, data []*response.MerchantResponseYearlyTotalAmount)
}

type MerchantStatisticByMerchantCache interface {
	SetMonthlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseMonthlyPaymentMethod)
	GetMonthlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, bool)

	SetYearlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseYearlyPaymentMethod)
	GetYearlyPaymentMethodByMerchantsCache(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, bool)

	SetMonthlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseMonthlyAmount)
	GetMonthlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, bool)

	SetYearlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseYearlyAmount)
	GetYearlyAmountByMerchantsCache(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, bool)

	SetMonthlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseMonthlyTotalAmount)
	GetMonthlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, bool)

	SetYearlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseYearlyTotalAmount)
	GetYearlyTotalAmountByMerchantsCache(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, bool)
}

type MerchantStatisticByApikeyCache interface {
	SetMonthlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseMonthlyPaymentMethod)
	GetMonthlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, bool)

	SetYearlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseYearlyPaymentMethod)
	GetYearlyPaymentMethodByApikeysCache(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, bool)

	SetMonthlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseMonthlyAmount)
	GetMonthlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, bool)

	SetYearlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseYearlyAmount)
	GetYearlyAmountByApikeysCache(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, bool)

	SetMonthlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseMonthlyTotalAmount)
	GetMonthlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, bool)

	SetYearlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseYearlyTotalAmount)
	GetYearlyTotalAmountByApikeysCache(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, bool)
}

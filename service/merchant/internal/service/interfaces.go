package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantQueryService interface {
	FindAll(req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByApiKey(api_key string) (*response.MerchantResponse, *response.ErrorResponse)
	FindByMerchantUserId(user_id int) ([]*response.MerchantResponse, *response.ErrorResponse)
}

type MerchantDocumentQueryService interface {
	FindAll(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse)

	FindById(document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
}

type MerchantTransactionService interface {
	FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
	FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
	FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
}

type MerchantCommandService interface {
	CreateMerchant(request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse)
	TrashedMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)
	RestoreMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)
	DeleteMerchantPermanent(merchant_id int) (bool, *response.ErrorResponse)

	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	TrashedMerchantDocument(document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	RestoreMerchantDocument(document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	DeleteMerchantDocumentPermanent(document_id int) (bool, *response.ErrorResponse)

	RestoreAllMerchantDocument() (bool, *response.ErrorResponse)
	DeleteAllMerchantDocumentPermanent() (bool, *response.ErrorResponse)
}

type MerchantStatisticService interface {
	FindMonthlyPaymentMethodsMerchant(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodMerchant(year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountMerchant(year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountMerchant(year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)

	FindMonthlyTotalAmountMerchant(year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}

type MerchantStatisticByMerchantService interface {
	FindMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
	FindMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}

type MerchantStatisticByApikeyService interface {
	FindMonthlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
	FindMonthlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}

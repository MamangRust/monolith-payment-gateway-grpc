package repository

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type UserRepository interface {
	FindById(user_id int) (*record.UserRecord, error)
}

type MerchantQueryRepository interface {
	FindAllMerchants(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
	FindByActive(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
	FindById(merchant_id int) (*record.MerchantRecord, error)
	FindByApiKey(api_key string) (*record.MerchantRecord, error)
	FindByName(name string) (*record.MerchantRecord, error)
	FindByMerchantUserId(user_id int) ([]*record.MerchantRecord, error)
}

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
	FindById(id int) (*record.MerchantDocumentRecord, error)

	FindByActive(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
}

type MerchantTransactionRepository interface {
	FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error)
	FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error)
	FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error)
	TrashedMerchant(merchantId int) (*record.MerchantRecord, error)
	RestoreMerchant(merchant_id int) (*record.MerchantRecord, error)
	DeleteMerchantPermanent(merchant_id int) (bool, error)
	RestoreAllMerchant() (bool, error)
	DeleteAllMerchantPermanent() (bool, error)
}

type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)
	UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)
	UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error)
	TrashedMerchantDocument(merchant_document_id int) (*record.MerchantDocumentRecord, error)
	RestoreMerchantDocument(merchant_document_id int) (*record.MerchantDocumentRecord, error)
	DeleteMerchantDocumentPermanent(merchant_document_id int) (bool, error)
	RestoreAllMerchantDocument() (bool, error)
	DeleteAllMerchantDocumentPermanent() (bool, error)
}

type MerchantStatisticRepository interface {
	GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error)
}

type MerchantStatisticByMerchantRepository interface {
	GetMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error)
}

type MerchantStatisticByApiKeyRepository interface {
	GetMonthlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error)
}

// type MerchantRepository interface {
// 	FindAllMerchants(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
// 	FindByActive(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
// 	FindByTrashed(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)

// 	FindById(merchant_id int) (*record.MerchantRecord, error)
// 	GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error)
// 	GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error)

// 	FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error)
// 	GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error)
// 	GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error)
// 	GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error)
// 	GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error)

// 	FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error)
// 	GetMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error)
// 	GetYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error)
// 	GetMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error)
// 	GetYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error)
// 	GetMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error)
// 	GetYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error)

// 	FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error)
// 	GetMonthlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error)
// 	GetYearlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error)
// 	GetMonthlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error)
// 	GetYearlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error)
// 	GetMonthlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error)
// 	GetYearlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error)

// 	FindByApiKey(api_key string) (*record.MerchantRecord, error)
// 	FindByName(name string) (*record.MerchantRecord, error)
// 	FindByMerchantUserId(user_id int) ([]*record.MerchantRecord, error)

// 	CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
// 	UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
// 	UpdateMerchantStatus(request *requests.UpdateMerchantStatus) (*record.MerchantRecord, error)

// 	TrashedMerchant(merchantId int) (*record.MerchantRecord, error)
// 	RestoreMerchant(merchant_id int) (*record.MerchantRecord, error)
// 	DeleteMerchantPermanent(merchant_id int) (bool, error)

// 	RestoreAllMerchant() (bool, error)
// 	DeleteAllMerchantPermanent() (bool, error)
// }

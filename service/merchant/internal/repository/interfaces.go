package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// UserRepository is an interface that defines methods for interacting with user data.
type UserRepository interface {
	// FindById retrieves a user record by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to retrieve.
	//
	// Returns:
	//   - *record.UserRecord: The user record if found.
	//   - error: An error if any occurred during the query.
	FindByUserId(ctx context.Context, user_id int) (*record.UserRecord, error)
}

type MerchantQueryRepository interface {
	// FindAllMerchants retrieves all merchants based on the given request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination, filters, and search query.
	//
	// Returns:
	//   - []*record.MerchantRecord: The list of merchant records.
	//   - *int: The total count of records.
	//   - error: An error if any occurred during the query.
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)

	// FindByActive retrieves only active merchants based on the request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination and filters.
	//
	// Returns:
	//   - []*record.MerchantRecord: The list of active merchant records.
	//   - *int: The total count of active records.
	//   - error: An error if any occurred during the query.
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)

	// FindByTrashed retrieves only trashed (soft-deleted) merchants based on the request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination and filters.
	//
	// Returns:
	//   - []*record.MerchantRecord: The list of trashed merchant records.
	//   - *int: The total count of trashed records.
	//   - error: An error if any occurred during the query.
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)

	// FindById retrieves a single merchant by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to retrieve.
	//
	// Returns:
	//   - *record.MerchantRecord: The merchant record if found.
	//   - error: An error if any occurred during the query.
	FindById(ctx context.Context, merchant_id int) (*record.MerchantRecord, error)

	// FindByApiKey retrieves a merchant by its API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - api_key: The API key of the merchant to retrieve.
	//
	// Returns:
	//   - *record.MerchantRecord: The merchant record if found.
	//   - error: An error if any occurred during the query.
	FindByApiKey(ctx context.Context, api_key string) (*record.MerchantRecord, error)

	// FindByName retrieves a merchant by its name.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - name: The name of the merchant to retrieve.
	//
	// Returns:
	//   - *record.MerchantRecord: The merchant record if found.
	//   - error: An error if any occurred during the query.
	FindByName(ctx context.Context, name string) (*record.MerchantRecord, error)

	// FindByMerchantUserId retrieves all merchants linked to a specific user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The user ID associated with the merchant(s).
	//
	// Returns:
	//   - []*record.MerchantRecord: The list of merchant records owned by the user.
	//   - error: An error if any occurred during the query.
	FindByMerchantUserId(ctx context.Context, user_id int) ([]*record.MerchantRecord, error)
}

type MerchantDocumentQueryRepository interface {
	// FindAllDocuments retrieves all merchant documents based on the request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination, filters, and search query.
	//
	// Returns:
	//   - []*record.MerchantDocumentRecord: The list of document records.
	//   - *int: The total count of records.
	//   - error: An error if any occurred during the query.
	FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)

	// FindById retrieves a single merchant document by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the document to retrieve.
	//
	// Returns:
	//   - *record.MerchantDocumentRecord: The document record if found.
	//   - error: An error if any occurred during the query.
	FindByIdDocument(ctx context.Context, id int) (*record.MerchantDocumentRecord, error)

	// FindByActive retrieves only active (non-deleted) merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters.
	//
	// Returns:
	//   - []*record.MerchantDocumentRecord: The list of active document records.
	//   - *int: The total count of active records.
	//   - error: An error if any occurred during the query.
	FindByActiveDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)

	// FindByTrashed retrieves only trashed (soft-deleted) merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters.
	//
	// Returns:
	//   - []*record.MerchantDocumentRecord: The list of trashed document records.
	//   - *int: The total count of trashed records.
	//   - error: An error if any occurred during the query.
	FindByTrashedDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
}

type MerchantTransactionRepository interface {
	// FindAllTransactions retrieves all merchant transactions based on the request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters and pagination.
	//
	// Returns:
	//   - []*record.MerchantTransactionsRecord: The list of transaction records.
	//   - *int: The total count of records.
	//   - error: An error if any occurred during the query.
	FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error)

	// FindAllTransactionsByMerchant retrieves transactions for a specific merchant by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing merchant ID and optional filters.
	//
	// Returns:
	//   - []*record.MerchantTransactionsRecord: The list of transaction records.
	//   - *int: The total count of records.
	//   - error: An error if any occurred during the query.
	FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error)

	// FindAllTransactionsByApikey retrieves transactions based on merchant API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing API key and optional filters.
	//
	// Returns:
	//   - []*record.MerchantTransactionsRecord: The list of transaction records.
	//   - *int: The total count of records.
	//   - error: An error if any occurred during the query.
	FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error)
}

type MerchantCommandRepository interface {
	// CreateMerchant inserts a new merchant record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request object containing merchant creation data.
	//
	// Returns:
	//   - *record.MerchantRecord: The created merchant record.
	//   - error: An error if any occurred during the insert.
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)

	// UpdateMerchant updates an existing merchant's details.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request object containing updated merchant data.
	//
	// Returns:
	//   - *record.MerchantRecord: The updated merchant record.
	//   - error: An error if any occurred during the update.
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)

	// UpdateMerchantStatus updates only the status field of a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request object containing merchant ID and new status.
	//
	// Returns:
	//   - *record.MerchantRecord: The updated merchant record.
	//   - error: An error if any occurred during the update.
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error)

	// TrashedMerchant performs a soft delete (trash) on a merchant record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchantId: The ID of the merchant to trash.
	//
	// Returns:
	//   - *record.MerchantRecord: The trashed merchant record.
	//   - error: An error if any occurred during the update.
	TrashedMerchant(ctx context.Context, merchantId int) (*record.MerchantRecord, error)

	// RestoreMerchant restores a soft-deleted (trashed) merchant record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to restore.
	//
	// Returns:
	//   - *record.MerchantRecord: The restored merchant record.
	//   - error: An error if any occurred during the restore.
	RestoreMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error)

	// DeleteMerchantPermanent permanently deletes a merchant record from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to delete permanently.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: An error if any occurred during deletion.
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)

	// RestoreAllMerchant restores all soft-deleted merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the restore was successful.
	//   - error: An error if any occurred during the operation.
	RestoreAllMerchant(ctx context.Context) (bool, error)

	// DeleteAllMerchantPermanent permanently deletes all trashed merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: An error if any occurred during the operation.
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentCommandRepository interface {
	// CreateMerchantDocument inserts a new merchant document record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request object containing document creation data.
	//
	// Returns:
	//   - *record.MerchantDocumentRecord: The created merchant document record.
	//   - error: An error if any occurred during the insert.
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)

	// UpdateMerchantDocument updates an existing merchant document's fields.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request object containing updated document data.
	//
	// Returns:
	//   - *record.MerchantDocumentRecord: The updated merchant document record.
	//   - error: An error if any occurred during the update.
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)

	// UpdateMerchantDocumentStatus updates only the status field of a merchant document.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request object containing document ID and new status.
	//
	// Returns:
	//   - *record.MerchantDocumentRecord: The updated merchant document record.
	//   - error: An error if any occurred during the update.
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error)

	// TrashedMerchantDocument performs a soft delete (trash) on a merchant document.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_document_id: The ID of the document to trash.
	//
	// Returns:
	//   - *record.MerchantDocumentRecord: The trashed merchant document record.
	//   - error: An error if any occurred during the operation.
	TrashedMerchantDocument(ctx context.Context, merchant_document_id int) (*record.MerchantDocumentRecord, error)

	// RestoreMerchantDocument restores a soft-deleted (trashed) merchant document.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_document_id: The ID of the document to restore.
	//
	// Returns:
	//   - *record.MerchantDocumentRecord: The restored document record.
	//   - error: An error if any occurred during the restore.
	RestoreMerchantDocument(ctx context.Context, merchant_document_id int) (*record.MerchantDocumentRecord, error)

	// DeleteMerchantDocumentPermanent permanently deletes a merchant document from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_document_id: The ID of the document to permanently delete.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: An error if any occurred during the deletion.
	DeleteMerchantDocumentPermanent(ctx context.Context, merchant_document_id int) (bool, error)

	// RestoreAllMerchantDocument restores all soft-deleted merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the restore operation was successful.
	//   - error: An error if any occurred during the operation.
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)

	// DeleteAllMerchantDocumentPermanent permanently deletes all trashed merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the deletion operation was successful.
	//   - error: An error if any occurred during the operation.
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}

type MerchantStatisticRepository interface {
}

type MerchantStatisticByMerchantRepository interface {
}

type MerchantStatisticByApiKeyRepository interface {
}

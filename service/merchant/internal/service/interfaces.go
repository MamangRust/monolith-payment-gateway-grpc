package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// MerchantQueryService defines methods for querying merchant data such as all merchants,
// specific merchant by ID, API key, user ID, and soft-deleted (trashed) records.
type MerchantQueryService interface {
	// FindAll retrieves a list of merchants with pagination and optional filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters for filtering and pagination.
	//
	// Returns:
	//   - []*response.MerchantResponse: A list of merchant data.
	//   - *int: The total count of matched merchants.
	//   - *response.ErrorResponse: An error if the operation fails.
	FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, *response.ErrorResponse)

	// FindById retrieves a merchant by its unique merchant ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to be retrieved.
	//
	// Returns:
	//   - *response.MerchantResponse: The merchant data.
	//   - *response.ErrorResponse: An error if the merchant is not found or retrieval fails.
	FindById(ctx context.Context, merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)

	// FindByActive retrieves all active (non-deleted) merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters for filtering and pagination.
	//
	// Returns:
	//   - []*response.MerchantResponseDeleteAt: A list of active merchant records.
	//   - *int: The total count of matched active merchants.
	//   - *response.ErrorResponse: An error if the operation fails.
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves all soft-deleted (trashed) merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters for filtering and pagination.
	//
	// Returns:
	//   - []*response.MerchantResponseDeleteAt: A list of trashed merchant records.
	//   - *int: The total count of matched trashed merchants.
	//   - *response.ErrorResponse: An error if the operation fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByApiKey retrieves a merchant based on the provided API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - api_key: The API key associated with the merchant.
	//
	// Returns:
	//   - *response.MerchantResponse: The matched merchant.
	//   - *response.ErrorResponse: An error if not found or retrieval fails.
	FindByApiKey(ctx context.Context, api_key string) (*response.MerchantResponse, *response.ErrorResponse)

	// FindByMerchantUserId retrieves merchants associated with a given user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user who owns the merchants.
	//
	// Returns:
	//   - []*response.MerchantResponse: A list of merchants owned by the user.
	//   - *response.ErrorResponse: An error if retrieval fails.
	FindByMerchantUserId(ctx context.Context, user_id int) ([]*response.MerchantResponse, *response.ErrorResponse)
}

// MerchantDocumentQueryService defines methods for querying merchant documents,
// including active, trashed, and all documents, as well as finding by document ID.
type MerchantDocumentQueryService interface {
	// FindAll retrieves all merchant documents with optional filtering and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters including filters and pagination options.
	//
	// Returns:
	//   - []*response.MerchantDocumentResponse: A list of merchant documents.
	//   - *int: The total count of matched documents.
	//   - *response.ErrorResponse: An error if retrieval fails.
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)

	// FindByActive retrieves all active (non-deleted) merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters including filters and pagination options.
	//
	// Returns:
	//   - []*response.MerchantDocumentResponse: A list of active merchant documents.
	//   - *int: The total count of matched active documents.
	//   - *response.ErrorResponse: An error if retrieval fails.
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)

	// FindByTrashed retrieves all trashed (soft-deleted) merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters including filters and pagination options.
	//
	// Returns:
	//   - []*response.MerchantDocumentResponseDeleteAt: A list of trashed merchant documents.
	//   - *int: The total count of matched trashed documents.
	//   - *response.ErrorResponse: An error if retrieval fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse)

	// FindById retrieves a merchant document by its unique document ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - document_id: The ID of the document to retrieve.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The document data.
	//   - *response.ErrorResponse: An error if the document is not found or retrieval fails.
	FindById(ctx context.Context, document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
}

// MerchantTransactionService defines methods for retrieving merchant transaction data,
// including all transactions, and filtered by merchant ID or API key.
type MerchantTransactionService interface {
	// FindAllTransactions retrieves all merchant transactions with optional filters and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters for filters, search, and pagination.
	//
	// Returns:
	//   - []*response.MerchantTransactionResponse: A list of transaction records.
	//   - *int: The total number of matched transactions.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)

	// FindAllTransactionsByMerchant retrieves all transactions for a specific merchant by merchant ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing the merchant ID and other filters.
	//
	// Returns:
	//   - []*response.MerchantTransactionResponse: A list of transaction records.
	//   - *int: The total number of matched transactions.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)

	// FindAllTransactionsByApikey retrieves all transactions for a merchant identified by an API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing the API key and other filters.
	//
	// Returns:
	//   - []*response.MerchantTransactionResponse: A list of transaction records.
	//   - *int: The total number of matched transactions.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
}

// MerchantCommandService defines command operations related to merchants,
// including create, update, soft-delete, restore, and permanent deletion.
type MerchantCommandService interface {
	// CreateMerchant creates a new merchant with the provided request data.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The merchant creation request payload.
	//
	// Returns:
	//   - *response.MerchantResponse: The created merchant's data.
	//   - *response.ErrorResponse: An error if creation fails.
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)

	// UpdateMerchant updates an existing merchant's data.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing updated merchant data.
	//
	// Returns:
	//   - *response.MerchantResponse: The updated merchant's data.
	//   - *response.ErrorResponse: An error if update fails.
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)

	// UpdateMerchantStatus updates the status (active/inactive) of a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing status update information.
	//
	// Returns:
	//   - *response.MerchantResponse: The updated merchant data.
	//   - *response.ErrorResponse: An error if status update fails.
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse)

	// TrashedMerchant soft-deletes a merchant by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to be soft-deleted.
	//
	// Returns:
	//   - *response.MerchantResponse: The trashed merchant data.
	//   - *response.ErrorResponse: An error if the operation fails.
	TrashedMerchant(ctx context.Context, merchant_id int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse)

	// RestoreMerchant restores a soft-deleted merchant by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to restore.
	//
	// Returns:
	//   - *response.MerchantResponse: The restored merchant data.
	//   - *response.ErrorResponse: An error if restoration fails.
	RestoreMerchant(ctx context.Context, merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)

	// DeleteMerchantPermanent permanently deletes a merchant by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant to delete permanently.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - *response.ErrorResponse: An error if deletion fails.
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, *response.ErrorResponse)

	// RestoreAllMerchant restores all soft-deleted merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the restoration was successful.
	//   - *response.ErrorResponse: An error if restoration fails.
	RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllMerchantPermanent permanently deletes all soft-deleted merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - *response.ErrorResponse: An error if deletion fails.
	DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

// MerchantDocumentCommandService defines command operations for managing merchant documents,
// including create, update, soft delete, restore, and permanent deletion.
type MerchantDocumentCommandService interface {
	// CreateMerchantDocument creates a new merchant document.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The merchant document creation request payload.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The created merchant document.
	//   - *response.ErrorResponse: An error if creation fails.
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	// UpdateMerchantDocument updates an existing merchant document.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The update request containing new document data.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The updated merchant document.
	//   - *response.ErrorResponse: An error if update fails.
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	// UpdateMerchantDocumentStatus updates the status (e.g., verified, rejected) of a merchant document.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing status update data.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The updated merchant document.
	//   - *response.ErrorResponse: An error if status update fails.
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	// TrashedMerchantDocument soft-deletes a merchant document by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - document_id: The ID of the document to be soft-deleted.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The trashed document.
	//   - *response.ErrorResponse: An error if the operation fails.
	TrashedMerchantDocument(ctx context.Context, document_id int) (*response.MerchantDocumentResponseDeleteAt, *response.ErrorResponse)

	// RestoreMerchantDocument restores a soft-deleted merchant document by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - document_id: The ID of the document to restore.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The restored document.
	//   - *response.ErrorResponse: An error if restoration fails.
	RestoreMerchantDocument(ctx context.Context, document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	// DeleteMerchantDocumentPermanent permanently deletes a merchant document by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - document_id: The ID of the document to delete.
	//
	// Returns:
	//   - bool: True if the deletion was successful.
	//   - *response.ErrorResponse: An error if the deletion fails.
	DeleteMerchantDocumentPermanent(ctx context.Context, document_id int) (bool, *response.ErrorResponse)

	// RestoreAllMerchantDocument restores all soft-deleted merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all documents were restored successfully.
	//   - *response.ErrorResponse: An error if restoration fails.
	RestoreAllMerchantDocument(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllMerchantDocumentPermanent permanently deletes all soft-deleted merchant documents.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all documents were deleted successfully.
	//   - *response.ErrorResponse: An error if deletion fails.
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

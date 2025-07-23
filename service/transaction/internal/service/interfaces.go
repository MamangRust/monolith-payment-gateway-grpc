package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// TransactionQueryService handles queries related to transactions.
type TransactionQueryService interface {
	// FindAll retrieves a paginated list of all transactions based on the given filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination info.
	//
	// Returns:
	//   - []*response.TransactionResponse: List of transactions.
	//   - *int: Total number of transactions.
	//   - *response.ErrorResponse: Error response if query fails.
	FindAll(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse)

	// FindAllByCardNumber retrieves all transactions associated with a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number and pagination.
	//
	// Returns:
	//   - []*response.TransactionResponse: List of transactions for the card number.
	//   - *int: Total number of transactions.
	//   - *response.ErrorResponse: Error response if query fails.
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse)

	// FindById retrieves a transaction by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transactionID: The ID of the transaction to retrieve.
	//
	// Returns:
	//   - *response.TransactionResponse: The transaction data.
	//   - *response.ErrorResponse: Error response if query fails.
	FindById(ctx context.Context, transactionID int) (*response.TransactionResponse, *response.ErrorResponse)

	// FindByActive retrieves active transactions with soft-delete not applied.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination info.
	//
	// Returns:
	//   - []*response.TransactionResponseDeleteAt: List of active transactions.
	//   - *int: Total number of active transactions.
	//   - *response.ErrorResponse: Error response if query fails.
	FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves transactions that have been soft-deleted.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination info.
	//
	// Returns:
	//   - []*response.TransactionResponseDeleteAt: List of trashed transactions.
	//   - *int: Total number of trashed transactions.
	//   - *response.ErrorResponse: Error response if query fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)

	// FindTransactionByMerchantId retrieves transactions by merchant ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant whose transactions are requested.
	//
	// Returns:
	//   - []*response.TransactionResponse: List of transactions for the merchant.
	//   - *response.ErrorResponse: Error response if query fails.
	FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*response.TransactionResponse, *response.ErrorResponse)
}

// TransactionCommandService handles command operations related to transactions.
type TransactionCommandService interface {
	// Create creates a new transaction based on the provided request and API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - apiKey: The API key for merchant authorization.
	//   - request: The transaction creation request payload.
	//
	// Returns:
	//   - *response.TransactionResponse: The created transaction response.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	Create(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)

	// Update updates an existing transaction with the given request and API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - apiKey: The API key for merchant authorization.
	//   - request: The transaction update request payload.
	//
	// Returns:
	//   - *response.TransactionResponse: The updated transaction response.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	Update(ctx context.Context, apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)

	// TrashedTransaction moves the transaction to the trash (soft delete).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transaction_id: The ID of the transaction to be trashed.
	//
	// Returns:
	//   - *response.TransactionResponse: The trashed transaction response.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	TrashedTransaction(ctx context.Context, transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)

	// RestoreTransaction restores a previously trashed transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transaction_id: The ID of the transaction to be restored.
	//
	// Returns:
	//   - *response.TransactionResponse: The restored transaction response.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	RestoreTransaction(ctx context.Context, transaction_id int) (*response.TransactionResponse, *response.ErrorResponse)

	// DeleteTransactionPermanent permanently deletes a transaction from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transaction_id: The ID of the transaction to delete permanently.
	//
	// Returns:
	//   - bool: Whether the operation was successful.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, *response.ErrorResponse)

	// RestoreAllTransaction restores all trashed transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the operation was successful.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	RestoreAllTransaction(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllTransactionPermanent permanently deletes all trashed transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the operation was successful.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	DeleteAllTransactionPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// MerchantRepository defines methods to interact with merchant records in the database.
type MerchantRepository interface {
	// FindByApiKey retrieves a merchant by its API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - api_key: The API key associated with the merchant.
	//
	// Returns:
	//   - *record.MerchantRecord: The merchant record if found.
	//   - error: Error if something went wrong during the query.
	FindByApiKey(ctx context.Context, api_key string) (*record.MerchantRecord, error)
}

// SaldoRepository defines methods to access and update saldo records in the database.
type SaldoRepository interface {
	// FindByCardNumber retrieves saldo information by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to lookup.
	//
	// Returns:
	//   - *record.SaldoRecord: The saldo record if found.
	//   - error: Error if something went wrong during the query.
	FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error)

	// UpdateSaldoBalance updates the saldo balance based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The update request containing new saldo balance data.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: Error if something went wrong during the update.
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
}

// CardRepository defines methods to interact with card records in the database.
type CardRepository interface {
	// FindCardByUserId retrieves a card associated with a specific user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The user ID to lookup.
	//
	// Returns:
	//   - *record.CardRecord: The card record if found.
	//   - error: Error if something went wrong during the query.
	FindCardByUserId(ctx context.Context, user_id int) (*record.CardRecord, error)

	// FindUserCardByCardNumber retrieves a user's card including email by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to lookup.
	//
	// Returns:
	//   - *record.CardEmailRecord: The card and user email record.
	//   - error: Error if something went wrong during the query.
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error)

	// FindCardByCardNumber retrieves a card by its card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to lookup.
	//
	// Returns:
	//   - *record.CardRecord: The card record if found.
	//   - error: Error if something went wrong during the query.
	FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error)
}

// TransactionQueryRepository defines query methods for transaction records in the database.
type TransactionQueryRepository interface {
	// FindAllTransactions retrieves all transactions with optional filtering and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The filter and pagination parameters.
	//
	// Returns:
	//   - []*record.TransactionRecord: A list of transaction records.
	//   - *int: The total count of transactions.
	//   - error: Error if something went wrong during the query.
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)

	// FindByActive retrieves all active (non-deleted) transactions with optional filtering and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The filter and pagination parameters.
	//
	// Returns:
	//   - []*record.TransactionRecord: A list of active transaction records.
	//   - *int: The total count of active transactions.
	//   - error: Error if something went wrong during the query.
	FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)

	// FindByTrashed retrieves all soft-deleted (trashed) transactions with optional filtering and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The filter and pagination parameters.
	//
	// Returns:
	//   - []*record.TransactionRecord: A list of trashed transaction records.
	//   - *int: The total count of trashed transactions.
	//   - error: Error if something went wrong during the query.
	FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)

	// FindAllTransactionByCardNumber retrieves all transactions by a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and optional filters.
	//
	// Returns:
	//   - []*record.TransactionRecord: A list of transactions associated with the card number.
	//   - *int: The total count of transactions.
	//   - error: Error if something went wrong during the query.
	FindAllTransactionByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*record.TransactionRecord, *int, error)

	// FindById retrieves a transaction by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transaction_id: The ID of the transaction.
	//
	// Returns:
	//   - *record.TransactionRecord: The transaction record if found.
	//   - error: Error if something went wrong during the query.
	FindById(ctx context.Context, transaction_id int) (*record.TransactionRecord, error)

	// FindTransactionByMerchantId retrieves all transactions associated with a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: The ID of the merchant.
	//
	// Returns:
	//   - []*record.TransactionRecord: A list of transactions for the given merchant.
	//   - error: Error if something went wrong during the query.
	FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*record.TransactionRecord, error)
}

// TransactionCommandRepository defines commands to manipulate transaction records in the data store.
type TransactionCommandRepository interface {
	// CreateTransaction creates a new transaction record in the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The data to create the transaction.
	//
	// Returns:
	//   - *record.TransactionRecord: The created transaction record.
	//   - error: Error if operation fails.
	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*record.TransactionRecord, error)

	// UpdateTransaction updates an existing transaction record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The data to update the transaction.
	//
	// Returns:
	//   - *record.TransactionRecord: The updated transaction record.
	//   - error: Error if operation fails.
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error)

	// UpdateTransactionStatus updates only the status of a transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The status update request.
	//
	// Returns:
	//   - *record.TransactionRecord: The transaction with updated status.
	//   - error: Error if operation fails.
	UpdateTransactionStatus(ctx context.Context, request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error)

	// TrashedTransaction marks a transaction as soft-deleted.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transaction_id: ID of the transaction to trash.
	//
	// Returns:
	//   - *record.TransactionRecord: The trashed transaction record.
	//   - error: Error if operation fails.
	TrashedTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error)

	// RestoreTransaction restores a soft-deleted transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: ID of the transaction to restore.
	//
	// Returns:
	//   - *record.TransactionRecord: The restored transaction record.
	//   - error: Error if operation fails.
	RestoreTransaction(ctx context.Context, topup_id int) (*record.TransactionRecord, error)

	// DeleteTransactionPermanent deletes a transaction permanently.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: ID of the transaction to delete permanently.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - error: Error if operation fails.
	DeleteTransactionPermanent(ctx context.Context, topup_id int) (bool, error)

	// RestoreAllTransaction restores all soft-deleted transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if restore operation was successful.
	//   - error: Error if operation fails.
	RestoreAllTransaction(ctx context.Context) (bool, error)

	// DeleteAllTransactionPermanent deletes all transactions permanently.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - error: Error if operation fails.
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}

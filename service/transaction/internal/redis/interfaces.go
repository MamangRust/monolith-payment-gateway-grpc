package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// TransactionQueryCache defines methods for caching transaction query results.
type TransactionQueryCache interface {
	// GetCachedTransactionsCache retrieves cached paginated transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Filter and pagination parameters.
	//
	// Returns:
	//   - []*response.TransactionResponse: List of transactions.
	//   - *int: Total number of matching records.
	//   - bool: Whether the cache was found.
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, bool)

	// SetCachedTransactionsCache stores paginated transaction results in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Original request used as cache key.
	//   - data: List of transaction responses to cache.
	//   - total: Total number of records for pagination.
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransactions, data []*response.TransactionResponse, total *int)

	// GetCachedTransactionByCardNumberCache retrieves cached transactions for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and filters.
	//
	// Returns:
	//   - []*response.TransactionResponse: Transactions matching the card number.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedTransactionByCardNumberCache(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, bool)

	// SetCachedTransactionByCardNumberCache caches transactions by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and filters.
	//   - data: Transactions to cache.
	//   - total: Total number of matching records.
	SetCachedTransactionByCardNumberCache(ctx context.Context, req *requests.FindAllTransactionCardNumber, data []*response.TransactionResponse, total *int)

	// GetCachedTransactionActiveCache retrieves cached active transactions (not soft-deleted).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Filter and pagination request.
	//
	// Returns:
	//   - []*response.TransactionResponseDeleteAt: List of active transactions.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool)

	// SetCachedTransactionActiveCache caches active transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Original filter request.
	//   - data: Transactions to cache.
	//   - total: Total number of records.
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, total *int)

	// GetCachedTransactionTrashedCache retrieves cached trashed (soft-deleted) transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Filter and pagination request.
	//
	// Returns:
	//   - []*response.TransactionResponseDeleteAt: List of trashed transactions.
	//   - *int: Total records.
	//   - bool: Whether the cache was found.
	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, bool)

	// SetCachedTransactionTrashedCache stores trashed transactions in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Original request with filters.
	//   - data: Trashed transactions to cache.
	//   - total: Total number of records.
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransactions, data []*response.TransactionResponseDeleteAt, total *int)

	// GetCachedTransactionByMerchantIdCache retrieves cached transactions for a merchant ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: ID of the merchant.
	//
	// Returns:
	//   - []*response.TransactionResponse: Transactions related to the merchant.
	//   - bool: Whether the cache was found.
	GetCachedTransactionByMerchantIdCache(ctx context.Context, merchant_id int) ([]*response.TransactionResponse, bool)

	// SetCachedTransactionByMerchantIdCache caches transactions by merchant ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - merchant_id: ID of the merchant.
	//   - data: Transactions to cache.
	SetCachedTransactionByMerchantIdCache(ctx context.Context, merchant_id int, data []*response.TransactionResponse)

	// GetCachedTransactionCache retrieves a transaction by ID from cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The transaction ID.
	//
	// Returns:
	//   - *response.TransactionResponse: The transaction response.
	//   - bool: Whether the cache was found.
	GetCachedTransactionCache(ctx context.Context, id int) (*response.TransactionResponse, bool)

	// SetCachedTransactionCache caches a transaction by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - data: The transaction response to cache.
	SetCachedTransactionCache(ctx context.Context, data *response.TransactionResponse)
}

type TransactionCommandCache interface {
	// DeleteTransactionCache removes a cached transaction entry by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The transaction ID whose cache entry should be deleted.
	DeleteTransactionCache(ctx context.Context, id int)
}

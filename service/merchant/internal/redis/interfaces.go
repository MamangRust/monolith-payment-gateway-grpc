package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// MerchantQueryCache is an interface for caching merchant queries
type MerchantQueryCache interface {
	// GetCachedMerchants retrieves a list of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the page, page size, and search string.
	//
	// Returns:
	//   - []*response.MerchantResponse: The list of merchants.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, bool)

	// SetCachedMerchants stores a list of merchants into cache based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of merchants to cache.
	//   - total: The total records to cache.
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*response.MerchantResponse, total *int)

	// GetCachedMerchantActive retrieves a list of active merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the page, page size, and search string.
	//
	// Returns:
	//   - []*response.MerchantResponseDeleteAt: The list of active merchants.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool)

	// SetCachedMerchantActive stores a list of active merchants into cache based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of active merchants to cache.
	//   - total: The total records to cache.
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int)

	// GetCachedMerchantTrashed retrieves a list of trashed merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the page, page size, and search string.
	//
	// Returns:
	//   - []*response.MerchantResponseDeleteAt: The list of trashed merchants.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool)

	// SetCachedMerchantTrashed stores a list of trashed merchants into cache based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of trashed merchants to cache.
	//   - total: The total records to cache.
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int)

	// GetCachedMerchant retrieves a merchant from cache by its ID.
	// If the cache is found and contains a valid response, it will return the cached
	// merchant. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The merchant ID.
	//
	// Returns:
	//   - *response.MerchantResponse: The cached merchant.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchant(ctx context.Context, id int) (*response.MerchantResponse, bool)

	// SetCachedMerchant stores a merchant into cache based on its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - data: The merchant data to cache.
	SetCachedMerchant(ctx context.Context, data *response.MerchantResponse)

	// GetCachedMerchantsByUserId retrieves a list of merchants associated with the given user ID from cache.
	// If the cache is found and contains a valid response, it will return the cached list.
	// Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userId: The user ID to look up.
	//
	// Returns:
	//   - []*response.MerchantResponse: The cached merchants.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantsByUserId(ctx context.Context, userId int) ([]*response.MerchantResponse, bool)

	// SetCachedMerchantsByUserId stores a list of merchants into cache associated with a user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userId: The user ID.
	//   - data: The list of merchants to cache.
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*response.MerchantResponse)

	// GetCachedMerchantByApiKey retrieves a merchant from cache by its API key.
	// If the cache is found and contains a valid response, it will return the cached
	// merchant. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - apiKey: The API key associated with the merchant.
	//
	// Returns:
	//   - *response.MerchantResponse: The cached merchant.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantByApiKey(ctx context.Context, apiKey string) (*response.MerchantResponse, bool)

	// SetCachedMerchantByApiKey stores a merchant into cache based on its API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - apiKey: The API key associated with the merchant.
	//   - data: The merchant data to cache.
	SetCachedMerchantByApiKey(ctx context.Context, apiKey string, data *response.MerchantResponse)
}

// MerchantDocumentQueryCache is an interface for caching merchant document queries
type MerchantDocumentQueryCache interface {
	// GetCachedMerchantDocuments retrieves a list of merchant documents from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the page, page size, and search string.
	//
	// Returns:
	//   - []*response.MerchantDocumentResponse: The list of merchant documents.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool)

	// SetCachedMerchantDocuments stores a list of merchant documents into cache based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of merchant documents to cache.
	//   - total: The total records to cache.
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int)

	// GetCachedMerchantDocumentsActive retrieves a list of active merchant documents from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the page, page size, and search string.
	//
	// Returns:
	//   - []*response.MerchantDocumentResponseDeleteAt: The list of active merchant documents.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)

	// SetCachedMerchantDocumentsActive stores a list of active merchant documents into cache based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of active merchant documents to cache.
	//   - total: The total records to cache.
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)

	// GetCachedMerchantDocumentsTrashed retrieves a list of trashed merchant documents from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the page, page size, and search string.
	//
	// Returns:
	//   - []*response.MerchantDocumentResponseDeleteAt: The list of trashed merchant documents.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)

	// SetCachedMerchantDocumentsTrashed stores a list of trashed merchant documents into cache based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of trashed merchant documents to cache.
	//   - total: The total records to cache.
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)

	// GetCachedMerchantDocument retrieves a merchant document from cache by its ID.
	// If the cache is found and contains a valid response, it will return the cached
	// document. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The merchant document ID.
	//
	// Returns:
	//   - *response.MerchantDocumentResponse: The cached merchant document.
	//   - bool: Whether the cache is found and valid.
	GetCachedMerchantDocument(ctx context.Context, id int) (*response.MerchantDocumentResponse, bool)

	// SetCachedMerchantDocument stores a merchant document into cache based on its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The merchant document ID.
	//   - data: The merchant document to cache.
	SetCachedMerchantDocument(ctx context.Context, id int, data *response.MerchantDocumentResponse)
}

// MerchantCommandCache is an interface for caching merchant commands
type MerchantCommandCache interface {
	// DeleteCachedMerchant removes a cached merchant by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The merchant ID.
	DeleteCachedMerchant(ctx context.Context, id int)
}

// MerchantDocumentCommandCache is an interface for caching merchant document commands
type MerchantDocumentCommandCache interface {
	// DeleteCachedMerchantDocuments removes cached merchant documents by merchant ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The merchant ID.
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
}

// MerchantTransactionCache is an interface for caching merchant transactions
type MerchantTransactionCache interface {
	// GetCacheAllMerchantTransactions retrieves all merchant transactions from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters for transactions.
	//
	// Returns:
	//   - []*response.MerchantTransactionResponse: The list of transactions.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, bool)

	// SetCacheAllMerchantTransactions stores all merchant transactions into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The list of transactions to cache.
	//   - total: The total records to cache.
	SetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions, data []*response.MerchantTransactionResponse, total *int)

	// GetCacheMerchantTransactions retrieves transactions for a specific merchant ID from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing merchant ID filters.
	//
	// Returns:
	//   - []*response.MerchantTransactionResponse: The list of transactions.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, bool)

	// SetCacheMerchantTransactions stores transactions by merchant ID into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The list of transactions to cache.
	//   - total: The total records to cache.
	SetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById, data []*response.MerchantTransactionResponse, total *int)

	// GetCacheMerchantTransactionApikey retrieves transactions by merchant API key from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing the API key.
	//
	// Returns:
	//   - []*response.MerchantTransactionResponse: The list of transactions.
	//   - *int: The total records.
	//   - bool: Whether the cache is found and valid.
	GetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, bool)

	// SetCacheMerchantTransactionApikey stores transactions by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The list of transactions to cache.
	//   - total: The total records to cache.
	SetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey, data []*response.MerchantTransactionResponse, total *int)
}

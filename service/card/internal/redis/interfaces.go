package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardQueryCache is an interface that defines methods for retrieving card data from the cache store
//
// CardQueryCache defines the caching behavior for card query operations.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/cache.go
type CardQueryCache interface {

	// GetByIdCache retrieves a card from the cache by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - cardID: The ID of the card.
	//
	// Returns:
	//   - *response.CardResponse: The cached card data.
	//   - bool: Whether the data was found in the cache.
	GetByIdCache(ctx context.Context, cardID int) (*response.CardResponse, bool)

	// GetByUserIDCache retrieves a card from the cache by user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userID: The ID of the user who owns the card.
	//
	// Returns:
	//   - *response.CardResponse: The cached card data.
	//   - bool: Whether the data was found in the cache.
	GetByUserIDCache(ctx context.Context, userID int) (*response.CardResponse, bool)

	// GetByCardNumberCache retrieves a card from the cache by its card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - cardNumber: The card number.
	//
	// Returns:
	//   - *response.CardResponse: The cached card data.
	//   - bool: Whether the data was found in the cache.
	GetByCardNumberCache(ctx context.Context, cardNumber string) (*response.CardResponse, bool)

	// GetFindAllCache retrieves a paginated list of cards from the cache based on the given request filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing search filters such as keyword, page number, and page size.
	//
	// Returns:
	//   - []*response.CardResponse: Slice of cached card data.
	//   - *int: Total number of matching card records.
	//   - bool: Whether the data was found in the cache.
	GetFindAllCache(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponse, *int, bool)

	// GetByActiveCache retrieves a paginated list of active (non-deleted) cards from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing search filters such as keyword, page number, and page size.
	//
	// Returns:
	//   - []*response.CardResponseDeleteAt: Slice of cached active cards.
	//   - *int: Total number of matching records.
	//   - bool: Whether the data was found in the cache.
	GetByActiveCache(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool)

	// GetByTrashedCache retrieves a paginated list of trashed (soft-deleted) cards from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing search filters such as keyword, page number, and page size.
	//
	// Returns:
	//   - []*response.CardResponseDeleteAt: Slice of cached trashed cards.
	//   - *int: Total number of matching records.
	//   - bool: Whether the data was found in the cache.
	GetByTrashedCache(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, bool)

	// SetByIdCache stores a card in the cache by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - cardID: The ID of the card.
	//   - data: The card data to be cached.
	SetByIdCache(ctx context.Context, cardID int, data *response.CardResponse)

	// SetByUserIDCache stores a card in the cache by user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userID: The ID of the user.
	//   - data: The card data to be cached.
	SetByUserIDCache(ctx context.Context, userID int, data *response.CardResponse)

	// SetByCardNumberCache stores a card in the cache by its card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - cardNumber: The card number.
	//   - data: The card data to be cached.
	SetByCardNumberCache(ctx context.Context, cardNumber string, data *response.CardResponse)

	// SetFindAllCache stores a paginated list of cards in the cache based on the given request filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing search filters such as keyword, page number, and page size.
	//   - data: Slice of card data to be cached.
	//   - totalRecords: Total number of records matching the request.
	SetFindAllCache(ctx context.Context, req *requests.FindAllCards, data []*response.CardResponse, totalRecords *int)

	// SetByActiveCache stores a paginated list of active cards in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing search filters.
	//   - data: Slice of active card data to be cached.
	//   - totalRecords: Total number of matching records.
	SetByActiveCache(ctx context.Context, req *requests.FindAllCards, data []*response.CardResponseDeleteAt, totalRecords *int)

	// SetByTrashedCache stores a paginated list of trashed cards in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing search filters.
	//   - data: Slice of trashed card data to be cached.
	//   - totalRecords: Total number of matching records.
	SetByTrashedCache(ctx context.Context, req *requests.FindAllCards, data []*response.CardResponseDeleteAt, totalRecords *int)

	// DeleteByIdCache removes a card from the cache by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - cardID: The ID of the card to be deleted.
	DeleteByIdCache(ctx context.Context, cardID int)

	// DeleteByUserIDCache removes a card from the cache by user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userID: The ID of the user whose card should be removed.
	DeleteByUserIDCache(ctx context.Context, userID int)

	// DeleteByCardNumberCache removes a card from the cache by its card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - cardNumber: The card number of the card to be removed.
	DeleteByCardNumberCache(ctx context.Context, cardNumber string)
}

// CardCommandCache defines the caching behavior for card command operations.
type CardCommandCache interface {
	// DeleteCardCommandCache removes the cache entry associated with the specified card ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the card whose command cache entry should be removed.
	//
	// Behavior:
	//   - Formats the cache key using the given card ID.
	//   - Deletes the corresponding entry from the cache store.
	DeleteCardCommandCache(ctx context.Context, id int)
}

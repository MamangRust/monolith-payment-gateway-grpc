package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// TopupQueryCache is an interface that defines methods for querying and caching topup data.
type TopupQueryCache interface {
	// GetCachedTopupsCache retrieves cached list of topups based on the given filter request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter including pagination and search.
	//
	// Returns:
	//   - []*response.TopupResponse: Cached topup responses.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponse, *int, bool)

	// SetCachedTopupsCache stores the topup responses and total record count in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The original request used as the cache key.
	//   - data: The topup response data to cache.
	//   - total: The total number of records.
	SetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups, data []*response.TopupResponse, total *int)

	// GetCacheTopupByCardCache retrieves cached topups by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and optional filters.
	//
	// Returns:
	//   - []*response.TopupResponse: Cached topups associated with the card.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, bool)

	// SetCacheTopupByCardCache stores the topups associated with a card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request used to generate the cache key.
	//   - data: Topup response data to cache.
	//   - total: Total number of records.
	SetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber, data []*response.TopupResponse, total *int)

	// GetCachedTopupActiveCache retrieves cached list of active topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request used to generate the cache key.
	//
	// Returns:
	//   - []*response.TopupResponseDeleteAt: List of active (non-deleted) topups.
	//   - *int: Total records.
	//   - bool: Whether the cache was found.
	GetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool)

	// SetCachedTopupActiveCache stores the active topups in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The original request used as the cache key.
	//   - data: Topup response data to cache.
	//   - total: Total number of records.
	SetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int)

	// GetCachedTopupTrashedCache retrieves cached list of trashed (soft-deleted) topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request used to generate the cache key.
	//
	// Returns:
	//   - []*response.TopupResponseDeleteAt: List of trashed topups.
	//   - *int: Total records.
	//   - bool: Whether the cache was found.
	GetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, bool)

	// SetCachedTopupTrashedCache stores the trashed (soft-deleted) topups in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request used to generate the cache key.
	//   - data: Topup response data to cache.
	//   - total: Total number of records.
	SetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups, data []*response.TopupResponseDeleteAt, total *int)

	// GetCachedTopupCache retrieves a single topup record from the cache by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The unique topup ID.
	//
	// Returns:
	//   - *response.TopupResponse: The cached topup response.
	//   - bool: Whether the cache was found.
	GetCachedTopupCache(ctx context.Context, id int) (*response.TopupResponse, bool)

	// SetCachedTopupCache stores a single topup response in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - data: The topup response to be cached.
	SetCachedTopupCache(ctx context.Context, data *response.TopupResponse)
}

// TopupCommandCache is an interface for the topup command cache.
type TopupCommandCache interface {
	// DeleteCachedTopupCache removes cached topup data by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the topup to remove from cache.
	DeleteCachedTopupCache(ctx context.Context, id int)
}

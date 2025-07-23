package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// WithdrawQueryCache handles caching operations for withdraw query results.
type WithdrawQueryCache interface {
	// GetCachedWithdrawsCache retrieves cached list of all withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for finding all withdraws.
	//
	// Returns:
	//   - []*response.WithdrawResponse: List of withdraws.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, bool)

	// SetCachedWithdrawsCache stores a list of withdraws in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters used for caching.
	//   - data: The withdraw response data to cache.
	//   - total: Total number of records.
	SetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws, data []*response.WithdrawResponse, total *int)

	// GetCachedWithdrawByCardCache retrieves cached withdraws for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for finding withdraws by card number.
	//
	// Returns:
	//   - []*response.WithdrawResponse: List of withdraws.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, bool)

	// SetCachedWithdrawByCardCache stores withdraws for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters used for caching.
	//   - data: The withdraw response data to cache.
	//   - total: Total number of records.
	SetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber, data []*response.WithdrawResponse, total *int)

	// GetCachedWithdrawActiveCache retrieves cached active (non-deleted) withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for finding active withdraws.
	//
	// Returns:
	//   - []*response.WithdrawResponseDeleteAt: List of active withdraws.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool)

	// SetCachedWithdrawActiveCache stores active withdraws in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters used for caching.
	//   - data: The active withdraw response data to cache.
	//   - total: Total number of records.
	SetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int)

	// GetCachedWithdrawTrashedCache retrieves cached trashed (soft-deleted) withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for finding trashed withdraws.
	//
	// Returns:
	//   - []*response.WithdrawResponseDeleteAt: List of trashed withdraws.
	//   - *int: Total number of records.
	//   - bool: Whether the cache was found.
	GetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, bool)

	// SetCachedWithdrawTrashedCache stores trashed withdraws in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters used for caching.
	//   - data: The trashed withdraw response data to cache.
	//   - total: Total number of records.
	SetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws, data []*response.WithdrawResponseDeleteAt, total *int)

	// GetCachedWithdrawCache retrieves cached withdraw by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the withdraw to retrieve.
	//
	// Returns:
	//   - *response.WithdrawResponse: The withdraw response.
	//   - bool: Whether the cache was found.
	GetCachedWithdrawCache(ctx context.Context, id int) (*response.WithdrawResponse, bool)

	// SetCachedWithdrawCache stores a withdraw record in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - data: The withdraw response to cache.
	SetCachedWithdrawCache(ctx context.Context, data *response.WithdrawResponse)
}

// WithdrawCommandCache handles caching operations related to mutation of withdraw records.
type WithdrawCommandCache interface {
	// DeleteCachedWithdrawCache deletes a cached withdraw entry by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the withdraw cache to delete.
	DeleteCachedWithdrawCache(ctx context.Context, id int)
}

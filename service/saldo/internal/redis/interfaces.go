package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// SaldoQueryCache is an interface for the cache store for saldo query operations.
type SaldoQueryCache interface {
	// GetCachedSaldos retrieves a list of saldos from the cache based on filter parameters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination and search filters.
	//
	// Returns:
	//   - []*response.SaldoResponse: The list of saldos.
	//   - *int: The total number of records.
	//   - bool: Whether the cache was found and valid.
	GetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, bool)

	// SetCachedSaldos stores a list of saldos in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of saldos to be cached.
	//   - totalRecords: The total number of records.
	SetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos, data []*response.SaldoResponse, totalRecords *int)

	// GetCachedSaldoById retrieves a saldo by its ID from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo.
	//
	// Returns:
	//   - *response.SaldoResponse: The cached saldo data.
	//   - bool: Whether the cache was found and valid.
	GetCachedSaldoById(ctx context.Context, saldo_id int) (*response.SaldoResponse, bool)

	// SetCachedSaldoById stores a saldo by its ID in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo.
	//   - data: The saldo data to cache.
	SetCachedSaldoById(ctx context.Context, saldo_id int, data *response.SaldoResponse)

	// GetCachedSaldoByCardNumber retrieves a saldo by card number from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number.
	//
	// Returns:
	//   - *response.SaldoResponse: The cached saldo data.
	//   - bool: Whether the cache was found and valid.
	GetCachedSaldoByCardNumber(ctx context.Context, card_number string) (*response.SaldoResponse, bool)

	// SetCachedSaldoByCardNumber stores a saldo by card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number.
	//   - data: The saldo data to cache.
	SetCachedSaldoByCardNumber(ctx context.Context, card_number string, data *response.SaldoResponse)

	// GetCachedSaldoByActive retrieves a list of active (non-deleted) saldos from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filter parameters.
	//
	// Returns:
	//   - []*response.SaldoResponseDeleteAt: The list of active saldos.
	//   - *int: The total number of records.
	//   - bool: Whether the cache was found and valid.
	GetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool)

	// SetCachedSaldoByActive stores a list of active (non-deleted) saldos in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of active saldos to be cached.
	//   - totalRecords: The total number of records.
	SetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos, data []*response.SaldoResponseDeleteAt, totalRecords *int)

	// GetCachedSaldoByTrashed retrieves a list of trashed (soft-deleted) saldos from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filter parameters.
	//
	// Returns:
	//   - []*response.SaldoResponseDeleteAt: The list of trashed saldos.
	//   - *int: The total number of records.
	//   - bool: Whether the cache was found and valid.
	GetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, bool)

	// SetCachedSaldoByTrashed stores a list of trashed (soft-deleted) saldos in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The list of trashed saldos to be cached.
	//   - totalRecords: The total number of records.
	SetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos, data []*response.SaldoResponseDeleteAt, totalRecords *int)
}


// SaldoCommandCache defines the interface for caching saldo commands.
type SaldoCommandCache interface {
	// DeleteSaldoCache removes the cached saldo entry by saldo ID.
	//
	// This is typically called after a saldo record is updated, deleted,
	// or otherwise invalidated to ensure cache consistency.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo to be removed from cache.
	DeleteSaldoCache(ctx context.Context, saldo_id int)
}

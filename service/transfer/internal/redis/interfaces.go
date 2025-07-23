package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// TransferQueryCache defines cache operations for transfer-related queries.
type TransferQueryCache interface {
	// GetCachedTransfersCache retrieves cached list of transfers.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter for transfer data.
	//
	// Returns:
	//   - []*response.TransferResponse: List of transfers.
	//   - *int: Total count of transfers.
	//   - bool: Whether the cache was found.
	GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponse, *int, bool)

	// SetCachedTransfersCache stores list of transfers into the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter used to cache the data.
	//   - data: List of transfers to cache.
	//   - total: Total count of transfers.
	SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers, data []*response.TransferResponse, total *int)

	// GetCachedTransferActiveCache retrieves cached list of active (non-trashed) transfers.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter for active transfers.
	//
	// Returns:
	//   - []*response.TransferResponseDeleteAt: List of active transfers.
	//   - *int: Total count of active transfers.
	//   - bool: Whether the cache was found.
	GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, bool)

	// SetCachedTransferActiveCache stores list of active (non-trashed) transfers into the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter used to cache the data.
	//   - data: List of active transfers to cache.
	//   - total: Total count of active transfers.
	SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers, data []*response.TransferResponseDeleteAt, total *int)

	// GetCachedTransferTrashedCache retrieves cached list of trashed transfers.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter for trashed transfers.
	//
	// Returns:
	//   - []*response.TransferResponseDeleteAt: List of trashed transfers.
	//   - *int: Total count of trashed transfers.
	//   - bool: Whether the cache was found.
	GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, bool)

	// SetCachedTransferTrashedCache stores list of trashed transfers into the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request filter used to cache the data.
	//   - data: List of trashed transfers to cache.
	//   - total: Total count of trashed transfers.
	SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers, data []*response.TransferResponseDeleteAt, total *int)

	// GetCachedTransferCache retrieves a specific transfer by ID from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the transfer to retrieve.
	//
	// Returns:
	//   - *response.TransferResponse: Transfer response.
	//   - bool: Whether the cache was found.
	GetCachedTransferCache(ctx context.Context, id int) (*response.TransferResponse, bool)

	// SetCachedTransferCache stores a specific transfer into the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - data: The transfer data to cache.
	SetCachedTransferCache(ctx context.Context, data *response.TransferResponse)

	// GetCachedTransferByFrom retrieves cached transfers filtered by source card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - from: The card number from which the transfer was made.
	//
	// Returns:
	//   - []*response.TransferResponse: List of transfers.
	//   - bool: Whether the cache was found.
	GetCachedTransferByFrom(ctx context.Context, from string) ([]*response.TransferResponse, bool)

	// SetCachedTransferByFrom stores cached transfers by source card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - from: The card number from which the transfer was made.
	//   - data: List of transfers to cache.
	SetCachedTransferByFrom(ctx context.Context, from string, data []*response.TransferResponse)

	// GetCachedTransferByTo retrieves cached transfers filtered by destination card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - to: The card number to which the transfer was made.
	//
	// Returns:
	//   - []*response.TransferResponse: List of transfers.
	//   - bool: Whether the cache was found.
	GetCachedTransferByTo(ctx context.Context, to string) ([]*response.TransferResponse, bool)

	// SetCachedTransferByTo stores cached transfers by destination card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - to: The card number to which the transfer was made.
	//   - data: List of transfers to cache.
	SetCachedTransferByTo(ctx context.Context, to string, data []*response.TransferResponse)
}

type TransferCommandCache interface {
	// DeleteTransferCache removes a cached transfer by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The transfer ID.
	DeleteTransferCache(ctx context.Context, id int)
}

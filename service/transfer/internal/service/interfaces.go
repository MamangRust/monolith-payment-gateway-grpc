package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// TransferQueryService provides methods for querying transfer data.

// TransferQueryService defines query operations for retrieving transfer data.
type TransferQueryService interface {
	// FindAll retrieves all transfer records based on filter and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The filter and pagination parameters.
	//
	// Returns:
	//   - []*response.TransferResponse: List of transfer responses.
	//   - *int: Total number of transfer records.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponse, *int, *response.ErrorResponse)

	// FindById retrieves a single transfer record by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transferId: The ID of the transfer to retrieve.
	//
	// Returns:
	//   - *response.TransferResponse: The transfer response.
	//   - *response.ErrorResponse: Error response if not found or failed.
	FindById(ctx context.Context, transferId int) (*response.TransferResponse, *response.ErrorResponse)

	// FindByActive retrieves all active (non-deleted) transfer records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The filter and pagination parameters.
	//
	// Returns:
	//   - []*response.TransferResponseDeleteAt: List of active transfer responses with deleted_at info.
	//   - *int: Total number of active records.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves all trashed (soft-deleted) transfer records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The filter and pagination parameters.
	//
	// Returns:
	//   - []*response.TransferResponseDeleteAt: List of trashed transfer responses with deleted_at info.
	//   - *int: Total number of trashed records.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)

	// FindTransferByTransferFrom retrieves transfers by sender card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_from: The sender card number.
	//
	// Returns:
	//   - []*response.TransferResponse: List of transfer responses.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse)

	// FindTransferByTransferTo retrieves transfers by receiver card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_to: The receiver card number.
	//
	// Returns:
	//   - []*response.TransferResponse: List of transfer responses.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse)
}

// TransferCommandService handles business logic for creating, updating,
// deleting, and restoring transfer records.
type TransferCommandService interface {
	// CreateTransaction creates a new transfer transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing transfer details.
	//
	// Returns:
	//   - *response.TransferResponse: The created transfer data.
	//   - *response.ErrorResponse: Error details if operation fails.
	CreateTransaction(ctx context.Context, request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse)

	// UpdateTransaction updates an existing transfer transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing updated transfer details.
	//
	// Returns:
	//   - *response.TransferResponse: The updated transfer data.
	//   - *response.ErrorResponse: Error details if operation fails.
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse)

	// TrashedTransfer marks a transfer transaction as trashed (soft delete).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_id: The ID of the transfer to be trashed.
	//
	// Returns:
	//   - *response.TransferResponse: The trashed transfer data.
	//   - *response.ErrorResponse: Error details if operation fails.
	TrashedTransfer(ctx context.Context, transfer_id int) (*response.TransferResponseDeleteAt, *response.ErrorResponse)

	// RestoreTransfer restores a previously trashed transfer transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_id: The ID of the transfer to be restored.
	//
	// Returns:
	//   - *response.TransferResponse: The restored transfer data.
	//   - *response.ErrorResponse: Error details if operation fails.
	RestoreTransfer(ctx context.Context, transfer_id int) (*response.TransferResponse, *response.ErrorResponse)

	// DeleteTransferPermanent permanently deletes a transfer transaction.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_id: The ID of the transfer to be permanently deleted.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - *response.ErrorResponse: Error details if operation fails.
	DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, *response.ErrorResponse)

	// RestoreAllTransfer restores all trashed transfer transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the restore operation was successful.
	//   - *response.ErrorResponse: Error details if operation fails.
	RestoreAllTransfer(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllTransferPermanent permanently deletes all trashed transfer transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - *response.ErrorResponse: Error details if operation fails.
	DeleteAllTransferPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

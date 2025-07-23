package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// WithdrawQueryService defines query operations for fetching withdraw data.
type WithdrawQueryService interface {
	// FindAll retrieves all withdraws based on the given request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filters like pagination, date range, etc.
	//
	// Returns:
	//   - []*response.WithdrawResponse: List of withdraws.
	//   - *int: Total number of records matching the filter.
	//   - *response.ErrorResponse: Error details if any.
	FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)

	// FindAllByCardNumber retrieves all withdraws filtered by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number and other filters.
	//
	// Returns:
	//   - []*response.WithdrawResponse: List of withdraws for the specified card.
	//   - *int: Total number of records found.
	//   - *response.ErrorResponse: Error details if any.
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)

	// FindById retrieves a single withdraw record by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - withdrawID: The ID of the withdraw to retrieve.
	//
	// Returns:
	//   - *response.WithdrawResponse: The withdraw data if found.
	//   - *response.ErrorResponse: Error details if any.
	FindById(ctx context.Context, withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse)

	// FindByActive retrieves active withdraw records based on the request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filters like pagination, etc.
	//
	// Returns:
	//   - []*response.WithdrawResponseDeleteAt: List of active withdraws.
	//   - *int: Total number of active records found.
	//   - *response.ErrorResponse: Error details if any.
	FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves soft-deleted withdraw records based on the request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filters like pagination, etc.
	//
	// Returns:
	//   - []*response.WithdrawResponseDeleteAt: List of trashed withdraws.
	//   - *int: Total number of trashed records found.
	//   - *response.ErrorResponse: Error details if any.
	FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)
}

// WithdrawCommandService handles all command operations for withdraw records,
// including creation, update, soft-delete, restore, and permanent deletion.
type WithdrawCommandService interface {
	// Create creates a new withdraw record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request data to create a withdraw.
	//
	// Returns:
	//   - *response.WithdrawResponse: The created withdraw response.
	//   - *response.ErrorResponse: Error information if any occurred.
	Create(ctx context.Context, request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)

	// Update updates an existing withdraw record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request data to update the withdraw.
	//
	// Returns:
	//   - *response.WithdrawResponse: The updated withdraw response.
	//   - *response.ErrorResponse: Error information if any occurred.
	Update(ctx context.Context, request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)

	// TrashedWithdraw soft-deletes a withdraw by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - withdraw_id: The ID of the withdraw to soft-delete.
	//
	// Returns:
	//   - *response.WithdrawResponse: The soft-deleted withdraw response.
	//   - *response.ErrorResponse: Error information if any occurred.
	TrashedWithdraw(ctx context.Context, withdraw_id int) (*response.WithdrawResponseDeleteAt, *response.ErrorResponse)

	// RestoreWithdraw restores a soft-deleted withdraw by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - withdraw_id: The ID of the withdraw to restore.
	//
	// Returns:
	//   - *response.WithdrawResponse: The restored withdraw response.
	//   - *response.ErrorResponse: Error information if any occurred.
	RestoreWithdraw(ctx context.Context, withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)

	// DeleteWithdrawPermanent permanently deletes a withdraw by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - withdraw_id: The ID of the withdraw to delete permanently.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - *response.ErrorResponse: Error information if any occurred.
	DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, *response.ErrorResponse)

	// RestoreAllWithdraw restores all soft-deleted withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all records were successfully restored.
	//   - *response.ErrorResponse: Error information if any occurred.
	RestoreAllWithdraw(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllWithdrawPermanent permanently deletes all soft-deleted withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all records were successfully deleted.
	//   - *response.ErrorResponse: Error information if any occurred.
	DeleteAllWithdrawPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// TopupQueryService defines the read-only operations for querying topup data.
type TopupQueryService interface {
	// FindAll retrieves all topup records based on the given filter request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination criteria.
	//
	// Returns:
	//   - []*response.TopupResponse: List of topup responses.
	//   - *int: Total number of matching records.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindAll(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponse, *int, *response.ErrorResponse)

	// FindAllByCardNumber retrieves all topup records filtered by a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number and filter criteria.
	//
	// Returns:
	//   - []*response.TopupResponse: List of topup responses.
	//   - *int: Total number of matching records.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, *response.ErrorResponse)

	// FindById retrieves a topup record by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topupID: The ID of the topup to be retrieved.
	//
	// Returns:
	//   - *response.TopupResponse: The topup data if found.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindById(ctx context.Context, topupID int) (*response.TopupResponse, *response.ErrorResponse)

	// FindByActive retrieves all active (non-deleted) topup records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination criteria.
	//
	// Returns:
	//   - []*response.TopupResponseDeleteAt: List of active topup records.
	//   - *int: Total number of matching records.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves all soft-deleted (trashed) topup records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination criteria.
	//
	// Returns:
	//   - []*response.TopupResponseDeleteAt: List of trashed topup records.
	//   - *int: Total number of matching records.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
}

// TopupCommandService handles commands for creating, updating, and managing topup records.
type TopupCommandService interface {
	// CreateTopup creates a new topup record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The topup creation request payload.
	//
	// Returns:
	//   - *response.TopupResponse: The newly created topup response.
	//   - *response.ErrorResponse: Error details if creation fails.
	CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)

	// UpdateTopup updates an existing topup record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The topup update request payload.
	//
	// Returns:
	//   - *response.TopupResponse: The updated topup response.
	//   - *response.ErrorResponse: Error details if update fails.
	UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)

	// TrashedTopup marks a topup record as trashed (soft deleted).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The ID of the topup to trash.
	//
	// Returns:
	//   - *response.TopupResponseDeleteAt: The trashed topup response.
	//   - *response.ErrorResponse: Error details if trashing fails.
	TrashedTopup(ctx context.Context, topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse)

	// RestoreTopup restores a previously trashed topup record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The ID of the topup to restore.
	//
	// Returns:
	//   - *response.TopupResponseDeleteAt: The restored topup response.
	//   - *response.ErrorResponse: Error details if restoration fails.
	RestoreTopup(ctx context.Context, topup_id int) (*response.TopupResponse, *response.ErrorResponse)

	// DeleteTopupPermanent permanently deletes a topup record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The ID of the topup to delete.
	//
	// Returns:
	//   - bool: True if deletion succeeded.
	//   - *response.ErrorResponse: Error details if deletion fails.
	DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, *response.ErrorResponse)

	// RestoreAllTopup restores all trashed topup records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all records were successfully restored.
	//   - *response.ErrorResponse: Error details if restoration fails.
	RestoreAllTopup(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllTopupPermanent permanently deletes all trashed topup records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all records were successfully deleted.
	//   - *response.ErrorResponse: Error details if deletion fails.
	DeleteAllTopupPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

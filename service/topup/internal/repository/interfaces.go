package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// SaldoRepository is an interface that defines the methods for interacting with the saldo records in the database
type SaldoRepository interface {
	// FindByCardNumber retrieves a saldo record by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The unique card number to find the saldo.
	//
	// Returns:
	//   - *record.SaldoRecord: The saldo record found.
	//   - error: Error if the query fails or saldo is not found.
	FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error)

	// UpdateSaldoBalance updates the balance of a saldo record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing saldo ID and new balance.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: Error if the update fails.
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
}

// TopupQueryRepository is an interface that defines the methods for querying the topup records in the database
type TopupQueryRepository interface {
	// FindAllTopups retrieves a paginated list of all topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for pagination and filters.
	//
	// Returns:
	//   - []*record.TopupRecord: List of topup records.
	//   - *int: Total number of records.
	//   - error: Error if the query fails.
	FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)

	// FindByActive retrieves a paginated list of active topup records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for pagination and filters.
	//
	// Returns:
	//   - []*record.TopupRecord: List of active topup records.
	//   - *int: Total number of records.
	//   - error: Error if the query fails.
	FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)

	// FindByTrashed retrieves a paginated list of trashed topup records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for pagination and filters.
	//
	// Returns:
	//   - []*record.TopupRecord: List of trashed topup records.
	//   - *int: Total number of records.
	//   - error: Error if the query fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)

	// FindAllTopupByCardNumber retrieves all topups associated with a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and pagination info.
	//
	// Returns:
	//   - []*record.TopupRecord: List of topups associated with the card.
	//   - *int: Total number of records.
	//   - error: Error if the query fails.
	FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*record.TopupRecord, *int, error)

	// FindById retrieves a topup record by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The unique ID of the topup.
	//
	// Returns:
	//   - *record.TopupRecord: The found topup record.
	//   - error: Error if the query fails or topup is not found.
	FindById(ctx context.Context, topup_id int) (*record.TopupRecord, error)
}

// TopupCommandRepository defines methods for performing write operations on topup records.
type TopupCommandRepository interface {
	// CreateTopup inserts a new topup record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The data needed to create a new topup.
	//
	// Returns:
	//   - *record.TopupRecord: The created topup record.
	//   - error: Error if creation fails.
	CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*record.TopupRecord, error)

	// UpdateTopup updates an existing topup record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The data used to update the topup.
	//
	// Returns:
	//   - *record.TopupRecord: The updated topup record.
	//   - error: Error if update fails.
	UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*record.TopupRecord, error)

	// UpdateTopupAmount updates the amount of a specific topup.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The new amount data.
	//
	// Returns:
	//   - *record.TopupRecord: The updated topup record.
	//   - error: Error if update fails.
	UpdateTopupAmount(ctx context.Context, request *requests.UpdateTopupAmount) (*record.TopupRecord, error)

	// UpdateTopupStatus updates the status of a topup (e.g., success, failed).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The status update data.
	//
	// Returns:
	//   - *record.TopupRecord: The updated topup record.
	//   - error: Error if update fails.
	UpdateTopupStatus(ctx context.Context, request *requests.UpdateTopupStatus) (*record.TopupRecord, error)

	// TrashedTopup soft deletes a topup by marking it as trashed.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The ID of the topup to trash.
	//
	// Returns:
	//   - *record.TopupRecord: The trashed topup record.
	//   - error: Error if trashing fails.
	TrashedTopup(ctx context.Context, topup_id int) (*record.TopupRecord, error)

	// RestoreTopup restores a previously trashed topup.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The ID of the topup to restore.
	//
	// Returns:
	//   - *record.TopupRecord: The restored topup record.
	//   - error: Error if restoration fails.
	RestoreTopup(ctx context.Context, topup_id int) (*record.TopupRecord, error)

	// DeleteTopupPermanent permanently deletes a topup from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - topup_id: The ID of the topup to delete.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: Error if deletion fails.
	DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error)

	// RestoreAllTopup restores all trashed topups in the system.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the restoration was successful.
	//   - error: Error if operation fails.
	RestoreAllTopup(ctx context.Context) (bool, error)

	// DeleteAllTopupPermanent permanently deletes all trashed topups from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: Error if deletion fails.
	DeleteAllTopupPermanent(ctx context.Context) (bool, error)
}

// CardRepository defines methods for interacting with card records in the database.
type CardRepository interface {
	// FindUserCardByCardNumber retrieves card data with user email by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to search.
	//
	// Returns:
	//   - *record.CardEmailRecord: The card record including associated user email.
	//   - error: Error if retrieval fails.
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error)

	// FindCardByCardNumber retrieves a card record by its card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to search.
	//
	// Returns:
	//   - *record.CardRecord: The card record data.
	//   - error: Error if retrieval fails.
	FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error)

	// UpdateCard updates an existing card's data.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The updated card information.
	//
	// Returns:
	//   - *record.CardRecord: The updated card record.
	//   - error: Error if update fails.
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*record.CardRecord, error)
}

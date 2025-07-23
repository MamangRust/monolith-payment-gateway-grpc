package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// SaldoRepository handles operations related to saldo records.
type SaldoRepository interface {
	// FindByCardNumber retrieves a saldo record by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to search for.
	//
	// Returns:
	//   - *record.SaldoRecord: The found saldo record.
	//   - error: An error if the operation fails.
	FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error)

	// UpdateSaldoBalance updates the saldo balance for a specific card.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The update request containing balance data.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: An error if the update fails.
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)

	// UpdateSaldoWithdraw updates the saldo balance after a withdrawal.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The update request related to withdrawal.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: An error if the update fails.
	UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error)
}

// WithdrawQueryRepository handles query operations for withdraw records.
type WithdrawQueryRepository interface {
	// FindAll retrieves all withdraw records with pagination and filtering.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination information.
	//
	// Returns:
	//   - []*record.WithdrawRecord: List of withdraw records.
	//   - *int: Total count.
	//   - error: An error if the operation fails.
	FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)

	// FindByActive retrieves active (non-deleted) withdraw records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination information.
	//
	// Returns:
	//   - []*record.WithdrawRecord: List of active withdraw records.
	//   - *int: Total count.
	//   - error: An error if the operation fails.
	FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)

	// FindByTrashed retrieves soft-deleted withdraw records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination information.
	//
	// Returns:
	//   - []*record.WithdrawRecord: List of trashed withdraw records.
	//   - *int: Total count.
	//   - error: An error if the operation fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)

	// FindAllByCardNumber retrieves all withdraw records associated with a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number and filter info.
	//
	// Returns:
	//   - []*record.WithdrawRecord: List of withdraw records for the card.
	//   - *int: Total count.
	//   - error: An error if the operation fails.
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*record.WithdrawRecord, *int, error)

	// FindById retrieves a withdraw record by its unique ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the withdraw record.
	//
	// Returns:
	//   - *record.WithdrawRecord: The withdraw record if found.
	//   - error: An error if the operation fails or the record is not found.
	FindById(ctx context.Context, id int) (*record.WithdrawRecord, error)
}

// WithdrawCommandRepository defines the command operations for withdraw records in the database.
type WithdrawCommandRepository interface {
	// CreateWithdraw inserts a new withdraw record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing withdraw data.
	//
	// Returns:
	//   - *record.WithdrawRecord: The newly created withdraw record.
	//   - error: An error if the operation fails.
	CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error)

	// UpdateWithdraw modifies an existing withdraw record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing updated withdraw data.
	//
	// Returns:
	//   - *record.WithdrawRecord: The updated withdraw record.
	//   - error: An error if the operation fails.
	UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error)

	// UpdateWithdrawStatus updates the status of a withdraw record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing the new status.
	//
	// Returns:
	//   - *record.WithdrawRecord: The updated withdraw record with the new status.
	//   - error: An error if the operation fails.
	UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error)

	// TrashedWithdraw soft-deletes a withdraw record by marking it as trashed.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - WithdrawID: The ID of the withdraw to be trashed.
	//
	// Returns:
	//   - *record.WithdrawRecord: The trashed withdraw record.
	//   - error: An error if the operation fails.
	TrashedWithdraw(ctx context.Context, WithdrawID int) (*record.WithdrawRecord, error)

	// RestoreWithdraw restores a previously trashed withdraw record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - WithdrawID: The ID of the withdraw to be restored.
	//
	// Returns:
	//   - *record.WithdrawRecord: The restored withdraw record.
	//   - error: An error if the operation fails.
	RestoreWithdraw(ctx context.Context, WithdrawID int) (*record.WithdrawRecord, error)

	// DeleteWithdrawPermanent permanently deletes a withdraw record from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - WithdrawID: The ID of the withdraw to be permanently deleted.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: An error if the operation fails.
	DeleteWithdrawPermanent(ctx context.Context, WithdrawID int) (bool, error)

	// RestoreAllWithdraw restores all trashed withdraw records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the restoration was successful.
	//   - error: An error if the operation fails.
	RestoreAllWithdraw(ctx context.Context) (bool, error)

	// DeleteAllWithdrawPermanent permanently deletes all trashed withdraw records from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - error: An error if the operation fails.
	DeleteAllWithdrawPermanent(ctx context.Context) (bool, error)
}

type CardRepository interface {
	// FindUserCardByCardNumber retrieves a card record along with associated user email by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to look up.
	//
	// Returns:
	//   - *record.CardEmailRecord: The card and user email record if found.
	//   - error: An error if the operation fails or no record is found.
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error)
}

package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// SaldoRepository defines operations related to saldo records in the database.
type SaldoRepository interface {
	// FindByCardNumber retrieves a saldo record by the given card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to search the saldo for.
	//
	// Returns:
	//   - *record.SaldoRecord: The retrieved saldo record.
	//   - error: Any error encountered during the operation.
	FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error)

	// UpdateSaldoBalance updates the saldo balance based on the provided request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing card number and the new balance.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: Any error encountered during the operation.
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
}

// CardRepository defines operations related to card records in the database.
type CardRepository interface {
	// FindUserCardByCardNumber retrieves user-related card info by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to search for.
	//
	// Returns:
	//   - *record.CardEmailRecord: The card and user email info.
	//   - error: Any error encountered during the operation.
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error)

	// FindCardByCardNumber retrieves a card record by its card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to search for.
	//
	// Returns:
	//   - *record.CardRecord: The retrieved card record.
	//   - error: Any error encountered during the operation.
	FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error)
}

// TransferQueryRepository defines operations for querying transfer records from the database.
type TransferQueryRepository interface {
	// FindAll retrieves all transfer records with optional filtering and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for filtering and pagination.
	//
	// Returns:
	//   - []*record.TransferRecord: List of transfer records.
	//   - *int: Total number of records (for pagination).
	//   - error: Any error encountered during the operation.
	FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*record.TransferRecord, *int, error)

	// FindByActive retrieves all active (non-trashed) transfer records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for filtering and pagination.
	//
	// Returns:
	//   - []*record.TransferRecord: List of active transfer records.
	//   - *int: Total number of records (for pagination).
	//   - error: Any error encountered during the operation.
	FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*record.TransferRecord, *int, error)

	// FindById retrieves a single transfer record by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the transfer record.
	//
	// Returns:
	//   - *record.TransferRecord: The transfer record, if found.
	//   - error: Any error encountered during the operation.
	FindById(ctx context.Context, id int) (*record.TransferRecord, error)

	// FindByTrashed retrieves all soft-deleted (trashed) transfer records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request parameters for filtering and pagination.
	//
	// Returns:
	//   - []*record.TransferRecord: List of trashed transfer records.
	//   - *int: Total number of records (for pagination).
	//   - error: Any error encountered during the operation.
	FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*record.TransferRecord, *int, error)

	// FindTransferByTransferFrom retrieves all transfer records where the given card is the sender.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_from: The sender card number.
	//
	// Returns:
	//   - []*record.TransferRecord: List of transfer records from the specified sender.
	//   - error: Any error encountered during the operation.
	FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*record.TransferRecord, error)

	// FindTransferByTransferTo retrieves all transfer records where the given card is the receiver.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_to: The receiver card number.
	//
	// Returns:
	//   - []*record.TransferRecord: List of transfer records to the specified receiver.
	//   - error: Any error encountered during the operation.
	FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*record.TransferRecord, error)
}

// TransferCommandRepository defines command operations for managing transfer records.
type TransferCommandRepository interface {
	// CreateTransfer inserts a new transfer record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The data needed to create the transfer.
	//
	// Returns:
	//   - *record.TransferRecord: The created transfer record.
	//   - error: Any error encountered during the operation.
	CreateTransfer(ctx context.Context, request *requests.CreateTransferRequest) (*record.TransferRecord, error)

	// UpdateTransfer updates the details of an existing transfer.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The updated transfer data.
	//
	// Returns:
	//   - *record.TransferRecord: The updated transfer record.
	//   - error: Any error encountered during the operation.
	UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*record.TransferRecord, error)

	// UpdateTransferAmount updates only the amount field of a transfer.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The new amount information.
	//
	// Returns:
	//   - *record.TransferRecord: The updated transfer record.
	//   - error: Any error encountered during the operation.
	UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error)

	// UpdateTransferStatus updates the status of a transfer.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The new status for the transfer.
	//
	// Returns:
	//   - *record.TransferRecord: The updated transfer record.
	//   - error: Any error encountered during the operation.
	UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*record.TransferRecord, error)

	// TrashedTransfer marks a transfer as deleted (soft delete).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_id: The ID of the transfer to be trashed.
	//
	// Returns:
	//   - *record.TransferRecord: The trashed transfer record.
	//   - error: Any error encountered during the operation.
	TrashedTransfer(ctx context.Context, transfer_id int) (*record.TransferRecord, error)

	// RestoreTransfer restores a previously trashed transfer.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_id: The ID of the transfer to be restored.
	//
	// Returns:
	//   - *record.TransferRecord: The restored transfer record.
	//   - error: Any error encountered during the operation.
	RestoreTransfer(ctx context.Context, transfer_id int) (*record.TransferRecord, error)

	// DeleteTransferPermanent permanently deletes a transfer from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - transfer_id: The ID of the transfer to be permanently deleted.
	//
	// Returns:
	//   - bool: Indicates if the deletion was successful.
	//   - error: Any error encountered during the operation.
	DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error)

	// RestoreAllTransfer restores all trashed transfers.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Indicates if the restore operation was successful.
	//   - error: Any error encountered during the operation.
	RestoreAllTransfer(ctx context.Context) (bool, error)

	// DeleteAllTransferPermanent permanently deletes all trashed transfers.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Indicates if the deletion was successful.
	//   - error: Any error encountered during the operation.
	DeleteAllTransferPermanent(ctx context.Context) (bool, error)
}

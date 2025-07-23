package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// SaldoQueryRepository is an interface for the repository that handles saldo query operations.
type SaldoQueryRepository interface {
	// FindAllSaldos retrieves all saldo records based on provided filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination, filtering, or search criteria.
	//
	// Returns:
	//   - []*record.SaldoRecord: The list of saldo records.
	//   - *int: The total number of records found.
	//   - error: An error if the query fails.
	FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)

	// FindByActive retrieves all active saldo records (not soft-deleted).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination or filtering options.
	//
	// Returns:
	//   - []*record.SaldoRecord: The list of active saldo records.
	//   - *int: The total number of active records.
	//   - error: An error if the query fails.
	FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)

	// FindByTrashed retrieves all trashed saldo records (soft-deleted).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination or filtering options.
	//
	// Returns:
	//   - []*record.SaldoRecord: The list of trashed saldo records.
	//   - *int: The total number of trashed records.
	//   - error: An error if the query fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)

	// FindById retrieves a saldo record by its unique ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The unique saldo ID to query.
	//
	// Returns:
	//   - *record.SaldoRecord: The saldo record if found.
	//   - error: An error if the record is not found or query fails.
	FindById(ctx context.Context, saldo_id int) (*record.SaldoRecord, error)

	// FindByCardNumber retrieves a saldo record by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number associated with the saldo.
	//
	// Returns:
	//   - *record.SaldoRecord: The saldo record if found.
	//   - error: An error if the record is not found or query fails.
	FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error)
}

// SaldoCommandRepository is an interface for the repository that handles saldo command operations.
type SaldoCommandRepository interface {
	// CreateSaldo inserts a new saldo record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing saldo data to be created.
	//
	// Returns:
	//   - *record.SaldoRecord: The created saldo record.
	//   - error: An error if the insert operation fails.
	CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*record.SaldoRecord, error)

	// UpdateSaldo updates the saldo record by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing updated saldo data.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: An error if the update operation fails.
	UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*record.SaldoRecord, error)

	// UpdateSaldoBalance updates the saldo balance (e.g., after top-up or transaction).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing the updated balance info.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: An error if the update fails.
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)

	// UpdateSaldoWithdraw updates the saldo after a withdrawal operation.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request containing withdrawal update information.
	//
	// Returns:
	//   - *record.SaldoRecord: The updated saldo record.
	//   - error: An error if the update fails.
	UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error)

	// TrashedSaldo marks a saldo record as soft-deleted.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldoID: The ID of the saldo to be soft-deleted.
	//
	// Returns:
	//   - *record.SaldoRecord: The trashed saldo record.
	//   - error: An error if the operation fails.
	TrashedSaldo(ctx context.Context, saldoID int) (*record.SaldoRecord, error)

	// RestoreSaldo restores a previously trashed saldo record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldoID: The ID of the saldo to be restored.
	//
	// Returns:
	//   - *record.SaldoRecord: The restored saldo record.
	//   - error: An error if the operation fails.
	RestoreSaldo(ctx context.Context, saldoID int) (*record.SaldoRecord, error)

	// DeleteSaldoPermanent permanently deletes a saldo record from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo to be permanently deleted.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - error: An error if the operation fails.
	DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error)

	// RestoreAllSaldo restores all trashed saldo records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all saldo records were successfully restored.
	//   - error: An error if the operation fails.
	RestoreAllSaldo(ctx context.Context) (bool, error)

	// DeleteAllSaldoPermanent permanently deletes all trashed saldo records from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if all trashed saldo records were successfully deleted.
	//   - error: An error if the operation fails.
	DeleteAllSaldoPermanent(ctx context.Context) (bool, error)
}

// CardRepository is an interface for the repository that handles card operations.
type CardRepository interface {
	// FindCardByCardNumber retrieves a card record based on the card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number to look up.
	//
	// Returns:
	//   - *record.CardRecord: The found card record if exists.
	//   - error: An error if the card is not found or the query fails.
	FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error)
}

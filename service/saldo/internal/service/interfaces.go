package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// SaldoQueryService is an interface that defines methods for querying saldo data.
type SaldoQueryService interface {
	// FindAll retrieves all saldo records with optional pagination and filtering.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filters such as pagination, status, etc.
	//
	// Returns:
	//   - []*response.SaldoResponse: The list of saldo responses.
	//   - *int: The total number of records found.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindAll(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, *response.ErrorResponse)

	// FindById retrieves a saldo by its unique ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo to retrieve.
	//
	// Returns:
	//   - *response.SaldoResponse: The saldo response if found.
	//   - *response.ErrorResponse: An error response if the saldo is not found or an error occurs.
	FindById(ctx context.Context, saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)

	// FindByCardNumber retrieves a saldo by its associated card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - card_number: The card number associated with the saldo.
	//
	// Returns:
	//   - *response.SaldoResponse: The saldo response if found.
	//   - *response.ErrorResponse: An error response if the saldo is not found or an error occurs.
	FindByCardNumber(ctx context.Context, card_number string) (*response.SaldoResponse, *response.ErrorResponse)

	// FindByActive retrieves all active saldo records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter options such as page and page size.
	//
	// Returns:
	//   - []*response.SaldoResponseDeleteAt: The list of active saldo records.
	//   - *int: The total number of active records.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves all trashed (soft-deleted) saldo records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter options such as page and page size.
	//
	// Returns:
	//   - []*response.SaldoResponseDeleteAt: The list of trashed saldo records.
	//   - *int: The total number of trashed records.
	//   - *response.ErrorResponse: An error response if the operation fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
}

// SaldoCommandService is an interface that defines methods for creating, updating, and deleting saldo records.
type SaldoCommandService interface {
	// CreateSaldo creates a new saldo record in the system.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing saldo creation data.
	//
	// Returns:
	//   - *response.SaldoResponse: The created saldo response.
	//   - *response.ErrorResponse: An error response if creation fails.
	CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)

	// UpdateSaldo updates an existing saldo record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing updated saldo data.
	//
	// Returns:
	//   - *response.SaldoResponse: The updated saldo response.
	//   - *response.ErrorResponse: An error response if update fails.
	UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)

	// TrashSaldo moves a saldo to the trash (soft delete).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo to trash.
	//
	// Returns:
	//   - *response.SaldoResponseDeleteAt: The trashed saldo response.
	//   - *response.ErrorResponse: An error response if trashing fails.
	TrashSaldo(ctx context.Context, saldo_id int) (*response.SaldoResponseDeleteAt, *response.ErrorResponse)

	// RestoreSaldo restores a previously trashed saldo.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo to restore.
	//
	// Returns:
	//   - *response.SaldoResponse: The restored saldo response.
	//   - *response.ErrorResponse: An error response if restoring fails.
	RestoreSaldo(ctx context.Context, saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)

	// DeleteSaldoPermanent permanently deletes a saldo record.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - saldo_id: The ID of the saldo to delete permanently.
	//
	// Returns:
	//   - bool: True if the deletion is successful.
	//   - *response.ErrorResponse: An error response if deletion fails.
	DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, *response.ErrorResponse)

	// RestoreAllSaldo restores all trashed saldo records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if restoration is successful.
	//   - *response.ErrorResponse: An error response if operation fails.
	RestoreAllSaldo(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllSaldoPermanent permanently deletes all trashed saldo records.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if deletion is successful.
	//   - *response.ErrorResponse: An error response if operation fails.
	DeleteAllSaldoPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

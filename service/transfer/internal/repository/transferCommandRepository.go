package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer"
)

// transferCommandRepository is a struct that implements the TransferCommandRepository interface
type transferCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferCommandRecordMapper
}

// NewTransferCommandRepository initializes a new instance of transferCommandRepository with the provided
// database queries, context, and transfer record mapper. This repository is responsible for executing
// command operations related to transfer records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A TransferRecordMapping that provides methods to map database rows to TransferRecord domain models.
//
// Returns:
//   - A pointer to the newly created transferCommandRepository instance.
func NewTransferCommandRepository(db *db.Queries, mapper recordmapper.TransferCommandRecordMapper) TransferCommandRepository {
	return &transferCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateTransfer inserts a new transfer record into the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The data needed to create the transfer.
//
// Returns:
//   - *record.TransferRecord: The created transfer record.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) CreateTransfer(ctx context.Context, request *requests.CreateTransferRequest) (*record.TransferRecord, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrCreateTransferFailed
	}

	so := r.mapper.ToTransferRecord(res)

	return so, nil
}

// UpdateTransfer updates the details of an existing transfer.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The updated transfer data.
//
// Returns:
//   - *record.TransferRecord: The updated transfer record.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferParams{
		TransferID:     int32(*request.TransferID),
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransfer(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferFailed
	}

	so := r.mapper.ToTransferRecord(res)

	return so, nil

}

// UpdateTransferAmount updates only the amount field of a transfer.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The new amount information.
//
// Returns:
//   - *record.TransferRecord: The updated transfer record.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransferAmount(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferAmountFailed
	}

	so := r.mapper.ToTransferRecord(res)

	return so, nil
}

// UpdateTransferStatus updates the status of a transfer.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The new status for the transfer.
//
// Returns:
//   - *record.TransferRecord: The updated transfer record.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*record.TransferRecord, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateTransferStatus(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferStatusFailed
	}

	so := r.mapper.ToTransferRecord(res)

	return so, nil
}

// TrashedTransfer marks a transfer as deleted (soft delete).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_id: The ID of the transfer to be trashed.
//
// Returns:
//   - *record.TransferRecord: The trashed transfer record.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) TrashedTransfer(ctx context.Context, transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.TrashTransfer(ctx, int32(transfer_id))

	if err != nil {
		return nil, transfer_errors.ErrTrashedTransferFailed
	}

	so := r.mapper.ToTransferRecord(res)

	return so, nil
}

// RestoreTransfer restores a previously trashed transfer.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_id: The ID of the transfer to be restored.
//
// Returns:
//   - *record.TransferRecord: The restored transfer record.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) RestoreTransfer(ctx context.Context, transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.RestoreTransfer(ctx, int32(transfer_id))
	if err != nil {
		return nil, transfer_errors.ErrRestoreTransferFailed
	}

	so := r.mapper.ToTransferRecord(res)

	return so, nil
}

// DeleteTransferPermanent permanently deletes a transfer from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_id: The ID of the transfer to be permanently deleted.
//
// Returns:
//   - bool: Indicates if the deletion was successful.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(ctx, int32(transfer_id))
	if err != nil {
		return false, transfer_errors.ErrDeleteTransferPermanentFailed
	}
	return true, nil
}

// RestoreAllTransfer restores all trashed transfers.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Indicates if the restore operation was successful.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) RestoreAllTransfer(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransfers(ctx)

	if err != nil {
		return false, transfer_errors.ErrRestoreAllTransfersFailed
	}

	return true, nil
}

// DeleteAllTransferPermanent permanently deletes all trashed transfers.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Indicates if the deletion was successful.
//   - error: Any error encountered during the operation.
func (r *transferCommandRepository) DeleteAllTransferPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(ctx)

	if err != nil {
		return false, transfer_errors.ErrDeleteAllTransfersPermanentFailed
	}

	return true, nil
}

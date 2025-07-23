package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw"
)

// withdrawCommandRepository is a struct that implements the WithdrawCommandRepository interface
type withdrawCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.WithdrawCommandRecordMapping
}

// NewWithdrawCommandRepository initializes a new instance of withdrawCommandRepository.
// This repository is responsible for executing command operations related to withdraw records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - ctx: The context to be used for database operations, allowing for cancellation and timeout.
//   - mapper: A WithdrawRecordMapping that provides methods to map database rows to WithdrawRecord domain models.
//
// Returns:
//   - A pointer to the newly created withdrawCommandRepository instance.
func NewWithdrawCommandRepository(db *db.Queries, mapper recordmapper.WithdrawCommandRecordMapping) WithdrawCommandRepository {
	return &withdrawCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateWithdraw inserts a new withdraw record into the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing withdraw data.
//
// Returns:
//   - *record.WithdrawRecord: The newly created withdraw record.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.CreateWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.CreateWithdraw(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrCreateWithdrawFailed
	}

	so := r.mapper.ToWithdrawRecord(res)

	return so, nil
}

// UpdateWithdraw modifies an existing withdraw record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing updated withdraw data.
//
// Returns:
//   - *record.WithdrawRecord: The updated withdraw record.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawParams{
		WithdrawID:     int32(*request.WithdrawID),
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.UpdateWithdraw(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawFailed
	}

	so := r.mapper.ToWithdrawRecord(res)

	return so, nil
}

// UpdateWithdrawStatus updates the status of a withdraw record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing the new status.
//
// Returns:
//   - *record.WithdrawRecord: The updated withdraw record with the new status.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawStatusParams{
		WithdrawID: int32(request.WithdrawID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateWithdrawStatus(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawStatusFailed
	}

	so := r.mapper.ToWithdrawRecord(res)

	return so, nil
}

// TrashedWithdraw soft-deletes a withdraw record by marking it as trashed.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - WithdrawID: The ID of the withdraw to be trashed.
//
// Returns:
//   - *record.WithdrawRecord: The trashed withdraw record.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) TrashedWithdraw(ctx context.Context, withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.TrashWithdraw(ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrTrashedWithdrawFailed
	}

	so := r.mapper.ToWithdrawRecord(res)

	return so, nil
}

// RestoreWithdraw restores a previously trashed withdraw record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - WithdrawID: The ID of the withdraw to be restored.
//
// Returns:
//   - *record.WithdrawRecord: The restored withdraw record.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) RestoreWithdraw(ctx context.Context, withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.RestoreWithdraw(ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrRestoreWithdrawFailed
	}

	so := r.mapper.ToWithdrawRecord(res)

	return so, nil
}

// DeleteWithdrawPermanent permanently deletes a withdraw record from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - WithdrawID: The ID of the withdraw to be permanently deleted.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(ctx, int32(withdraw_id))

	if err != nil {
		return false, withdraw_errors.ErrDeleteWithdrawPermanentFailed
	}

	return true, nil
}

// RestoreAllWithdraw restores all trashed withdraw records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the restoration was successful.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) RestoreAllWithdraw(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllWithdraws(ctx)

	if err != nil {
		return false, withdraw_errors.ErrRestoreAllWithdrawsFailed
	}

	return true, nil
}

// DeleteAllWithdrawPermanent permanently deletes all trashed withdraw records from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - error: An error if the operation fails.
func (r *withdrawCommandRepository) DeleteAllWithdrawPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(ctx)

	if err != nil {
		return false, withdraw_errors.ErrDeleteAllWithdrawsPermanentFailed
	}

	return true, nil
}

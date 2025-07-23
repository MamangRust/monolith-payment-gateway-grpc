package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup"
)

// topupCommandRepository is a struct that implements the TopupCommandRepository interface
type topupCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupCommandRecordMapping
}

// NewTopupCommandRepository initializes a new instance of topupCommandRepository with the provided
// database queries, context, and topup record mapper. This repository is responsible for executing
// command operations related to topup records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A TopupRecordMapping that provides methods to map database rows to TopupRecord domain models.
//
// Returns:
//   - A pointer to the newly created topupCommandRepository instance.
func NewTopupCommandRepository(db *db.Queries, mapper recordmapper.TopupCommandRecordMapping) TopupCommandRepository {
	return &topupCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateTopup inserts a new topup record into the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The data needed to create a new topup.
//
// Returns:
//   - *record.TopupRecord: The created topup record.
//   - error: Error if creation fails.
func (r *topupCommandRepository) CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*record.TopupRecord, error) {
	req := db.CreateTopupParams{
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.CreateTopup(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrCreateTopupFailed
	}

	return r.mapper.ToTopupRecord(res), nil
}

// UpdateTopup updates an existing topup record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The data used to update the topup.
//
// Returns:
//   - *record.TopupRecord: The updated topup record.
//   - error: Error if update fails.
func (r *topupCommandRepository) UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*record.TopupRecord, error) {
	req := db.UpdateTopupParams{
		TopupID:     int32(*request.TopupID),
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.UpdateTopup(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupFailed
	}

	return r.mapper.ToTopupRecord(res), nil
}

// UpdateTopupAmount updates the amount of a specific topup.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The new amount data.
//
// Returns:
//   - *record.TopupRecord: The updated topup record.
//   - error: Error if update fails.
func (r *topupCommandRepository) UpdateTopupAmount(ctx context.Context, request *requests.UpdateTopupAmount) (*record.TopupRecord, error) {
	req := db.UpdateTopupAmountParams{
		TopupID:     int32(request.TopupID),
		TopupAmount: int32(request.TopupAmount),
	}

	res, err := r.db.UpdateTopupAmount(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupAmountFailed
	}

	return r.mapper.ToTopupRecord(res), nil
}

// UpdateTopupStatus updates the status of a topup (e.g., success, failed).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The status update data.
//
// Returns:
//   - *record.TopupRecord: The updated topup record.
//   - error: Error if update fails.
func (r *topupCommandRepository) UpdateTopupStatus(ctx context.Context, request *requests.UpdateTopupStatus) (*record.TopupRecord, error) {
	req := db.UpdateTopupStatusParams{
		TopupID: int32(request.TopupID),
		Status:  request.Status,
	}

	res, err := r.db.UpdateTopupStatus(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupStatusFailed
	}

	return r.mapper.ToTopupRecord(res), nil
}

// TrashedTopup soft deletes a topup by marking it as trashed.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topup_id: The ID of the topup to trash.
//
// Returns:
//   - *record.TopupRecord: The trashed topup record.
//   - error: Error if trashing fails.
func (r *topupCommandRepository) TrashedTopup(ctx context.Context, topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.TrashTopup(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrTrashedTopupFailed
	}
	return r.mapper.ToTopupRecord(res), nil
}

// RestoreTopup restores a previously trashed topup.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topup_id: The ID of the topup to restore.
//
// Returns:
//   - *record.TopupRecord: The restored topup record.
//   - error: Error if restoration fails.
func (r *topupCommandRepository) RestoreTopup(ctx context.Context, topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.RestoreTopup(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrRestoreTopupFailed
	}
	return r.mapper.ToTopupRecord(res), nil
}

// DeleteTopupPermanent permanently deletes a topup from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topup_id: The ID of the topup to delete.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - error: Error if deletion fails.
func (r *topupCommandRepository) DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error) {
	err := r.db.DeleteTopupPermanently(ctx, int32(topup_id))
	if err != nil {
		return false, topup_errors.ErrDeleteTopupPermanentFailed
	}
	return true, nil
}

// RestoreAllTopup restores all trashed topups in the system.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the restoration was successful.
//   - error: Error if operation fails.
func (r *topupCommandRepository) RestoreAllTopup(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTopups(ctx)

	if err != nil {
		return false, topup_errors.ErrRestoreAllTopupFailed
	}

	return true, nil
}

// DeleteAllTopupPermanent permanently deletes all trashed topups from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - error: Error if deletion fails.
func (r *topupCommandRepository) DeleteAllTopupPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTopups(ctx)

	if err != nil {
		return false, topup_errors.ErrDeleteAllTopupPermanentFailed
	}

	return true, nil
}

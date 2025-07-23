package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

// saldoCommandRepository is a struct that implements the SaldoCommandRepository interface
type saldoCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.SaldoCommandRecordMapping
}

// NewSaldoCommandRepository initializes a new instance of saldoCommandRepository with the provided
// database queries, context, and saldo record mapper. This repository is responsible for executing
// command operations related to saldo records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A SaldoRecordMapping that provides methods to map database rows to SaldoRecord domain models.
//
// Returns:
//   - A pointer to the newly created saldoCommandRepository instance.
func NewSaldoCommandRepository(db *db.Queries, mapper recordmapper.SaldoCommandRecordMapping) SaldoCommandRepository {
	return &saldoCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateSaldo inserts a new saldo record into the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing saldo data to be created.
//
// Returns:
//   - *record.SaldoRecord: The created saldo record.
//   - error: An error if the insert operation fails.
func (r *saldoCommandRepository) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*record.SaldoRecord, error) {
	req := db.CreateSaldoParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}
	res, err := r.db.CreateSaldo(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrCreateSaldoFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// UpdateSaldo updates the saldo record by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing updated saldo data.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: An error if the update operation fails.
func (r *saldoCommandRepository) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoParams{
		SaldoID:      int32(*request.SaldoID),
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldo(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// UpdateSaldoBalance updates the saldo balance (e.g., after top-up or transaction).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing the updated balance info.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: An error if the update fails.
func (r *saldoCommandRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// UpdateSaldoWithdraw updates the saldo after a withdrawal operation.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing withdrawal update information.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: An error if the update fails.
func (r *saldoCommandRepository) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error) {
	withdrawAmount := sql.NullInt32{
		Int32: int32(*request.WithdrawAmount),
		Valid: request.WithdrawAmount != nil,
	}
	var withdrawTime sql.NullTime
	if request.WithdrawTime != nil {
		withdrawTime = sql.NullTime{
			Time:  *request.WithdrawTime,
			Valid: true,
		}
	}

	req := db.UpdateSaldoWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: withdrawAmount,
		WithdrawTime:   withdrawTime,
	}

	res, err := r.db.UpdateSaldoWithdraw(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoWithdrawFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// TrashedSaldo marks a saldo record as soft-deleted.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldoID: The ID of the saldo to be soft-deleted.
//
// Returns:
//   - *record.SaldoRecord: The trashed saldo record.
//   - error: An error if the operation fails.
func (r *saldoCommandRepository) TrashedSaldo(ctx context.Context, saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.TrashSaldo(ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrTrashSaldoFailed
	}
	return r.mapper.ToSaldoRecord(res), nil
}

// RestoreSaldo restores a previously trashed saldo record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldoID: The ID of the saldo to be restored.
//
// Returns:
//   - *record.SaldoRecord: The restored saldo record.
//   - error: An error if the operation fails.
func (r *saldoCommandRepository) RestoreSaldo(ctx context.Context, saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.RestoreSaldo(ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrRestoreSaldoFailed
	}
	return r.mapper.ToSaldoRecord(res), nil
}

// DeleteSaldoPermanent permanently deletes a saldo record from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo to be permanently deleted.
//
// Returns:
//   - bool: True if deletion was successful.
//   - error: An error if the operation fails.
func (r *saldoCommandRepository) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error) {
	err := r.db.DeleteSaldoPermanently(ctx, int32(saldo_id))
	if err != nil {
		return false, saldo_errors.ErrDeleteSaldoPermanentFailed
	}
	return true, nil
}

// RestoreAllSaldo restores all trashed saldo records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if all saldo records were successfully restored.
//   - error: An error if the operation fails.
func (r *saldoCommandRepository) RestoreAllSaldo(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllSaldos(ctx)

	if err != nil {
		return false, saldo_errors.ErrRestoreAllSaldosFailed
	}

	return true, nil
}

// DeleteAllSaldoPermanent permanently deletes all trashed saldo records from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if all trashed saldo records were successfully deleted.
//   - error: An error if the operation fails.
func (r *saldoCommandRepository) DeleteAllSaldoPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentSaldos(ctx)

	if err != nil {
		return false, saldo_errors.ErrDeleteAllSaldosPermanentFailed
	}

	return true, nil
}

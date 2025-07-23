package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction"
)

// transactionCommandRepository is a repository for handling transaction operations
type transactionCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionCommandRecordMapper
}

// NewTransactionCommandRepository initializes a new instance of transactionCommandRepository with the provided
// database queries, context, and transaction record mapper. This repository is responsible for executing
// command operations related to transaction records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A TransactionRecordMapping that provides methods to map database rows to TransactionRecord domain models.
//
// Returns:
//   - A pointer to the newly created transactionCommandRepository instance.
func NewTransactionCommandRepository(db *db.Queries, mapper recordmapper.TransactionCommandRecordMapper) TransactionCommandRepository {
	return &transactionCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateTransaction creates a new transaction record in the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The data to create the transaction.
//
// Returns:
//   - *record.TransactionRecord: The created transaction record.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.CreateTransactionParams{
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.CreateTransaction(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed
	}

	return r.mapper.ToTransactionRecord(res), nil
}

// UpdateTransaction updates an existing transaction record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The data to update the transaction.
//
// Returns:
//   - *record.TransactionRecord: The updated transaction record.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionParams{
		TransactionID:   int32(*request.TransactionID),
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.UpdateTransaction(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionFailed
	}

	return r.mapper.ToTransactionRecord(res), nil
}

// UpdateTransactionStatus updates only the status of a transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The status update request.
//
// Returns:
//   - *record.TransactionRecord: The transaction with updated status.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) UpdateTransactionStatus(ctx context.Context, request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionStatusParams{
		TransactionID: int32(request.TransactionID),
		Status:        request.Status,
	}

	res, err := r.db.UpdateTransactionStatus(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed
	}

	return r.mapper.ToTransactionRecord(res), nil
}

// TrashedTransaction marks a transaction as soft-deleted.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transaction_id: ID of the transaction to trash.
//
// Returns:
//   - *record.TransactionRecord: The trashed transaction record.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) TrashedTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrTrashedTransactionFailed
	}
	return r.mapper.ToTransactionRecord(res), nil
}

// RestoreTransaction restores a soft-deleted transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topup_id: ID of the transaction to restore.
//
// Returns:
//   - *record.TransactionRecord: The restored transaction record.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransactionFailed
	}
	return r.mapper.ToTransactionRecord(res), nil
}

// DeleteTransactionPermanent deletes a transaction permanently.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topup_id: ID of the transaction to delete permanently.
//
// Returns:
//   - bool: True if deletion was successful.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(ctx, int32(transaction_id))
	if err != nil {

		return false, transaction_errors.ErrDeleteTransactionPermanentFailed
	}
	return true, nil
}

// RestoreAllTransaction restores all soft-deleted transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if restore operation was successful.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) RestoreAllTransaction(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransactions(ctx)

	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactionsFailed
	}

	return true, nil
}

// DeleteAllTransactionPermanent deletes all transactions permanently.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if deletion was successful.
//   - error: Error if operation fails.
func (r *transactionCommandRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(ctx)

	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionsPermanentFailed
	}
	return true, nil
}

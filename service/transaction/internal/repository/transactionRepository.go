package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transactionCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionCommandRepository {
	return &transactionCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionCommandRepository) CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.CreateTransactionParams{
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.CreateTransaction(r.ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionParams{
		TransactionID:   int32(*request.TransactionID),
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.UpdateTransaction(r.ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) UpdateTransactionStatus(request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionStatusParams{
		TransactionID: int32(request.TransactionID),
		Status:        request.Status,
	}

	res, err := r.db.UpdateTransactionStatus(r.ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) TrashedTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(r.ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrTrashedTransactionFailed
	}
	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) RestoreTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(r.ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransactionFailed
	}
	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) DeleteTransactionPermanent(transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(r.ctx, int32(transaction_id))
	if err != nil {

		return false, transaction_errors.ErrDeleteTransactionPermanentFailed
	}
	return true, nil
}

func (r *transactionCommandRepository) RestoreAllTransaction() (bool, error) {
	err := r.db.RestoreAllTransactions(r.ctx)

	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactionsFailed
	}

	return true, nil
}

func (r *transactionCommandRepository) DeleteAllTransactionPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(r.ctx)

	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionsPermanentFailed
	}
	return true, nil
}

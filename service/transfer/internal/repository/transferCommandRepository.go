package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transferCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransferRecordMapping
}

func NewTransferCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransferRecordMapping) *transferCommandRepository {
	return &transferCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transferCommandRepository) CreateTransfer(request *requests.CreateTransferRequest) (*record.TransferRecord, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrCreateTransferFailed
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferCommandRepository) UpdateTransfer(request *requests.UpdateTransferRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferParams{
		TransferID:     int32(*request.TransferID),
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransfer(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferFailed
	}

	return r.mapping.ToTransferRecord(res), nil

}

func (r *transferCommandRepository) UpdateTransferAmount(request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransferAmount(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferAmountFailed
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferCommandRepository) UpdateTransferStatus(request *requests.UpdateTransferStatus) (*record.TransferRecord, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateTransferStatus(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferStatusFailed
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferCommandRepository) TrashedTransfer(transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.TrashTransfer(r.ctx, int32(transfer_id))

	if err != nil {
		return nil, transfer_errors.ErrTrashedTransferFailed
	}
	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferCommandRepository) RestoreTransfer(transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.RestoreTransfer(r.ctx, int32(transfer_id))
	if err != nil {
		return nil, transfer_errors.ErrRestoreTransferFailed
	}
	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferCommandRepository) DeleteTransferPermanent(transfer_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(r.ctx, int32(transfer_id))
	if err != nil {
		return false, transfer_errors.ErrDeleteTransferPermanentFailed
	}
	return true, nil
}

func (r *transferCommandRepository) RestoreAllTransfer() (bool, error) {
	err := r.db.RestoreAllTransfers(r.ctx)

	if err != nil {
		return false, transfer_errors.ErrRestoreAllTransfersFailed
	}

	return true, nil
}

func (r *transferCommandRepository) DeleteAllTransferPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(r.ctx)

	if err != nil {
		return false, transfer_errors.ErrDeleteAllTransfersPermanentFailed
	}

	return true, nil
}

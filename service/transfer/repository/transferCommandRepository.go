package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferCommandRepository struct {
	db *db.Queries
}

func NewTransferCommandRepository(db *db.Queries) TransferCommandRepository {
	return &transferCommandRepository{
		db: db,
	}
}

func (r *transferCommandRepository) CreateTransfer(ctx context.Context, request *requests.CreateTransferRequest) (*db.CreateTransferRow, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrCreateTransferFailed
	}

	return res, nil
}

func (r *transferCommandRepository) UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*db.UpdateTransferRow, error) {
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

	return res, nil

}

func (r *transferCommandRepository) UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*db.UpdateTransferAmountRow, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransferAmount(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferAmountFailed
	}

	return res, nil
}

func (r *transferCommandRepository) UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*db.UpdateTransferStatusRow, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateTransferStatus(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferStatusFailed
	}

	return res, nil
}

func (r *transferCommandRepository) TrashedTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error) {
	res, err := r.db.TrashTransfer(ctx, int32(transfer_id))

	if err != nil {
		return nil, transfer_errors.ErrTrashedTransferFailed
	}
	return res, nil
}

func (r *transferCommandRepository) RestoreTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error) {
	res, err := r.db.RestoreTransfer(ctx, int32(transfer_id))
	if err != nil {
		return nil, transfer_errors.ErrRestoreTransferFailed
	}
	return res, nil
}

func (r *transferCommandRepository) DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(ctx, int32(transfer_id))
	if err != nil {
		return false, transfer_errors.ErrDeleteTransferPermanentFailed
	}
	return true, nil
}

func (r *transferCommandRepository) RestoreAllTransfer(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransfers(ctx)

	if err != nil {
		return false, transfer_errors.ErrRestoreAllTransfersFailed
	}

	return true, nil
}

func (r *transferCommandRepository) DeleteAllTransferPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(ctx)

	if err != nil {
		return false, transfer_errors.ErrDeleteAllTransfersPermanentFailed
	}

	return true, nil
}

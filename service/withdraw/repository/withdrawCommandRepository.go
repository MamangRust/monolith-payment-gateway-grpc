package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawCommandRepository struct {
	db *db.Queries
}

func NewWithdrawCommandRepository(db *db.Queries) WithdrawCommandRepository {
	return &withdrawCommandRepository{
		db: db,
	}
}

func (r *withdrawCommandRepository) CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*db.CreateWithdrawRow, error) {
	req := db.CreateWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.CreateWithdraw(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrCreateWithdrawFailed
	}

	return res, nil
}

func (r *withdrawCommandRepository) UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*db.UpdateWithdrawRow, error) {
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

	return res, nil
}

func (r *withdrawCommandRepository) UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*db.UpdateWithdrawStatusRow, error) {
	req := db.UpdateWithdrawStatusParams{
		WithdrawID: int32(request.WithdrawID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateWithdrawStatus(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawStatusFailed
	}

	return res, nil
}

func (r *withdrawCommandRepository) TrashedWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error) {
	res, err := r.db.TrashWithdraw(ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrTrashedWithdrawFailed
	}

	return res, nil
}

func (r *withdrawCommandRepository) RestoreWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error) {
	res, err := r.db.RestoreWithdraw(ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrRestoreWithdrawFailed
	}

	return res, nil
}

func (r *withdrawCommandRepository) DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(ctx, int32(withdraw_id))

	if err != nil {
		return false, withdraw_errors.ErrDeleteWithdrawPermanentFailed
	}

	return true, nil
}

func (r *withdrawCommandRepository) RestoreAllWithdraw(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllWithdraws(ctx)

	if err != nil {
		return false, withdraw_errors.ErrRestoreAllWithdrawsFailed
	}

	return true, nil
}

func (r *withdrawCommandRepository) DeleteAllWithdrawPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(ctx)

	if err != nil {
		return false, withdraw_errors.ErrDeleteAllWithdrawsPermanentFailed
	}

	return true, nil
}

package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type withdrawCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.WithdrawRecordMapping
}

func NewWithdrawCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.WithdrawRecordMapping) *withdrawCommandRepository {
	return &withdrawCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *withdrawCommandRepository) CreateWithdraw(request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.CreateWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.CreateWithdraw(r.ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrCreateWithdrawFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawCommandRepository) UpdateWithdraw(request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawParams{
		WithdrawID:     int32(*request.WithdrawID),
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.UpdateWithdraw(r.ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawCommandRepository) UpdateWithdrawStatus(request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawStatusParams{
		WithdrawID: int32(request.WithdrawID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateWithdrawStatus(r.ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawStatusFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawCommandRepository) TrashedWithdraw(withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.TrashWithdraw(r.ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrTrashedWithdrawFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawCommandRepository) RestoreWithdraw(withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.RestoreWithdraw(r.ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrRestoreWithdrawFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawCommandRepository) DeleteWithdrawPermanent(withdraw_id int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(r.ctx, int32(withdraw_id))

	if err != nil {
		return false, withdraw_errors.ErrDeleteWithdrawPermanentFailed
	}

	return true, nil
}

func (r *withdrawCommandRepository) RestoreAllWithdraw() (bool, error) {
	err := r.db.RestoreAllWithdraws(r.ctx)

	if err != nil {
		return false, withdraw_errors.ErrRestoreAllWithdrawsFailed
	}

	return true, nil
}

func (r *withdrawCommandRepository) DeleteAllWithdrawPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(r.ctx)

	if err != nil {
		return false, withdraw_errors.ErrDeleteAllWithdrawsPermanentFailed
	}

	return true, nil
}

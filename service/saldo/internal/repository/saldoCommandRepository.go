package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type saldoCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SaldoRecordMapping
}

func NewSaldoCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SaldoRecordMapping) *saldoCommandRepository {
	return &saldoCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *saldoCommandRepository) CreateSaldo(request *requests.CreateSaldoRequest) (*record.SaldoRecord, error) {
	req := db.CreateSaldoParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}
	res, err := r.db.CreateSaldo(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrCreateSaldoFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoCommandRepository) UpdateSaldo(request *requests.UpdateSaldoRequest) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoParams{
		SaldoID:      int32(*request.SaldoID),
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldo(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoCommandRepository) UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoCommandRepository) UpdateSaldoWithdraw(request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error) {
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

	res, err := r.db.UpdateSaldoWithdraw(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoWithdrawFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoCommandRepository) TrashedSaldo(saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.TrashSaldo(r.ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrTrashSaldoFailed
	}
	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoCommandRepository) RestoreSaldo(saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.RestoreSaldo(r.ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrRestoreSaldoFailed
	}
	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoCommandRepository) DeleteSaldoPermanent(saldo_id int) (bool, error) {
	err := r.db.DeleteSaldoPermanently(r.ctx, int32(saldo_id))
	if err != nil {
		return false, saldo_errors.ErrDeleteSaldoPermanentFailed
	}
	return true, nil
}

func (r *saldoCommandRepository) RestoreAllSaldo() (bool, error) {
	err := r.db.RestoreAllSaldos(r.ctx)

	if err != nil {
		return false, saldo_errors.ErrRestoreAllSaldosFailed
	}

	return true, nil
}

func (r *saldoCommandRepository) DeleteAllSaldoPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentSaldos(r.ctx)

	if err != nil {
		return false, saldo_errors.ErrDeleteAllSaldosPermanentFailed
	}

	return true, nil
}
